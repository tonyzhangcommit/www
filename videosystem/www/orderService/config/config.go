package config

type Configuration struct {
	App       `mapstructure:"app"`
	Database  `mapstructure:"databases"`
	Redis     `mapstructure:"redis"`
	RabbitMQ  `mapstructure:"rabbitmq"`
	LocalLogs `mapstructure:"locallogs"`
}
