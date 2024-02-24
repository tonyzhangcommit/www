package utils

// 消息结构体
type LogMessage struct {
	ServiceName string `json:"serviceName"` // 服务名称
	LogLevel    string `json:"logLevel"`
}
