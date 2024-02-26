package config

// 用户管理服务接口配置
type UserServiceApi struct {
	Name      string     `mapstructure:"name"`
	BaseUrl   string     `mapstructure:"baseurl"`
	Timeout   int        `mapstructure:"timeout"`
	ClientUrl Userclient `mapstructure:"client"`
	AdminUrl  Useradmin  `mapstructure:"admin"`
}

// 客户端接口配置
type Userclient struct {
	Login        string `mapstructure:"login"`
	Register     string `mapstructure:"register"`
	Getvirifcode string `mapstructure:"getvirifcode"`
}

// 管理端接口配置
type Useradmin struct {
}
