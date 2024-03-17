package config

/*
	rabbitMQ 配置,这里除基础配置（host port vhost等）以外，以下是队列作用：
	1. Userinfoqueue: 日志info级别传送队列
	2. Usererrorqueue: 日志error级别传送队列
	3. FlashEventReqQueue: 秒杀活动请求削峰队列
	4. FlashEventResQueue: 秒杀活动完成请求队列
	5. FlashEventDeadQueue: 秒杀活动死信队列
*/

type RabbitMQ struct {
	Host                     string `mapstructure:"host"`
	Port                     string `mapstructure:"port"`
	Vhost                    string `mapstructure:"vhost"`
	User                     string `mapstructure:"user"`
	Password                 string `mapstructure:"password"`
	ExchangeName             string `mapstructure:"exchangename"`
	Userinfoqueue            string `mapstructure:"userinfoqueue"`
	Usererrorqueue           string `mapstructure:"usererrorqueue"`
	FlashEventReqQueue       string `mapstructure:"flasheventreqqueue"`
	FlashEventReqBackupQueue string `mapstructure:"flasheventreqbackupqueue"`
	FlashEventResQueue       string `mapstructure:"flasheventresqueue"`
}
