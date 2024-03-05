package utils

import (
	"fmt"
	"log"
	"os"
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
