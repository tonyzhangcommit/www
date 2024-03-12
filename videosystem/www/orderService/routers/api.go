package routers

import (
	"userservice/app/management"

	"github.com/gin-gonic/gin"
)

func SetOrderGroupRouter(router *gin.RouterGroup) {
	router.POST("/takeflashorder", management.TakeFlashOrder)
}
