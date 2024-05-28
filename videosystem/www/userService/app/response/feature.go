package response

import (
	"time"
	"userservice/models"
)

/*
	保存非常规返回值格式
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

type UserInfo struct {
	Username        string    `json:"username"`
	PhoneNumber     string    `json:"phonenumber"`
	Address         string    `json:"address"`
	Identification  string    `json:"idcard"`
	Email           string    `json:"email"`
	VIP             bool      `json:"vip"`
	TypeVip         string    `json:"typevip"`
	ExpVipDate      time.Time `json:"expvipdate"`
	Preferences     string    `json:"preferences"`
	Sex             uint      `json:"sex"`
	Roles           []string  `json:"roles"`
	AgentCode       string    `json:"agentcode"`
	ParentAgentCode string    `json:"parentagentcode"`
	IsBanned        bool      `json:"isbanned"`
}

type UserRoles struct {
	Roles []string `json:"roles"`
}

// 分页返回结构体
type UsersPagesInfo struct {
	Users       []models.User
	Total       int64
	Pagesize    int64
	CurrentPage int64
}
