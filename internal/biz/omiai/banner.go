package biz_omiai

import (
	"context"
	"omiai-server/internal/biz"
	"time"
)

const (
	BannerStatusEnable  = 1
	BannerStatusDisable = 0
)

// Banner 轮播图模型
type Banner struct {
	ID        uint64    `json:"id" gorm:"column:id"`                 // 主键ID
	Title     string    `json:"title" gorm:"column:title"`           // 轮播图标题
	ImageURL  string    `json:"image_url" gorm:"column:image_url"`   // 轮播图片URL
	SortOrder uint      `json:"sort_order" gorm:"column:sort_order"` // 排序序号，数字越小越靠前
	Status    int8      `json:"status" gorm:"column:status"`         // 状态：1-启用，2-禁用
	LinkUrl   string    `json:"link_url" gorm:"column:link_url"`     // 链接URL
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"` // 更新时间
}

// TableName 表名
func (t *Banner) TableName() string {
	return "banner"
}

type BannerInterface interface {
	Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*Banner, error)
	Create(ctx context.Context, banner *Banner) error
}
