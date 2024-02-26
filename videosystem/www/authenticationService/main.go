package main

import (
	"auth/bootstrap"
)

func main() {
	//  初始化配置文件
	bootstrap.InitializeConfig()
	// 初始化rabbitMQ
	bootstrap.InitializeRabbitMQ()
	// 初始化缓存
	bootstrap.InitializeRedis()
	// 初始化用户管理服务
	bootstrap.RunServer()
}
