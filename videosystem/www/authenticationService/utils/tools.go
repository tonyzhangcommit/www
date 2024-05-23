package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"time"
)

// 创建文件夹
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

// 生成密钥
func GenerateSecretKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	// Use Base64 encoding for ease of storage and use
	secretKey := base64.URLEncoding.EncodeToString(key)
	return secretKey, nil
}

// SHA256 加密处理字符串
func Sha256(str []byte) string {
	h := sha256.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// 创建锁标识
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65 // 生成 65 - 65 + 26 之间的随机整数，对应ascii为大写字母A-Z
		bytes[i] = byte(b)
	}
	return string(bytes)
}
