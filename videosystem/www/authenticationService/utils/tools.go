package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
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
