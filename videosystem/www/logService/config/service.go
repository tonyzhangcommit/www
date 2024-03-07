package config

/*
这里是各个子服务配置信息
BaseServiceLog: 基础日志配置信息
LocalServiceLog: 当前日志服务
UserServiceLog: 用户管理日志服务
*/
type BaseServiceLog struct {
	RootDir       string `mapstructure:"root_dir"`
	ShowLine      bool   `mapstructure:"show_line"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
	DefaultErrorL string `mapstructure:"defaultErrorL"`
}

type LocalServiceLog struct {
	ServiceName string `mapstructure:"servicename"`
	RootDir     string `mapstructure:"root_dir"`
	Level       string `mapstructure:"level"`
	JsonFormat  bool   `mapstructure:"jsonformat"`
	Info        string `mapstructure:"info"`
	Error       string `mapstructure:"error"`
}

// userlog本地配置
type UserServiceLog struct {
	ServiceName string `mapstructure:"servicename"`
	RootDir     string `mapstructure:"root_dir"`
	Level       string `mapstructure:"level"`
	JsonFormat  bool   `mapstructure:"jsonformat"`
	Info        string `mapstructure:"info"`
	Error       string `mapstructure:"error"`
}

type AuthServiceLog struct {
	ServiceName string `mapstructure:"servicename"`
	RootDir     string `mapstructure:"root_dir"`
	Level       string `mapstructure:"level"`
	JsonFormat  bool   `mapstructure:"jsonformat"`
	Info        string `mapstructure:"info"`
	Error       string `mapstructure:"error"`
}

type ProductServiceLog struct {
	ServiceName string `mapstructure:"servicename"`
	RootDir     string `mapstructure:"root_dir"`
	Level       string `mapstructure:"level"`
	JsonFormat  bool   `mapstructure:"jsonformat"`
	Info        string `mapstructure:"info"`
	Error       string `mapstructure:"error"`
}

type OrderServiceLog struct {
	ServiceName string `mapstructure:"servicename"`
	RootDir     string `mapstructure:"root_dir"`
	Level       string `mapstructure:"level"`
	JsonFormat  bool   `mapstructure:"jsonformat"`
	Info        string `mapstructure:"info"`
	Error       string `mapstructure:"error"`
}
