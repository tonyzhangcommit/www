package request

import (
	"auth/global"
	"auth/utils"

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

// 商品信息预热
func (f *flashevent) PreheatProductInfo(c *gin.Context) {

}

// 下单
func (f *flashevent) PlaceOrder(c *gin.Context) {

}
