// 将项目对外信息序列化到结构体中
package config

type App struct {
	Env     string `yaml:"env"`
	Port    string `yaml:"port"`
	AppName string `yaml:"appname"`
}
