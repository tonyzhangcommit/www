package config

// 用户管理服务接口配置
type UserServiceApi struct {
	Name           string     `mapstructure:"name"`
	BaseUrl        string     `mapstructure:"baseurl"`
	Timeout        int        `mapstructure:"timeout"`
	ClientPath     string     `mapstructure:"clientpath"`
	AdminPath      string     `mapstructure:"adminpath"`
	FlashEventPath string     `mapstructure:"flashpath"`
	ClientUrl      Userclient `mapstructure:"client"`
	AdminUrl       Useradmin  `mapstructure:"admin"`
}

// 客户端接口配置

type Userclient struct {
	Login        string `mapstructure:"login"`
	Register     string `mapstructure:"register"`
	Getverifcode string `mapstructure:"getverifcode"`
	GetuserInfo  string `mapstructure:"getuserinfo"`
	InproveInfo  string `mapstructure:"inproveinfo"`
	Getuvip      string `mapstructure:"getuvip"`
}

// 管理端接口配置
type Useradmin struct {
	Login   string `mapstructure:"login"`
	Preheat string `mapstructure:"preheat"`
}

// 商品管理服务接口配置
type ProductServiceApi struct {
	Name           string `mapstructure:"name"`
	BaseUrl        string `mapstructure:"baseurl"`
	Timeout        int    `mapstructure:"timeout"`
	FlashGetFEinfo string `mapstructure:"flashgetfeinfo"`
	FlashGetEventP string `mapstructure:"flashgeteventp"`
}

// 订单管理服务接口配置
type OrderServiceApi struct {
	Name           string `mapstructure:"name"`
	BaseUrl        string `mapstructure:"baseurl"`
	Timeout        int    `mapstructure:"timeout"`
	TakeFalshOrder string `mapstructure:"takefalshorder"`
}
