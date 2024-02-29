package bootstrap

import (
	"auth/utils"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

/*
定制Validator 属性

*/

func InitializeValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 说明当前引擎使用的validator,开始注册
		_ = v.RegisterValidation("mobile", utils.ValidateMobile)
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
