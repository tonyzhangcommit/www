package request

import (
	"auth/global"
	"auth/response"
	"auth/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type flashevent struct{}

var FlashEvent = new(flashevent)

// 用户信息预热
func (f *flashevent) PreheatUserInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.AdminUrl.Preheat)
	GetRequest(c, 10, remoteurl)
}

// 获取用户等级信息
func (f *flashevent) GetUserLevelInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Getuvip)
	PostRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 活动信息预热
func (f *flashevent) PreheatProductInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.ProductServiceApi.BaseUrl, global.App.Config.ProductServiceApi.FlashGetFEinfo)
	PostRequest(c, global.App.Config.ProductServiceApi.Timeout, remoteurl)
}

// 秒杀活动&商品展示
func (f *flashevent) PreheatEventProductShow(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.ProductServiceApi.BaseUrl, global.App.Config.ProductServiceApi.FlashGetEventP)
	PostRequest(c, global.App.Config.ProductServiceApi.Timeout, remoteurl)
}

/*
秒杀活动下单，这里为了防止认证服务宕机，请求前需要对请求进行限流
*/
func (f *flashevent) PlaceOrder(c *gin.Context) {
	// 首先对参数做基本验证,验证之前需要保存请求体
	body, errofreadbody := io.ReadAll(c.Request.Body)

	if errofreadbody != nil {
		response.IllegalRequestFail(c)
		return
	}
	// 将请求参数重新放到请求中
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	// 验证参数
	var form TakeFlashOrder
	if errvarify := c.ShouldBindJSON(&form); errvarify != nil {
		response.VarifyErrorFail(c, GetErrorMsg(form, errvarify))
		return
	}
	// 将请求参数重新放到请求中
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 获取用户信息&活动信息
	userremoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Getuvip)
	eventremoteurl := utils.JoinStrings(global.App.Config.ProductServiceApi.BaseUrl, global.App.Config.ProductServiceApi.FlashGetFEinfo)

	var userRes response.UserVipType
	var eventRes response.FlashEventRes
	var userResError, eventResError error

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		userRes, userResError = getUserLevel(userremoteurl, form.UserID)
	}()
	go func() {
		defer wg.Done()
		eventRes, eventResError = getBaseEvent(eventremoteurl, form.EventID)
	}()
	wg.Wait()

	// 首次判断请求是否合法
	if err := checkIslawful(userRes, eventRes, &form, userResError, eventResError); err != nil {
		response.PermissionFail(c, err.Error())
	} else {
		remoteurl := utils.JoinStrings(global.App.Config.OrderServiceApi.BaseUrl, global.App.Config.OrderServiceApi.TakeFalshOrder)
		PostRequest(c, global.App.Config.ProductServiceApi.Timeout, remoteurl)
	}
}

// 获取用户会员级别
func getUserLevel(remoteurl string, userid uint) (res response.UserVipType, err error) {
	type form struct {
		Uid uint `json:"uid"`
	}
	data := form{
		Uid: userid,
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		err = errors.New("非法请求")
	}
	// 创建客户端
	requester := global.NewRequestClient(3 * time.Second)
	// 创建请求
	req, err := http.NewRequest("POST", remoteurl, bytes.NewBuffer(jsondata))
	if err != nil {
		global.SendLogs("error", "获取用户会员级别-创建请求错误", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发出请求
	resp, err := requester.Client.Do(req)
	if err != nil {
		global.SendLogs("error", "获取用户会员级别-发送请求错误", err)
		return
	}
	// 读取请求体
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		global.SendLogs("error", "获取用户会员级别-读取response错误", err)
		return
	}
	// 序列化至指定结构体
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		global.SendLogs("error", "获取用户会员级别-解析response错误", err)
		return
	}
	return
}

// 获取活动基本信息
func getBaseEvent(remoteurl string, eventid uint) (res response.FlashEventRes, err error) {
	type form struct {
		EventId uint `form:"eventid" json:"eventid"`
	}
	data := form{
		EventId: eventid,
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		err = errors.New("非法请求")
	}
	requester := global.NewRequestClient(3 * time.Second)
	// 创建请求
	req, err := http.NewRequest("POST", remoteurl, bytes.NewBuffer(jsondata))
	if err != nil {
		global.SendLogs("error", "请求秒杀活动基本信息-创建请求错误", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发出请求
	resp, err := requester.Client.Do(req)
	if err != nil {
		global.SendLogs("error", "请求秒杀活动基本信息-发送请求错误", err)
		return
	}
	// 读取请求体
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		global.SendLogs("error", "请求秒杀活动基本信息-读取response错误", err)
		return
	}
	// 序列化至指定结构体
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		global.SendLogs("error", "请求秒杀活动基本信息-解析response错误", err)
		return
	}
	return
}

// 判断当前请求是否合法，主要过滤部分不合理请求，将符合要求的请求转发至订单服务中
func checkIslawful(userres response.UserVipType, eventres response.FlashEventRes, form *TakeFlashOrder, userErr, eventErr error) (err error) {
	if userErr != nil || eventErr != nil {
		err = userErr
		return
	}
	now := time.Now()
	// 判断活动状态（未开始 进行中 已结束）
	key := "flashevent:limituserreq:" + strconv.Itoa(int(form.UserID))
	val := "1"
	dur := time.Hour * 5
	result, err := global.App.Redis.SetNX(context.Background(), key, val, dur).Result()
	if err != nil {
		global.SendLogs("error", "秒杀活动限制用户请求次数错误", err)
		err = errors.New("内部错误")
		return
	}
	if !result {
		// err = errors.New("您已参加过活动了呦~~~")
		return
	} else if now.Before(eventres.Data.StartTime) {
		err = errors.New("活动未开始")
	} else if now.After(eventres.Data.EndTime) {
		err = errors.New("活动已结束")
	} else if global.UserVIPM[userres.Vtype] < global.UserVIPM[eventres.Data.Condition] {
		err = errors.New("非常抱歉，您没有该活动参与权限！")
	} else if form.Count != int(eventres.Data.Count) {
		err = errors.New("参数错误，非法请求") // 这里后续可以将用户信息加入黑名单
	}
	return
}
