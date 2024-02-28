package config

type Jwt struct {
	Secretkey string `mapstructure:"secretkey"` // 密钥
	Jwttil    int64  `mapstructure:"jwttil"`    // token有效期
}
