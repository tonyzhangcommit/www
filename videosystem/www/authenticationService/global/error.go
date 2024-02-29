package global

/*
	封装错误码
*/

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	LocalServiceError CustomError // 认证服务发生错误
	UserServiceError  CustomError // 用户管理服务发生错误
	PermissionError   CustomError // 权限错误
	IllegalRequest    CustomError // 非法请求
	FrequentRequest   CustomError // 频繁访问
	UnknownError      CustomError // 未知错误
	JWTToknError      CustomError // 鉴权错误
}

var Errors = CustomErrors{
	LocalServiceError: CustomError{50000, "内部服务错误"},
	UserServiceError:  CustomError{51000, "远程服务错误"},
	PermissionError:   CustomError{52000, "没有权限访问"},
	IllegalRequest:    CustomError{53000, "非法请求"},
	FrequentRequest:   CustomError{54000, "访问过于频繁，请稍后再试"},
	UnknownError:      CustomError{55000, "未知错误"},
	JWTToknError:      CustomError{56000, "Token 验证失败"},
}
