package utils

import (
	"math/rand"
	"strings"
	"time"
)

// 生成随机识别码
func GenerateRCode(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var str strings.Builder
	for i := 0; i < length; i++ {
		index := rnd.Intn(len(charset))
		str.WriteByte(charset[index])
	}
	return str.String()
}

// 生成随机验证码
func GenerateNumberCode(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	charset := "0123456789"
	var str strings.Builder
	for i := 0; i < length; i++ {
		index := rnd.Intn(len(charset))
		str.WriteByte(charset[index])
	}
	return str.String()
}
