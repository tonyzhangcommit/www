package test



type Response struct {
	ErrorCode int         `json:"errorCode"`
	Data      interface{} `json:"data"`
	Message   string      `json:"msg"`
}

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type LoginResponse struct {
	ErrorCode int         `json:"errorCode"`
	Data      TokenOutPut `json:"data"`
	Message   string      `json:"msg"`
}

type TaskFlashOrderOrderServiceRes struct {
	ErrorCode int    `json:"errorCode"`
	Data      string `json:"data"`
	Message   string `json:"msg"`
}

type TaskFlashOrder struct {
	ErrorCode int                           `json:"errorCode"`
	Data      TaskFlashOrderOrderServiceRes `json:"data"`
	Message   string                        `json:"msg"`
}
