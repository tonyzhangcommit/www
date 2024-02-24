package config

/*
	数据库配置
*/

type Database struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// 这里是初始化role permission 配置项

type Role struct {
	NameList []string `mapstructure:"namelist"`
}

type Permission struct {
	SuperAdmin   []string `mapstructure:"superAdmin"`
	Admin        []string `mapstructure:"admin"`
	RegularUser  []string `mapstructure:"regularUser"`
	MonthlyVip   []string `mapstructure:"monthlyVip"`
	QuarterlyVip []string `mapstructure:"quarterlyVip"`
	AnnualVip    []string `mapstructure:"annualVip"`
}
