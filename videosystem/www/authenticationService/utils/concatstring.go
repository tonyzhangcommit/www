package utils

import "strings"

func JoinStrings(stringsToJoin ...string) string {
	var builder strings.Builder // 创建一个strings.Builder实例

	// 遍历所有传入的字符串参数
	for _, str := range stringsToJoin {
		builder.WriteString(str) // 使用WriteString方法将字符串添加到builder
	}

	return builder.String() // 将拼接后的字符串结果返回
}
