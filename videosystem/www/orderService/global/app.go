package global

import (
	"context"
	"fmt"
	"sync"
	"time"
	"userservice/config"

	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
全局变量结构体
*/
type Application struct {
	Viper       *viper.Viper
	LocalLogger zap.Logger
	Config      config.Configuration
	DB          *gorm.DB
	Redis       *redis.Client
}

var App = new(Application)

type RMQ struct {
	Channel  *amqp.Channel
	Exchange string
	Conn     *amqp.Connection
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
		lm.Level+"."+App.Config.App.AppName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		App.LocalLogger.Error(fmt.Sprintf("日志发送失败: %s\n", body))
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

/*
全局本地缓存结构体
*/
type Cache struct {
	store sync.Map
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	if _, exist := c.Get(key); exist {
		c.store.Delete(key)
	}
	c.store.Store(key, value)
	if duration > 0 {
		go func() {
			<-time.After(duration)
			c.store.Delete(key)
		}()
	}
}
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Load(key)
}

var Store = new(Cache)
