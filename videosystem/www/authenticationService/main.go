package main

import (
	"auth/bootstrap"
	timingtask "auth/timingTask"
)

func main() {
	//  初始化配置文件
	bootstrap.InitializeConfig()
	// 初始化自定义验证器
	bootstrap.InitializeValidator()
	// 初始化rabbitMQ
	bootstrap.InitializeRabbitMQ()
	// 初始化缓存
	bootstrap.InitializeRedis()
	// 开启定时任务
	go timingtask.CleanupLimiters()
	// 初始化用户管理服务
	bootstrap.RunServer()
}
