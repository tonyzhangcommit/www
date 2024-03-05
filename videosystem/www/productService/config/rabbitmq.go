package config

/*
	rabbitMQ 配置
*/

type RabbitMQ struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Vhost          string `mapstructure:"vhost"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	ExchangeName   string `mapstructure:"exchangename"`
	Userinfoqueue  string `mapstructure:"userinfoqueue"`
	Usererrorqueue string `mapstructure:"usererrorqueue"`
}
