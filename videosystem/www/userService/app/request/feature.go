package request

/*
	基础功能相关请求参数格式定义
	功能为：登录，注册，找回密码，编辑个人信息，退出登录，注销
*/

// 注册，用户名（唯一），密码，电话（唯一），注册码（可选）
type Resister struct {
	Name        string `form:"name" json:"name" binding:"required,username"`
	Password    string `form:"password" json:"password" binding:"required,password"`
	Phonenumber string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
	VarifiCode  string `form:"varificode" json:"varificode" binding:"required"`
	AgentCode   string `form:"agentcode" json:"agentcode" binding:"agentcode"`
}

func (register Resister) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":       "姓名不能为空",
		"name.username":       "用户名格式错误，长度应该在2-16位，不包含特殊字符",
		"password.required":   "密码不能为空",
		"password.password":   "密码格式错误，应同时包含大小写字母，数字，特殊字符，长度6-12",
		"phonenum.required":   "手机号不能为空",
		"phonenum.mobile":     "手机号格式错误",
		"varificode.required": "请输入验证码",
		"agentcode.agentcode": "邀请码格式错误",
	}
}

// 登录  用户名 手机号 二选一， 密码必选, 验证码
type LoginNP struct {
	Name     string `form:"name" json:"name" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func (loginnp LoginNP) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "姓名不能为空",
		"name.username":     "用户名格式错误，长度应该在2-16位，不包含特殊字符",
		"password.required": "密码不能为空",
		"password.password": "密码格式错误，应同时包含大小写字母，数字，特殊字符，长度6-12",
	}
}

// 登录  用户名 手机号 二选一， 密码必选, 验证码
type LoginPVC struct {
	Phonenumber      string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
	VerificationCode string `form:"varificode" json:"varificode" binding:"required"`
}

func (Loginpvc LoginPVC) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"phonenum.required":   "手机号不能为空",
		"varificode.required": "验证码不能为空",
		"phonenum.mobile":     "手机号格式错误",
	}
}

// 编辑个人信息
type InproveInfo struct {
	UserID         uint   `form:"uid" json:"uid" binding:"required"`
	Address        string `form:"address" json:"address" `
	Sex            uint   `form:"sex" json:"sex" validate:"oneof=0 1"`
	Identification string `form:"identification" json:"identification" `
	Email          string `form:"email" json:"email" `
}

// 获取验证码
type GetVirifCode struct {
	Phonenumber string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
}

func (getvirifcode GetVirifCode) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"phonenum.required": "手机号不能为空",
		"phonenum.mobile":   "手机号格式错误",
	}
}
