package request

/*
	基础功能相关请求参数格式定义
	功能为：登录,注册,找回密码,编辑个人信息,退出登录,注销
*/

// 注册,用户名（唯一）,密码,电话（唯一）,注册码（可选）
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
		"name.username":       "用户名格式错误,长度应该在2-16位,不包含特殊字符",
		"password.required":   "密码不能为空",
		"password.password":   "密码格式错误,应同时包含大小写字母,数字,特殊字符,长度6-12",
		"phonenum.required":   "手机号不能为空",
		"phonenum.mobile":     "手机号格式错误",
		"varificode.required": "请输入验证码",
		"agentcode.agentcode": "邀请码格式错误",
	}
}

// 登录  用户名 手机号 二选一, 密码必选, 验证码
type LoginNP struct {
	Name          string `form:"name" json:"name" binding:"required,username"`
	Password      string `form:"password" json:"password" binding:"required"`
	RealIPAddress string `form:"realipaddress" json:"realipaddress"` // 等部署功能后,可以统计登录信息的IP,便于分析
}

func (loginnp LoginNP) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "姓名不能为空",
		"name.username":     "用户名格式错误,长度应该在2-16位,不包含特殊字符",
		"password.required": "密码不能为空",
		"password.password": "密码格式错误,应同时包含大小写字母,数字,特殊字符,长度6-12",
	}
}

// 登录  用户名 手机号 二选一, 密码必选, 验证码
type LoginPVC struct {
	Phonenumber      string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
	VerificationCode string `form:"varificode" json:"varificode" binding:"required"`
	RealIPAddress    string `form:"realipaddress" json:"realipaddress"` // 等部署功能后,可以统计登录信息的IP,便于分析
}

func (Loginpvc LoginPVC) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"phonenum.required":   "手机号不能为空",
		"varificode.required": "验证码不能为空",
		"phonenum.mobile":     "手机号格式错误",
	}
}

// 登出
type LogOut struct {
	UserId uint64 `form:"userid" json:"userid" binding:"required"`
}

func (logout LogOut) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"userid.required": "用户ID不能为空",
	}
}

// 获取个人信息
type GetPersonInfo struct {
	Phonenumber string `form:"phonenum" json:"phonenum" binding:"mobile"`
	Uid         string `form:"uid" json:"uid" `
}

func (getpersoninfo GetPersonInfo) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.username":   "用户名格式错误,长度应该在2-16位,不包含特殊字符",
		"phonenum.mobile": "手机号格式错误",
	}
}

// 编辑个人信息
// dev 暂时不用 binding:"idcard" 身份验证
type InproveInfo struct {
	UserID         uint   `form:"uid" json:"uid" binding:"required"`
	Address        string `form:"address" json:"address" `
	Sex            uint   `form:"sex" json:"sex"`
	Identification string `form:"identification" json:"identification"`
	Email          string `form:"email" json:"email" binding:"customemail"`
	Preferences    string `form:"Preferences" json:"Preferences" `
}

func (inproveinfo InproveInfo) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"uid.required": "用户id不能为空",
		// "uid.alphanum":          "id格式错误",
		"identification.idcard": "身份证格式错误",
		"email.email":           "邮件格式错误",
	}
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

// 获取用户有效角色
type GetVipType struct {
	Uid uint `form:"uid" json:"uid" binding:"required"`
}

func (u GetVipType) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"uid.required": "用户id不能为空",
	}
}

// 创建管理员
type AdminCreateManager struct {
	Uid         uint   `form:"uid" json:"uid" binding:"required"`
	UserName    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	Phonenumber string `form:"phonenum" json:"phonenum" binding:"required,mobile"`
}

func (a *AdminCreateManager) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "姓名不能为空",
		"password.required": "密码不能为空",
		"phonenum.required": "手机号不能为空",
		"phonenum.mobile":   "手机号格式错误",
	}
}

// 删除用户
type DeleteUser struct {
	Uid       uint `form:"uid" json:"uid" binding:"required"`
	DeleteUid uint `form:"deleteuid" json:"deleteuid"  binding:"required"`
}

func (d *DeleteUser) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"uid.required":       "uid不能为空",
		"deleteuid.required": "deleteuid不能为空",
	}
}

// 管理员获取用户列表
type AdminGetUserList struct {
	Uid            uint  `form:"uid" json:"uid" binding:"required"`
	IsFirstRequest bool  `form:"isfirstrequest" json:"isfirstrequest"`
	PageSize       int64 `form:"pagesize" json:"pagesize"`
	CurrentPage    int64 `form:"currentpage" json:"currentpage"`
}

func (a *AdminGetUserList) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"uid.required": "uid不能为空",
	}
}

// 管理员获取单个用户信息
type AdminGetUserinfo struct {
	Uid       uint `form:"uid" json:"uid" binding:"required"`
	TargetUid uint `form:"targetuid" json:"targetuid" binding:"required"`
	Action    bool `form:"action" json:"action"` // 封禁接口复用，0封禁，1表示解禁
}

func (a *AdminGetUserinfo) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"uid.required":       "uid不能为空",
		"targetuid.required": "targetuid不能为空",
	}
}
