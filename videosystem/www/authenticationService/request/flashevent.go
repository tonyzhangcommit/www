package request

import (
	"auth/global"
	"auth/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type flashevent struct{}

var FlashEvent = new(flashevent)

// 用户信息预热
func (f *flashevent) PreheatUserInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.AdminUrl.Preheat)
	GetRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 获取用户等级信息
func (f *flashevent) GetUserLevelInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Getuvip)
	PostRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 活动信息预热
func (f *flashevent) PreheatProductInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.ProductServiceApi.BaseUrl, global.App.Config.ProductServiceApi.FlashGetFEinfo)
	fmt.Println(remoteurl)
	PostRequest(c, global.App.Config.ProductServiceApi.Timeout, remoteurl)
}

// 秒杀活动&商品展示
func (f *flashevent) PreheatEventProductShow(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.ProductServiceApi.BaseUrl, global.App.Config.ProductServiceApi.FlashGetEventP)
	PostRequest(c, global.App.Config.ProductServiceApi.Timeout, remoteurl)
}

// 下单
func (f *flashevent) PlaceOrder(c *gin.Context) {

}
