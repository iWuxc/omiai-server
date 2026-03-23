package tenant

import (
	"context"
	"time"
)

// Tenant 租户(机构)模型
type Tenant struct {
	ID           uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"column:name;size:128;uniqueIndex;comment:机构名称"`
	ContactName  string    `json:"contact_name" gorm:"column:contact_name;size:64;comment:联系人姓名"`
	ContactPhone string    `json:"contact_phone" gorm:"column:contact_phone;size:20;comment:联系电话"`
	Status       int8      `json:"status" gorm:"column:status;default:1;comment:状态 1正常 2停用"`
	ExpireAt     time.Time `json:"expire_at" gorm:"column:expire_at;comment:SaaS订阅到期时间"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (t *Tenant) TableName() string {
	return "tenant"
}

// TenantInterface 租户数据层接口
type TenantInterface interface {
	Create(ctx context.Context, tenant *Tenant) error
	Get(ctx context.Context, id uint64) (*Tenant, error)
	List(ctx context.Context, offset, limit int) ([]*Tenant, int64, error)
	Update(ctx context.Context, tenant *Tenant) error
}
