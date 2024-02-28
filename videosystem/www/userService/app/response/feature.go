package response

import "time"

/*
	保存非常规格式返回值格式
*/

type LoginRes struct {
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
