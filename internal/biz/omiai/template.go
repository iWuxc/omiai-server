package biz_omiai

import (
	"time"

	"gorm.io/gorm"
)

// CommunicationTemplate 沟通话术模板
type CommunicationTemplate struct {
	ID         int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string         `json:"title" gorm:"size:64;not null;comment:模板标题"`
	Content    string         `json:"content" gorm:"type:text;not null;comment:模板内容"`
	Category   string         `json:"category" gorm:"size:32;not null;comment:分类(如:打招呼,邀约,回访)"`
	UsageCount int            `json:"usage_count" gorm:"default:0;comment:使用次数"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type TemplateRepo interface {
	Create(template *CommunicationTemplate) error
	Update(template *CommunicationTemplate) error
	Delete(id int64) error
	Get(id int64) (*CommunicationTemplate, error)
	List(category string, page, pageSize int) ([]*CommunicationTemplate, int64, error)
	IncrementUsage(id int64) error
}
