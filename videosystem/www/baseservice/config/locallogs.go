package config

/*
	本地日志设置
*/

type LocalLogs struct {
	Dir          string `mapstructure:"dir"`
	Logfilename  string `mapstructure:"logfilename"`
	Max_backups  int    `mapstructure:"max_backups"`
	Max_size     int    `mapstructure:"max_size"`
	IsJson       bool   `mapstructure:"isJson"`
	Compress     bool   `mapstructure:"compress"`
	DefaultLevel string `mapstructure:"defaultLevel"`
}
