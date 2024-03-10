package request

import "github.com/gin-gonic/gin"

type flashevent struct{}

var FlashEvent = new(flashevent)

// 用户信息预热
func (f *flashevent) PreheatUserInfo(c *gin.Context) {

}

// 商品信息预热
func (f *flashevent) PreheatProductInfo(c *gin.Context) {

}

// 下单
func (f *flashevent) PlaceOrder(c *gin.Context) {

}
