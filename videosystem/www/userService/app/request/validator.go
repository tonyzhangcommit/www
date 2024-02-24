package request

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

/*
封装获取验证信息功能函数
*/

// 定义验证错误信息类型
type ValidatorMessages map[string]string

type Validator interface {
	GetMessages() ValidatorMessages
}

func GetErrorMsg(request interface{}, err error) string {
	if _, ok := err.(validator.ValidationErrors); ok {
		_, isValidator := request.(Validator)
		var errorslice []string
		for _, v := range err.(validator.ValidationErrors) {
			if isValidator {
				if message, exist := request.(Validator).GetMessages()[v.Field()+"."+v.Tag()]; exist {
					errorslice = append(errorslice, message)
				}
			} else {
				errorslice = append(errorslice, v.Error())
			}
		}
		return strings.Join(errorslice, " - ")
	}
	return "参数错误"
}
