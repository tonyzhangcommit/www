package config

/*
	缓存配置
*/

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
