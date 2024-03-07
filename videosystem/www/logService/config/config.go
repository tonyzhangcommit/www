package config

/*
	config 包用于保存配置文件中分块信息
*/

type Config struct {
	ServiceInfo          `mapstructure:"app"`
	RabbitMQ             `mapstructure:"rabbitmq"`
	BaseServiceLog       `mapstructure:"baselogconfig"`
	LocalServiceLog      `mapstructure:"localserverlog"`
	UserServiceLog       `mapstructure:"userservicelog"`
	UserServiceConfig    `mapstructure:"usermanagerconfig"`
	AuthServiceLog       `mapstructure:"authsvicelog"`
	AuthServiceConfig    `mapstructure:"authserverconfig"`
	ProductServiceLog    `mapstructure:"productservicelog"`
	ProductServiceConfig `mapstructure:"productserverconfig"`
	OrderServiceLog      `mapstructure:"orderservicelog"`
	OrderServiceConfig   `mapstructure:"orderserverconfig"`
}
