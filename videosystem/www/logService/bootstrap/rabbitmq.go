package bootstrap

import (
	"logservice/global"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

/*
	这里封装处理不同微服务发送过来的日志功能封装
	1. 初始化rabbitMQ服务,包括初始化不同服务交换机，队列名
	2. 监听不同队列消息进行不同的处理
*/

func InitRabbitMQ() {
	global.App.LogsServiceLogger.Info("开始初始化rabbitMQ")
	rabbitUsername := global.App.Config.RabbitMQ.User
	rabbitPassword := global.App.Config.RabbitMQ.Password
	rabbitPort := global.App.Config.RabbitMQ.Port
	rabbitHost := global.App.Config.RabbitMQ.Host
	rabbitVhost := global.App.Config.RabbitMQ.Vhost
	amqpStr := "amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/" + rabbitVhost
	// 创建连接
	conn, err := amqp.Dial(amqpStr)
	if err != nil {
		global.App.LogsServiceLogger.Error(err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		global.App.LogsServiceLogger.Error(err.Error())
	}
	go InitMQService(ch, global.App.Config.UserviceConfig.Userexchange, global.App.Config.UserviceConfig.Userinfoqueue, global.App.Config.UserviceConfig.Usererrorqueue, global.App.Config.UserviceConfig.UserServiceName, global.App.UserServiceLogger)
	go InitMQService(ch, global.App.Config.AuthServiceConfig.Authexchange, global.App.Config.AuthServiceConfig.Authinfoqueue, global.App.Config.AuthServiceConfig.Autherrorqueue, global.App.Config.AuthServiceConfig.AuthServiceName, global.App.AuthServiceLogger)
	// go authService(ch)
}

// 封装读取日志函数
func handleLogMessage(logger *zap.Logger, level string, messages <-chan amqp.Delivery) {
	for msg := range messages {
		switch level {
		case "info":
			logger.Info(string(msg.Body))
		case "error":
			logger.Error(string(msg.Body))
		default:
			logger.Info(string(msg.Body))
		}
	}
}

// // 处理用户服务相关的日志
// func userService(ch *amqp.Channel) {
// 	// 声明交换机
// 	err := ch.ExchangeDeclare(
// 		global.App.Config.UserviceConfig.Userexchange,
// 		"topic",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 		return
// 	}
// 	// Declare info and error queues
// 	infoQueue, err := ch.QueueDeclare(
// 		global.App.Config.UserviceConfig.Userinfoqueue, // name
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}

// 	errorQueue, err := ch.QueueDeclare(
// 		global.App.Config.UserviceConfig.Usererrorqueue, // name
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	// 绑定info,error队列
// 	if err := ch.QueueBind(
// 		infoQueue.Name,
// 		"info."+global.App.Config.UserviceConfig.UserServiceName,
// 		global.App.Config.UserviceConfig.Userexchange,
// 		false,
// 		nil,
// 	); err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 	}
// 	// 绑定info,error队列
// 	if err := ch.QueueBind(
// 		errorQueue.Name,
// 		"error."+global.App.Config.UserviceConfig.UserServiceName,
// 		global.App.Config.UserviceConfig.Userexchange,
// 		false,
// 		nil,
// 	); err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 	}

// 	infoMsg, err := ch.Consume(
// 		global.App.Config.UserviceConfig.Userinfoqueue,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	errorMsg, err := ch.Consume(
// 		global.App.Config.UserviceConfig.Usererrorqueue,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	go handleLogMessage(global.App.UserServiceLogger, "info", infoMsg)
// 	go handleLogMessage(global.App.UserServiceLogger, "error", errorMsg)
// 	forever := make(chan bool)
// 	<-forever
// }

// // 处理用户服务相关的日志
// func authService(ch *amqp.Channel) {
// 	// 声明交换机
// 	err := ch.ExchangeDeclare(
// 		global.App.Config.AuthServiceConfig.Authexchange,
// 		"topic",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 		return
// 	}
// 	// Declare info and error queues
// 	infoQueue, err := ch.QueueDeclare(
// 		global.App.Config.AuthServiceConfig.Authinfoqueue,
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}

// 	errorQueue, err := ch.QueueDeclare(
// 		global.App.Config.AuthServiceConfig.Autherrorqueue,
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	// 绑定info,error队列
// 	if err := ch.QueueBind(
// 		infoQueue.Name,
// 		"info."+global.App.Config.AuthServiceConfig.AuthServiceName,
// 		global.App.Config.AuthServiceConfig.Authexchange,
// 		false,
// 		nil,
// 	); err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 	}
// 	// 绑定info,error队列
// 	if err := ch.QueueBind(
// 		errorQueue.Name,
// 		"error."+global.App.Config.AuthServiceConfig.AuthServiceName,
// 		global.App.Config.AuthServiceConfig.Authexchange,
// 		false,
// 		nil,
// 	); err != nil {
// 		global.App.LogsServiceLogger.Error(err.Error())
// 	}

// 	infoMsg, err := ch.Consume(
// 		global.App.Config.AuthServiceConfig.Authinfoqueue,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	errorMsg, err := ch.Consume(
// 		global.App.Config.AuthServiceConfig.Autherrorqueue,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		global.App.LogsServiceLogger.Fatal(err.Error())
// 	}
// 	go handleLogMessage(global.App.AuthServiceLogger, "info", infoMsg)
// 	go handleLogMessage(global.App.AuthServiceLogger, "error", errorMsg)
// 	forever := make(chan bool)
// 	<-forever
// }
