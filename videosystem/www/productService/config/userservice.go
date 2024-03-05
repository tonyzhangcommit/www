package config

/*
	服务对外配置
*/

type App struct {
	Env     string `mapstructure:"env"`
	Port    string `mapstructure:"port"`
	AppName string `mapstructure:"appname"`
}
