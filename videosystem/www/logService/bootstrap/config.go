package bootstrap

import (
	"fmt"
	"logservice/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/*
	初始化配置文件
*/

func InitializeConfig() {

	config_file := "config.yaml"
	v := viper.New()
	v.SetConfigFile(config_file)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败，err:%v\n", err))
	}
	if err := v.Unmarshal(&global.App.Config); err != nil {
		panic(fmt.Sprintf("序列化失败，err:%v\n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(&global.App.Config); err != nil {
			panic(fmt.Sprintf("序列化失败，err:%v\n", err))
		}
	})
	global.App.Viper = v
}
