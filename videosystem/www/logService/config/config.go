package config

/*
	config 包用于保存配置文件中分块信息
*/

type Config struct {
	ServiceInfo     ServiceInfo        `mapstructure:"app"`
	RabbitMQ        RabbitConfigServer `mapstructure:"rabbitmq"`
	BaseServiceLog  BaseServiceLog     `mapstructure:"baselogconfig"`
	LocalServiceLog LocalServiceLog    `mapstructure:"localserverlog"`
	UserServiceLog  UserServiceLog     `mapstructure:"userservicelog"`
	UserviceConfig  UserServiceConfig  `mapstructure:"usermanagerconfig"`
}
