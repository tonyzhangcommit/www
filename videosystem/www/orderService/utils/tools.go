package utils

import (
	"log"
	"os"
	"strings"
	"userservice/global"

	"github.com/google/uuid"
)

// 创建路径
func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755) // 0755 是文件夹的权限设置
		if err != nil {
			log.Fatalf("Error creating directory: %s", err)
		}
		log.Printf("Directory created: %s", path)
	} else {
		log.Printf("Directory already exists: %s", path)
	}
}

// 生成UUID
func GenerateUniqueID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		global.SendLogs("error", "生成uuid失败", err)
		return ""
	}
	return id.String()
}

// 拼接字符串
func JoinStrings(stringsToJoin ...string) string {
	var builder strings.Builder // 创建一个strings.Builder实例

	// 遍历所有传入的字符串参数
	for _, str := range stringsToJoin {
		builder.WriteString(str) // 使用WriteString方法将字符串添加到builder
	}

	return builder.String() // 将拼接后的字符串结果返回
}
