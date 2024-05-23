package middleware

import (
	"auth/global"
	"auth/response"
	"auth/utils"
	"context"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
jwt 鉴权中间件
包括触发jwt失效的逻辑实现
*/
func JWTAUTH(GuardName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取jwt token
		tokenStr := ctx.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.JwtTokenErrorFail(ctx, "非法请求，缺少token")
			ctx.Abort()
			return
		}
		// 这里获取token 的字符串为bearer **************
		tokenStr = tokenStr[len(response.TokenType)+1:]
		token, err := jwt.ParseWithClaims(tokenStr, &response.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secretkey), nil
		})

		// 增加黑名单校验，主要是用户提出退出登录
		if err != nil || isInBlacklist(tokenStr) {
			response.JwtTokenErrorFail(ctx, "登录验证失败，请重新登录")
			ctx.Abort()
			return
		}
		claims := token.Claims.(*response.CustomClaims)
		if claims.Issuer != GuardName {
			response.JwtTokenErrorFail(ctx, "非法token")
			ctx.Abort()
			return
		}
		// token 续签
		if claims.ExpiresAt-time.Now().Unix() < global.App.Config.Jwt.RefreshGracePeriod {
			// 生成一个锁（实现了Interface接口）
			refresh_token_key := "refresh_token_lock" + claims.Id
			lock := global.Lock(refresh_token_key, global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				// 只有第一次能成功
				tokenData, _ := response.JwtService.CreateJwtToken("app", claims.Id, claims.Roles)
				ctx.Header("new-token", tokenData.AccessToken)
				ctx.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
				// 将旧的token 放入黑名单
				_ = JoinBlackList(token)
			}
		}
		ctx.Set("token", token)
		ctx.Set("userid", claims.Id)
	}
}

// 获取jwt toeken hash key 值
func getBlackListKey(token string) string {
	return "jwt:black:list:" + utils.Sha256([]byte(token))
}

// 用户退出时，将token值加入黑名单
func JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*response.CustomClaims).ExpiresAt-nowUnix) * time.Second
	err = global.App.Redis.SetEX(context.Background(), getBlackListKey(token.Raw), nowUnix, timer).Err()
	return
}

// 判断token是否在黑名单中
func isInBlacklist(tokenStr string) bool {
	joinUnixStr, err := global.App.Redis.Get(context.Background(), getBlackListKey(tokenStr)).Result()
	if err != nil {
		return false
	}
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)

	if joinUnixStr == "" || err != nil {
		return false
	}
	if time.Now().Unix()-joinUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true
}
