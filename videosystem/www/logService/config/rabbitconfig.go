package config

type RabbitConfigServer struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Vhost    string `mapstructure:"vhost"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type UserServiceConfig struct {
	UserServiceName string `mapstructure:"serivcename"`
	Userexchange    string `mapstructure:"userexchange"`
	Userinfoqueue   string `mapstructure:"userinfoqueue"`
	Usererrorqueue  string `mapstructure:"usererrorqueue"`
}

type AuthServiceConfig struct {
	AuthServiceName string `mapstructure:"serivcename"`
	Authexchange    string `mapstructure:"authexchange"`
	Authinfoqueue   string `mapstructure:"authinfoqueue"`
	Autherrorqueue  string `mapstructure:"autherrorqueue"`
}
