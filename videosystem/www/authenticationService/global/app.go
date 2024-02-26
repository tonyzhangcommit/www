package global

import (
	"auth/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Application struct {
	ConfigViper *viper.Viper // 这个配置项的作用是，方便项目中动态进行配置文件增，改操作
	Config      config.Configuration
	Redis       *redis.Client
	LocalLogger *zap.Logger
}

var App = new(Application)
