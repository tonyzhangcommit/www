package config

type WebSocket struct {
	Maxdur            int  `mapstructure:"maxdur"`
	HandShakeTimeout  int  `mapstructure:"handshaketimeout"`
	EnableCompression bool `mapstructure:"enablecompression"`
	JwtSecretkey      string `mapstructure:"jwtsecretkey"`
}
