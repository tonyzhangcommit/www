package bootstrap

import (
	"userservice/global"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
	初始化rabbitMQ
	连接并创建channel
*/

func InitializeRabbitMQ() {
	// 首先初始化本地 local 配置
	InitLocalLogger()
	rabbitUsername := global.App.Config.RabbitMQ.User
	rabbitPassword := global.App.Config.RabbitMQ.Password
	rabbitPort := global.App.Config.RabbitMQ.Port
	rabbitHost := global.App.Config.RabbitMQ.Host
	rabbitVhost := global.App.Config.RabbitMQ.Vhost
	amqpStr := "amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/" + rabbitVhost
	conn, err := amqp.Dial(amqpStr)
	if err != nil {
		// 这里应该增加预警提示，比如发邮件，短信等通知
		global.App.LocalLogger.Error(err.Error())
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		global.App.LocalLogger.Error(err.Error())
		return
	}
	global.RabbitMQ.Channel = ch
	global.RabbitMQ.Exchange = global.App.Config.RabbitMQ.ExchangeName
}

/*
	定义秒杀服务相关队列和绑定相关交换机
	订单请求MQ
	处理结果MQ
	生成订单MQ
	死信MQ
	执行流程：
	请求->订单请求MQ-------->
*/


