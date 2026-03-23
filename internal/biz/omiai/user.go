package biz_omiai

import (
	"context"
	"time"
)

const (
	RoleSuperAdmin = "super_admin" // 平台超级管理员(唯一能看到所有机构的人)
	RoleAdmin      = "admin"       // 机构管理员(店长)
	RoleOperator   = "operator"    // 机构普通红娘
)

// User 系统用户模型 (B端红娘/管理员)
type User struct {
	ID        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID  uint64    `json:"tenant_id" gorm:"column:tenant_id;index;default:1;comment:所属租户(机构)ID"`
	Phone     string    `json:"phone" gorm:"column:phone;size:20;uniqueIndex;comment:手机号"`
	Password  string    `json:"-" gorm:"column:password;size:128;comment:密码"`
	Nickname  string    `json:"nickname" gorm:"column:nickname;size:64;comment:昵称"`
	Avatar    string    `json:"avatar" gorm:"column:avatar;size:255;comment:头像"`
	Role      string    `json:"role" gorm:"column:role;size:20;default:operator;comment:角色 admin/operator"`
	WxOpenID  string    `json:"wx_openid" gorm:"column:wx_openid;size:128;index;comment:微信OpenID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "user"
}

type UserInterface interface {
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetByWxOpenID(ctx context.Context, openID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
}
