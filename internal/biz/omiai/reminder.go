package biz_omiai

import (
	"time"

	"gorm.io/gorm"
)

// AutoReminderRule 自动提醒规则
type AutoReminderRule struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name             string         `json:"name" gorm:"size:64;not null;comment:规则名称"`
	TriggerType      string         `json:"trigger_type" gorm:"size:32;not null;comment:触发类型(NewClient, StatusChange, NoContact)"`
	TriggerCondition string         `json:"trigger_condition" gorm:"size:255;comment:触发条件(如:status=2)"`
	DelayDays        int            `json:"delay_days" gorm:"default:0;comment:延迟天数"`
	TemplateID       int64          `json:"template_id" gorm:"comment:关联的沟通模板ID"`
	IsEnabled        bool           `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ReminderTask 提醒任务
type ReminderTask struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientID    int64     `json:"client_id" gorm:"not null;index;comment:关联客户ID"`
	RuleID      int64     `json:"rule_id" gorm:"index;comment:关联规则ID"`
	Content     string    `json:"content" gorm:"type:text;comment:提醒内容/建议话术"`
	ScheduledAt time.Time `json:"scheduled_at" gorm:"index;comment:计划提醒时间"`
	Status      string    `json:"status" gorm:"size:20;default:'pending';comment:状态(pending, completed, cancelled)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReminderInterface interface {
	CreateRule(rule *AutoReminderRule) error
	ListRules() ([]*AutoReminderRule, error)
	GetRule(id int64) (*AutoReminderRule, error)
	UpdateRule(rule *AutoReminderRule) error

	CreateTask(task *ReminderTask) error
	ListPendingTasks() ([]*ReminderTask, error)
	CompleteTask(id int64) error
	GetTasksByClient(clientID int64) ([]*ReminderTask, error)

	GetTodayReminders(userID uint64) ([]*ReminderTask, error)
	GetPendingReminders(userID uint64) ([]*ReminderTask, error)
	MarkAsRead(id int64) error
	MarkAsDone(id int64) error
	Delete(id int64) error
	CountByUser(userID uint64, isDone int) (int64, error)
	ExistsByClientAndType(clientID uint64, triggerType string, start, end time.Time) (bool, error)
}
