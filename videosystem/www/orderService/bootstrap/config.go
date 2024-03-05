package bootstrap

import (
	"fmt"
	"userservice/global"

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
	// 测试读取文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件错误：%v\n", err))
	}
	// 绑定配置文件内容到结构体中
	if err := v.Unmarshal(&global.App.Config); err != nil {
		panic(fmt.Sprintf("解析文件失败：%v\n", err))
	}
	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(&global.App.Config); err != nil {
			panic(fmt.Sprintf("更改配置文件导致解析文件失败：%v\n", err))
		}
	})
	global.App.Viper = v
}
