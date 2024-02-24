package response

/*
	保存非常规格式返回值格式
*/

type LoginSuccess struct {
	Jwt string `json:"jwt"`
}
