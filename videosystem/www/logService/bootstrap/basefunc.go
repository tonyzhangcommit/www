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

// 初始化死信交换机和绑定死信队列
func InitDeadExchange(ch *amqp.Channel, exchange, queuename, routingkey string) {
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
		global.App.LogsServiceLogger.Fatal(err.Error())
		return
	}
	// 声明队列
	args := amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-message-ttl":             36000000,    // 消息存活时间 (ms)
		"x-max-length":              10000,       // 队列最大长度
		"x-max-length-bytes":        5000000,     // 队列最大占用空间 (byte)
		"x-overflow":                "drop-head", // 队列溢出行为
		"x-dead-letter-routing-key": routingkey,  // 死信路由键
	}

	queue, err := ch.QueueDeclare(
		queuename, // 队列名称
		true,      // 持久化
		false,     // 非排他性
		false,     // 自动删除
		false,     // 不等待
		args,      // 详细参数
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}
	// 将死信队列绑定到死信交换机
	err = ch.QueueBind(
		queue.Name, // 死信队列名称
		routingkey, // 死信路由键
		exchange,   // 死信交换机名称
		false,
		nil,
	)
	if err != nil {
		global.App.LogsServiceLogger.Fatal(err.Error())
	}

}
