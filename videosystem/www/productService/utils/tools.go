package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func CreateDir(path string) {
	fmt.Println(path)
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

func JoinStrings(stringsToJoin ...string) string {
	var builder strings.Builder // 创建一个strings.Builder实例

	// 遍历所有传入的字符串参数
	for _, str := range stringsToJoin {
		builder.WriteString(str) // 使用WriteString方法将字符串添加到builder
	}

	return builder.String() // 将拼接后的字符串结果返回
}
