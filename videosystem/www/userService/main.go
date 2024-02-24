package main

import (
	"userservice/bootstrap"
)

func main() {
	//初始化加载配置文件
	bootstrap.InitializeConfig()
	// 初始化验证器
	bootstrap.InitializeValidator()
	// 初始化rabbitMQ
	bootstrap.InitializeRabbitMQ()
	// 初始化 database & cache
	bootstrap.InitializeDatabase()
	bootstrap.InitializeRedis()
	// 启动微服务
	bootstrap.RunServer()
}
