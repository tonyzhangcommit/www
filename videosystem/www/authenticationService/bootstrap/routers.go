// 初始化项目启动路由配置
package bootstrap

import (
	"auth/global"
	"auth/routers"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	// 这里是定义与用户管理，商品管理，订单，支付服务交互的路由组
	apiGroup := router.Group("/api")
	// 设置用户服务-client路由组
	routers.SetUserServiceClientGroupRouter(apiGroup)
	// 设置用户服务-admin路由组
	routers.SetUserServiceManageGroupRouter(apiGroup)
	// 设置商品管理服务api路由组
	routers.SetProductServiceManageGroupRouter(apiGroup)
	// 设置订单服务api路由组
	routers.SetOrderServiceManageGroupRouter(apiGroup)
	// 设置秒杀活动api路由组
	routers.SetFlashEventServiceManageGroupRouter(apiGroup)
	return router
}

func RunServer() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := setupRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
