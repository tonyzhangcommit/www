package main

import (
	"userservice/app/services"
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
	// 启动生成订单消费者
	go services.FlashEventCustomer(10)
	go services.FlashEventsnapupresCustomer()
	// go services.FlashEventsnapupresCustomer()
	// go services.FlashEventsnapupresCustomer()

	// 启动微服务
	bootstrap.RunServer()
}
