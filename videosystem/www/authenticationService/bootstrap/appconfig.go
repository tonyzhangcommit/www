package bootstrap

import (
	"auth/global"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 初始化服务发现配置
func InitializeConfig() {
	// 这里以后可以增加配置文件地址，从系统变量获取路径
	config := "config.yaml"
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 绑定配置信息到结构体
	if err := v.Unmarshal(&global.App.Config); err != nil {
		panic(fmt.Sprintf("配置文件动态绑定失败: %v", err))
	}
	// 监听配置文件变化
	v.WatchConfig()
	// 配置文件更改后重新加载
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		if err := v.Unmarshal(&global.App.Config); err != nil {
			fmt.Println("配置文件动态绑定失败", err)
		}
	})

	global.App.ConfigViper = v
}
