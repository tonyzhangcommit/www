package config

// mapstructure
type ServiceInfo struct {
	Env         string `mapstructure:"env"`
	ServiceName string `mapstructure:"servicename"`
}
