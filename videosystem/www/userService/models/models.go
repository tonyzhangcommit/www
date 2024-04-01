/*
这里是数据库表设计
*/

package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Username          string       `json:"username" gorm:"index;not null;unique;column:username;comment:用户名"`
	Password          string       `json:"-" gorm:"not null;column:password;comment:密码"`
	PhoneNumber       string       `json:"phonenumber" gorm:"index;unique;not null;column:phonenumber;comment:手机号"`
	Roles             []Role       `json:"roles" gorm:"many2many:user_roles;"`
	ExtraPermissions  []Permission `json:"-" gorm:"many2many:user_extra_permissions;"`
	DeniedPermissions []Permission `json:"-" gorm:"many2many:user_denied_permissions;"`
	ParentID          *uint        `gorm:"column:parentid;"`
	Children          []User       `gorm:"foreignKey:ParentID"`
	AgentCode         string       `gorm:"uniqueIndex;column:agentcode;default:NULL"`
	ParentAgentCode   string       `json:"parentagentcode" gorm:"column:parentagentcode;comment:父级识别码;not null"`
	IsBanned          bool         `json:"isbanned" gorm:"not null;column:isbanned;default:false;comment:是否封禁"`
	Profile           Profile      `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	profile := Profile{UserID: u.ID}
	if u.Username != "desupadmin" {
		profile.TypeVip = "月会员" // 暂时测试
	}
	err = tx.Model(&Profile{}).Create(&profile).Error
	if err != nil {
		return err
	}
	return nil
}

// Profile 用户资料表
type Profile struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserID         uint
	Address        string    `json:"address" gorm:"column:address;comment:地区"`
	Sex            uint      `json:"sex" gorm:"column:sex;comment:性别"`
	Identification string    `json:"idcard" gorm:"uniqueIndex;column:idcard;comment:身份证号;default:NULL"`
	Email          string    `json:"email" gorm:"uniqueIndex;column:email;comment:邮箱;default:NULL"`
	VIP            bool      `json:"vip" gorm:"default:false;column:isvip;comment:是否会员"`
	TypeVip        string    `json:"typevip" gorm:"column:typevip;comment:会员类型"`
	ExpVipDate     time.Time `json:"expvipdate" gorm:"autoCreateTime;column:expvipdate;comment:会员过期时间"`
	Preferences    string    `json:"preferences" gorm:"column:preferences;comment:偏好"`
}

// 代理管理表
type AgentManagement struct {
	gorm.Model
	AgentUserID uint `json:"agentuserID" gorm:"column:agentuserID;comment:代理ID"`
	CanDevelop  bool `json:"candevelop" gorm:"column:candevelop;default:false;comment:能否发展自代理"`
	IsBanned    bool `json:"isbanned" gorm:"column:isbanned;default:false;comment:是否封禁"`
	ManagedBy   uint `json:"managedby" gorm:"column:managedby;comment:父级代理"`
}

type Role struct {
	gorm.Model
	RoleName    string       `json:"rolename" gorm:"unique;column:rolename;comment:角色名"`
	Description string       `json:"desc" gorm:"column:desc;comment:描述信息"`
	Users       []User       `json:"-" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"-" gorm:"many2many:role_permissions;"`
}

// Permission 权限表
type Permission struct {
	gorm.Model
	PermissionName string `json:"permissionname" gorm:"unique;column:permissionname;comment:权限名"`
	Description    string `json:"desc" gorm:"column:desc;comment:描述信息"`
}

// VerificationCodeRecord 验证码记录表
type VerificationCodeRecord struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	PhoneNumber  string `json:"phonenumber" gorm:"unique;column:phonenumber;comment:手机号"`
	BusinessType string `json:"businesstype" gorm:"column:businesstype;comment:业务类型"`
	Code         string `json:"code" gorm:"columncode:;comment:验证码"`
}

// SecurityInfo 安全信息表
type SecurityInfo struct {
	gorm.Model
	UserID      uint
	SessionID   string
	IPAddress   string
	Device      string
	LastLoginAt time.Time
	MFAEnabled  bool
}

// UserActivity 用户行为日志表
type UserActivity struct {
	gorm.Model
	UserID  uint
	Action  string
	Details string // 可以用JSON格式存储详细信息
}
