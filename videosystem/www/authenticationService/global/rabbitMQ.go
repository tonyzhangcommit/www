package global

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// 定义rabbitMQ全局channel（本质是TCP）
type RMQ struct {
	Channel  *amqp091.Channel
	Exchange string
}

var RabbitMQ = new(RMQ)

// 日志消息结构体，日志级别，日志内容
type LogMessage struct {
	Level   string
	Message interface{}
}

func (lm *LogMessage) SendInfoToRabbitMQ() {
	var body string
	switch msg := lm.Message.(type) {
	case error:
		body = msg.Error()
	default:
		body = fmt.Sprintf("%s", msg)
	}
	err := RabbitMQ.Channel.PublishWithContext(
		context.Background(),
		RabbitMQ.Exchange,
		lm.Level+"."+App.Config.AppName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		App.LocalLogger.Error(fmt.Sprintf("日志发送失败: %s", body))
		App.LocalLogger.Error(err.Error())
	}
}

func SendLogs(level string, msg interface{}, optionalErr ...error) {
	var logMessage interface{}
	if len(optionalErr) > 0 && optionalErr[0] != nil {
		err := optionalErr[0]
		logMessage = fmt.Sprintf("%v - 错误详情: %v", msg, err)
	} else {
		logMessage = msg
	}
	logmsg := LogMessage{
		Level:   level,
		Message: logMessage,
	}
	logmsg.SendInfoToRabbitMQ()
}
