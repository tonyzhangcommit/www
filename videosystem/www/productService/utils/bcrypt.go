package utils

import (
	"userservice/global"

	"golang.org/x/crypto/bcrypt"
)

/*
	密码加密解密功能
*/
// 加密
func BcryptMake(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		global.SendLogs("error", "加密失败", err)
		return ""
	}
	return string(hash)
}

// 验证
func BcryptMakeCheck(pwd []byte, hashePwd string) bool {
	byteHash := []byte(hashePwd)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	return err == nil
}
