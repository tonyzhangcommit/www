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
	global.RabbitMQ.Conn = conn
	global.RabbitMQ.Channel = ch
	global.RabbitMQ.Exchange = global.App.Config.RabbitMQ.ExchangeName
	// 初始化秒杀活动reqMQ
	genFlashEventReqMQ(ch)
	// 初始化秒杀活动reqBackupMQ
	genFlashEventReqBackupMQ(ch)
	// 定义订单请求处理结果MQ
	genFlashEventResMQ(ch)
	// 定义待支付订单MQ(主要服务于支付服务)
	genFlashEventPayOrderMQ(ch)
}

/*
	定义秒杀服务相关队列和绑定相关交换机
	订单请求MQ
	处理结果MQ
	死信MQ
*/

// 定义订单请求削峰MQ
func genFlashEventReqMQ(ch *amqp.Channel) {
	args := amqp.Table{
		"x-message-ttl": 6000000,     // 消息存活时间 (ms)
		"x-max-length":  150000,      // 队列最大长度
		"x-overflow":    "drop-head", // 队列溢出行为
	}
	queue, err := ch.QueueDeclare(
		global.App.Config.RabbitMQ.FlashEventReqQueue,
		true,  // 持久化
		false, // 非排他性
		false, // 自动删除
		false, // 不等待
		args,
	)
	if err != nil {
		global.SendLogs("error", "创建秒杀订单消息队列失败", err)
		return
	}
	// 绑定交换机
	if err := ch.QueueBind(
		queue.Name,
		"flashevent.req.row",
		global.App.Config.RabbitMQ.ExchangeName,
		false,
		nil,
	); err != nil {
		global.SendLogs("error", "绑定秒杀队列req失败", err)
		return
	}
}

// 定义订单请求削峰MQ(备份队列)
func genFlashEventReqBackupMQ(ch *amqp.Channel) {
	args := amqp.Table{
		"x-message-ttl": 6000000,     // 消息存活时间 (ms)
		"x-max-length":  100000,      // 队列最大长度
		"x-overflow":    "drop-head", // 队列溢出行为
	}
	queue, err := ch.QueueDeclare(
		global.App.Config.RabbitMQ.FlashEventReqBackupQueue,
		true,  // 持久化
		false, // 非排他性
		false, // 自动删除
		false, // 不等待
		args,
	)
	if err != nil {
		global.SendLogs("error", "创建秒杀订单消息备用队列失败", err)
		return
	}
	// 绑定交换机
	if err := ch.QueueBind(
		queue.Name,
		"flashevent.req.backup",
		global.App.Config.RabbitMQ.ExchangeName,
		false,
		nil,
	); err != nil {
		global.SendLogs("error", "绑定秒杀队列req backup失败", err)
		return
	}
}

// 定义订单请求处理结果MQ
func genFlashEventResMQ(ch *amqp.Channel) {
	args := amqp.Table{
		"x-message-ttl": 12000000,    // 消息存活时间 (ms)
		"x-max-length":  300000,      // 队列最大长度
		"x-overflow":    "drop-head", // 队列溢出行为
	}
	queue, err := ch.QueueDeclare(
		global.App.Config.RabbitMQ.FlashEventResQueue,
		true,  // 持久化
		false, // 非排他性
		false, // 自动删除
		false, // 不等待
		args,
	)
	if err != nil {
		global.SendLogs("error", "创建秒杀订单res消息队列失败", err)
		return
	}
	// 绑定交换机
	if err := ch.QueueBind(
		queue.Name,
		"flashevent.res.row",
		global.App.Config.RabbitMQ.ExchangeName,
		false,
		nil,
	); err != nil {
		global.SendLogs("error", "绑定秒杀队列res失败", err)
		return
	}
}

// 定义生成订单后待支付MQ
func genFlashEventPayOrderMQ(ch *amqp.Channel) {
	args := amqp.Table{
		"x-message-ttl": 12000000,    // 消息存活时间 (ms)
		"x-max-length":  300000,      // 队列最大长度 （可优化为商品抢购的最大数量）
		"x-overflow":    "drop-head", // 队列溢出行为
	}
	queue, err := ch.QueueDeclare(
		global.App.Config.RabbitMQ.FlashEventReadyPayQueue,
		true,  // 持久化
		false, // 非排他性
		false, // 自动删除
		false, // 不等待
		args,
	)
	if err != nil {
		global.SendLogs("error", "创建秒杀活动支付订单MQ失败", err)
		return
	}
	// 绑定交换机
	if err := ch.QueueBind(
		queue.Name,
		"flashevent.pay.order.queue",
		global.App.Config.RabbitMQ.ExchangeName,
		false,
		nil,
	); err != nil {
		global.SendLogs("error", "绑定秒杀活动支付订单MQ失败", err)
		return
	}
}
