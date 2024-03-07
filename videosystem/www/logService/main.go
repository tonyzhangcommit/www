package main

import (
	"fmt"
	"logservice/bootstrap"
)

func main() {
	fmt.Println("logservice start!")
	// 初始化配置文件
	bootstrap.InitializeConfig()
	// 初始日志服务
	bootstrap.InitializeLocalLogger()
	// 初始化用户服务日志器
	bootstrap.InitializeUserServiceLogger()
	// 初始化认证服务日志器
	bootstrap.InitializeAuthServiceLogger()
	// 初始化产品服务日至器
	bootstrap.InitializeProductLogger()
	// 初始化订单服务日志器
	bootstrap.InitializeOrderLogger()
	// 日志处理开始
	bootstrap.InitRabbitMQ()
	var ch chan int
	<-ch
}
