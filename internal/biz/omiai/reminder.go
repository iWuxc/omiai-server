package biz_omiai

import (
	"context"
	"time"
)

// 提醒类型常量
const (
	ReminderTypeFollowUp    = 1 // 回访提醒
	ReminderTypeBirthday    = 2 // 生日提醒
	ReminderTypeAnniversary = 3 // 纪念日提醒
	ReminderTypeChurnRisk   = 4 // 流失预警
)

// 优先级常量
const (
	ReminderPriorityLow    = 1 // 低
	ReminderPriorityMedium = 2 // 中
	ReminderPriorityHigh   = 3 // 高
)

// Reminder 提醒记录模型
type Reminder struct {
	ID            uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        uint64    `json:"user_id" gorm:"column:user_id;index;comment:用户ID（红娘ID）"`
	Type          int8      `json:"type" gorm:"column:type;comment:提醒类型：1=回访，2=生日，3=纪念日，4=流失预警"`
	ClientID      *uint64   `json:"client_id" gorm:"column:client_id;index;comment:关联客户ID"`
	MatchRecordID *uint64   `json:"match_record_id" gorm:"column:match_record_id;comment:关联匹配记录ID"`
	Title         string    `json:"title" gorm:"column:title;size:255;comment:提醒标题"`
	Content       string    `json:"content" gorm:"column:content;type:text;comment:提醒内容"`
	RemindAt      time.Time `json:"remind_at" gorm:"column:remind_at;index;comment:提醒时间"`
	IsRead        int8      `json:"is_read" gorm:"column:is_read;default:0;comment:是否已读：0=未读，1=已读"`
	IsDone        int8      `json:"is_done" gorm:"column:is_done;default:0;comment:是否已完成：0=未完成，1=已完成"`
	Priority      int8      `json:"priority" gorm:"column:priority;default:2;comment:优先级：1=低，2=中，3=高"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
	// 关联对象
	Client      *Client      `json:"client" gorm:"foreignKey:ClientID"`
	MatchRecord *MatchRecord `json:"match_record" gorm:"foreignKey:MatchRecordID"`
}

func (t *Reminder) TableName() string {
	return "reminder"
}

// ReminderInterface 提醒数据层接口
type ReminderInterface interface {
	// 基础CRUD
	Create(ctx context.Context, reminder *Reminder) error
	Update(ctx context.Context, reminder *Reminder) error
	Get(ctx context.Context, id uint64) (*Reminder, error)
	Delete(ctx context.Context, id uint64) error

	// 业务查询
	SelectByUser(ctx context.Context, userID uint64, isDone int8, offset, limit int) ([]*Reminder, error)
	CountByUser(ctx context.Context, userID uint64, isDone int8) (int64, error)
	GetTodayReminders(ctx context.Context, userID uint64) ([]*Reminder, error)
	GetPendingReminders(ctx context.Context, userID uint64) ([]*Reminder, error)

	// 标记操作
	MarkAsRead(ctx context.Context, id uint64) error
	MarkAsDone(ctx context.Context, id uint64) error

	// 批量检查是否存在
	ExistsByClientAndType(ctx context.Context, clientID uint64, reminderType int8, startTime, endTime time.Time) (bool, error)
}

// ReminderStats 提醒统计
type ReminderStats struct {
	Total        int64 `json:"total"`         // 总提醒数
	Pending      int64 `json:"pending"`       // 待处理
	Today        int64 `json:"today"`         // 今日提醒
	HighPriority int64 `json:"high_priority"` // 高优先级
}
