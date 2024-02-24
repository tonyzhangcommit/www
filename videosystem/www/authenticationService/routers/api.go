// 定义不同服务api接口

package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户管理服务api管理
func SetUserServiceGroupRouter(router *gin.RouterGroup) {
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
}

// 商品管理服务api管理