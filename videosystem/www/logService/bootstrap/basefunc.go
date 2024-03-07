package bootstrap

import (
	"logservice/global"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// 封装声明交换机，队列，以及监听队列函数
func InitMQService(ch *amqp.Channel, exchange, infoqueue, errorqueue, servicename string, logger *zap.Logger) {
	// 声明交换机
	err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.App.LogsServiceLogger.Error(err.Error())
		return
	}
	// Declare info and error queues
	infoQueue, err := ch.QueueDeclare(
		infoqueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}

	errorQueue, err := ch.QueueDeclare(
		errorqueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}
	// 绑定info,error队列
	if err := ch.QueueBind(
		infoQueue.Name,
		"info."+servicename,
		exchange,
		false,
		nil,
	); err != nil {
		global.App.LogsServiceLogger.Error(err.Error())
	}
	// 绑定info,error队列
	if err := ch.QueueBind(
		errorQueue.Name,
		"error."+servicename,
		exchange,
		false,
		nil,
	); err != nil {
		global.App.LogsServiceLogger.Error(err.Error())
	}

	infoMsg, err := ch.Consume(
		infoqueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}
	errorMsg, err := ch.Consume(
		errorqueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}
	go handleLogMessage(logger, "info", infoMsg)
	go handleLogMessage(logger, "error", errorMsg)
	forever := make(chan bool)
	<-forever
}
