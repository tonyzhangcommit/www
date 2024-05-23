package config

type Jwt struct {
	Secretkey               string `mapstructure:"secretkey"`               // 密钥
	Jwttil                  int64  `mapstructure:"jwttil"`                  // token有效期
	JwtBlacklistGracePeriod int64  `mapstructure:"jwtBlacklistGracePeriod"` // 黑名单过期时间（预防在退出登录时，其他并发请求失败的场景,单位秒）
}
