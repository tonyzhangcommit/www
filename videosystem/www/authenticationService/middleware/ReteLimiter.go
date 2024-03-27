package middleware

import (
	"auth/global"
	"auth/response"
	"auth/utils"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

/*
	限流中间件
	参数，次数和规定的时间，比如，一分钟限制请求3次，5分钟限制请求10次等
	限流这里考虑两个大方向：
	1. 服务限流
		针对某个服务，定义固定的qps
	2. 接口限流
		针对某个接口，定义访问频率，主要考虑特殊的接口，比如获取验证码，要求同一个IP,同一个手机号一天只能请求5次
*/

/*
令牌桶算法,对服务进行请求限制
1. 服务定制特定的qps
2. 指定ip指定服务设置访问次数（防止恶意访问）
参数：

	r: qps
	bucket 代表桶的最大容量
	ipsumcount 同一个ip在一天访问的时间量
*/

type Limiterinfo struct {
	Limiter  *rate.Limiter
	LastUsed time.Time
}

var (
	IpLimiter sync.Map
)

// ipKey 根据 IP 地址生成一个存储键，只包含 IP 的前两部分
func extractIPParts(ipAddr string) string {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return ""
	}

	if ip.To4() != nil { // IPv4
		parts := strings.Split(ip.String(), ".")
		if len(parts) >= 3 {
			return fmt.Sprintf("%s.%s", parts[0], parts[1])
		}
	} else { // IPv6
		parts := strings.Split(ip.String(), ":")
		if len(parts) >= 4 {
			return fmt.Sprintf("%s:%s:%s:%s", parts[0], parts[1], parts[2], parts[3])
		}
	}
	return ""
}

func ServiceLimit(servicename string, r rate.Limit, bucket int, ipsumcount int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		// 初始化令牌桶，这里为每个IP网段设置limiter,从而做调整
		iplimiterkey := extractIPParts(clientIP)
		if iplimiterkey == "" {
			global.SendLogs("error", fmt.Sprintf("%s:%s", "截取IP错误", clientIP))
		} else {
			val, ok := IpLimiter.Load(iplimiterkey)
			if !ok {
				limiter := rate.NewLimiter(r, bucket)
				val = &Limiterinfo{Limiter: limiter, LastUsed: time.Now()}
				IpLimiter.Store(iplimiterkey, val)
			}
			limiterInfo := val.(*Limiterinfo)
			// 更新最后一次使用时间
			limiterInfo.LastUsed = time.Now()

			if !limiterInfo.Limiter.Allow() {
				go global.SendLogs("error", utils.JoinStrings("IP ", clientIP, "访问频繁终止服务"))
				response.FrequentRequestFail(ctx, "访问频繁,请稍后再试")
				ctx.Abort()
				return
			}
		}
		// ip限制
		if ipsumcount > 0 {
			key := utils.JoinStrings(servicename, ":", "ipLimit:", clientIP)
			res, err := global.App.Redis.Incr(context.Background(), key).Result()
			if err != nil {
				global.SendLogs("error", "服务IP限流添加缓存失败", err)
				response.LocalErrorFail(ctx, "内部服务错误")
				ctx.Abort()
				return
			}
			if res == 1 {
				// 第一次
				now := time.Now()
				endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
				expire := endOfDay.Sub(now)
				err := global.App.Redis.Expire(context.Background(), key, expire).Err()
				if err != nil {
					global.SendLogs("error", "服务IP限流添加过期时间失败", err)
					response.LocalErrorFail(ctx, "内部服务错误")
					ctx.Abort()
				}
			}
			// 检测次数
			if res < int64(ipsumcount) {
				_, err := global.App.Redis.Incr(context.Background(), key).Result()
				if err != nil {
					global.SendLogs("error", "服务IP限流添加缓存失败", err)
					response.LocalErrorFail(ctx, "内部服务错误")
					ctx.Abort()
					return
				}
			} else {
				response.FrequentRequestFail(ctx, "当天请求已达到上限")
				ctx.Abort()
				return
			}
		}
	}
}

// 接口限流中间件,同时具有参数验证的功能
type virifcode struct {
	Phonenumber string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
}

func APIGetVerifCodeLimit(count int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 首先读取原始请求体
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			// 处理错误
			response.IllegalRequestFail(ctx)
			ctx.Abort()
			return
		}
		// 将请求体内容替换回去，以便后续处理可以再次读取
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// 参数检测
		var form virifcode
		if err := ctx.ShouldBindBodyWith(&form, binding.JSON); err != nil {

			response.IllegalRequestFail(ctx)
			ctx.Abort()
			return
		}
		// 首先获取请求IP
		clientIP := ctx.ClientIP()
		key := utils.JoinStrings("getVerifiCodeLimit:", form.Phonenumber, clientIP)
		// 存入redis,如果没有则放入，并且过期时间设定为当天最后一秒
		res, err := global.App.Redis.Incr(context.Background(), key).Result()
		if err != nil {
			global.SendLogs("error", "验证码限流添加缓存失败", err)
			response.LocalErrorFail(ctx, "内部服务错误")
			ctx.Abort()
			return
		}
		if res == 1 {
			// 第一次
			now := time.Now()
			endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
			expire := endOfDay.Sub(now)
			err := global.App.Redis.Expire(context.Background(), key, expire).Err()
			if err != nil {
				global.SendLogs("error", "验证码限流添加过期时间失败", err)
				response.LocalErrorFail(ctx, "内部服务错误")
				ctx.Abort()
			}
		} else if res-1 >= count {
			response.FrequentRequestFail(ctx, "当天请求验证码已达到上限")
			ctx.Abort()
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}
