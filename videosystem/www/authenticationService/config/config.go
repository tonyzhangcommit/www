package config

type Configuration struct {
	App            `mapstructure:"app"` // 这里的每一个项目都代表的是一个配置块信息，app 服务对外信息
	Redis          `mapstructure:"redis"`
	RabbitMQ       `mapstructure:"rabbitmq"`
	RabbitMQLog    `mapstructure:"rabbitmqlog"`
	LocalLogs      `mapstructure:"locallogs"`
	UserServiceApi `mapstructure:"userserviceapi"`
}
