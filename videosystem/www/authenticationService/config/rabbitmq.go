package config

/*
	消息队列配置
*/

// 系统配置
type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Vhost    string `mapstructure:"vhost"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// 日志交换机，队列配置
// 这里只是最基础的配置，后续可根据功能要求增加配置
type RabbitMQLog struct {
	ExchangeName   string `mapstructure:"exchangename"`
	Authinfoqueue  string `mapstructure:"userinfoqueue"`
	Autherrorqueue string `mapstructure:"usererrorqueue"`
}
