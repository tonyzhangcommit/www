package utils

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

/*
自定义验证规则，完善validator规则
*/

// 验证手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if mobile == "" {
		return true
	}
	ok, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	return ok
}

// 自定义密码规则  包含大小写字母，数字和特殊符号，长度6-12
func ValidatePassword(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 6 || len(pwd) > 12 {
		return false
	}

	// 验证小写字母
	lowerCase := regexp.MustCompile(`[a-z]`)
	if !lowerCase.MatchString(pwd) {
		return false
	}

	// 验证大写字母
	upperCase := regexp.MustCompile(`[A-Z]`)
	if !upperCase.MatchString(pwd) {
		return false
	}

	// 验证数字
	number := regexp.MustCompile(`\d`)
	if !number.MatchString(pwd) {
		return false
	}

	// 验证特殊字符
	specialChar := regexp.MustCompile(`[.)(*&^%$#@!]`)
	return specialChar.MatchString(pwd)
}

// 验证代理邀请码
func ValidateAgentCode(fl validator.FieldLevel) bool {
	agentCode := fl.Field().String()
	if len(agentCode) == 0 {
		return true
	}
	regexpPattern := `^[A-Za-z0-9]{6,8}$`
	re, _ := regexp.Compile(regexpPattern)
	return re.Match([]byte(agentCode))
}

// 自定义验证用户名
func ValidateUserName(fl validator.FieldLevel) bool {
	userName := fl.Field().String()
	if userName == "" {
		return true
	}
	regexpPattern := `^[\p{L}\p{N}]{2,16}$`
	re, _ := regexp.Compile(regexpPattern)
	return re.Match([]byte(userName))
}

// 验证身份证号合法性
func ValidateIDCard(fl validator.FieldLevel) bool {
	idcard := fl.Field().String()
	if idcard == "" {
		return true
	}
	if len(idcard) != 18 {
		return false
	}
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"

	sum := 0
	for i := 0; i < 17; i++ {
		num, err := strconv.Atoi(string(idcard[i]))
		if err != nil {
			return false // 前17位中有非数字字符
		}
		sum += num * weights[i]
	}
	checkIndex := sum % 11
	checkChar := checkCodes[checkIndex]

	return checkChar == idcard[17]
}

// 自定义验证器email
func CustomEmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true
	}
	// 使用 validator 内置的 Email 校验
	return validator.New().Var(email, "email") == nil
}

func OnlyNum(fl validator.FieldLevel) bool {
	return true
}
