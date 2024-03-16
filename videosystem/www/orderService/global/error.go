package global

/*
	封装错误码
*/

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError CustomError // 业务逻辑错误
	ValidateError CustomError // 参数错误
}

var Errors = CustomErrors{
	BusinessError: CustomError{40040, "内部错误"},
	ValidateError: CustomError{42040, "参数错误"},
}
