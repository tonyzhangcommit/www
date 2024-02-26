package global

import (
	"logservice/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Application struct {
	Config            config.Config
	Viper             *viper.Viper
	UserServiceLogger *zap.Logger
	AuthServiceLogger *zap.Logger
	LogsServiceLogger *zap.Logger
}

var App = new(Application)
