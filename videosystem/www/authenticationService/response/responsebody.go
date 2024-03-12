package response

import "time"

/*
	特定接口返回结构体
*/

// 用户登录信息
type LoginResInfo struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ParentID    uint      `json:"parentID"`
	Isbanned    bool      `json:"isbanned"`
	Phonenumber string    `json:"phonenumber"`
	AgentCode   string    `json:"agentcode"`
	Roles       []string  `json:"roles"`
	Username    string    `json:"username"`
}

// 登录接口返回结构体
type LoginRes struct {
	ErrorCode int          `json:"errorCode"`
	Data      LoginResInfo `json:"data"`
	Message   string       `json:"msg"`
}

// 获取用户VIP级别请求返回结构体
type UserVipType struct {
	ErrorCode int    `json:"errorCode"`
	Vtype     string `json:"data"`
	Message   string `json:"msg"`
}

// 获取活动信息（秒杀活动请求过滤使用）
type FlashEvent struct {
	Name      string    `json:"name"`
	Condition string    `json:"condition"`
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}

// 活动信息接口返回结构体
type FlashEventRes struct {
	ErrorCode int        `json:"errorCode"`
	Data      FlashEvent `json:"data"`
	Message   string     `json:"msg"`
}
