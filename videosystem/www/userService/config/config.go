package config

type Configuration struct {
	UserService UserService `mapstructure:"userservice"`
	Database    Database    `mapstructure:"databases"`
	Redis       Redis       `mapstructure:"redis"`
	Roles       Role        `mapstructure:"roles"`
	Permission  Permission  `mapstructure:"permission"`
	RabbitMQ    RabbitMQ    `mapstructure:"rabbitmq"`
	LocalLogs   LocalLogs   `mapstructure:"locallogs"`
}
