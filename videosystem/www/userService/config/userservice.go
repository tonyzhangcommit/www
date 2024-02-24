package config

/*
	服务对外配置
*/

type UserService struct {
	Env     string `mapstructure:"env"`
	Port    string `mapstructure:"port"`
	AppName string `mapstructure:"appname"`
}
