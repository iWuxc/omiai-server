package event

import (
	"context"
	"time"
)

// Event 线下活动表 (如：8分钟交友、剧本杀相亲)
type Event struct {
	ID          uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID    uint64    `json:"tenant_id" gorm:"column:tenant_id;index;default:1;comment:所属租户(机构)ID"`
	Title       string    `json:"title" gorm:"column:title;size:128;not null;comment:活动标题"`
	Cover       string    `json:"cover" gorm:"column:cover;size:255;comment:活动封面图"`
	Description string    `json:"description" gorm:"column:description;type:text;comment:活动详情(富文本)"`
	Address     string    `json:"address" gorm:"column:address;size:255;comment:活动地址"`
	StartTime   time.Time `json:"start_time" gorm:"column:start_time;comment:活动开始时间"`
	EndTime     time.Time `json:"end_time" gorm:"column:end_time;comment:活动结束时间"`
	PriceCoins  int       `json:"price_coins" gorm:"column:price_coins;comment:报名费用(红豆)"`
	MaxQuota    int       `json:"max_quota" gorm:"column:max_quota;comment:最大报名人数"`
	JoinedCount int       `json:"joined_count" gorm:"column:joined_count;default:0;comment:已报名人数"`
	Status      int8      `json:"status" gorm:"column:status;default:1;comment:状态 1报名中 2已满员 3进行中 4已结束 5已取消"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (t *Event) TableName() string {
	return "event"
}

// EventRegistration 活动报名记录表
type EventRegistration struct {
	ID        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID  uint64    `json:"tenant_id" gorm:"column:tenant_id;index;default:1;comment:所属租户(机构)ID"`
	EventID   uint64    `json:"event_id" gorm:"column:event_id;index;comment:活动ID"`
	ClientID  uint64    `json:"client_id" gorm:"column:client_id;index;comment:报名用户ID"`
	CostCoins int       `json:"cost_coins" gorm:"column:cost_coins;comment:实际支付红豆"`
	Status    int8      `json:"status" gorm:"column:status;default:1;comment:状态 1已报名 2已签到 3已取消"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (t *EventRegistration) TableName() string {
	return "event_registration"
}

// EventInterface 定义数据层接口
type EventInterface interface {
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Get(ctx context.Context, id uint64) (*Event, error)
	List(ctx context.Context, status int8, offset, limit int) ([]*Event, int64, error)
	
	// 报名相关
	Register(ctx context.Context, eventID, clientID uint64, costCoins int) error
	HasRegistered(ctx context.Context, eventID, clientID uint64) (bool, error)
}
