package biz_omiai

import (
	"context"
	"omiai-server/internal/biz"
	"time"
)

const (
	MatchStatusMatched = 1 // 已匹配
	MatchStatusBroken  = 2 // 已分手
	MatchStatusMarried = 3 // 已结婚
)

const (
	ClientStatusSingle   = 1 // 单身
	ClientStatusMatching = 2 // 匹配中
	ClientStatusMatched  = 3 // 已匹配
	ClientStatusStopped  = 4 // 停止服务
)

// MatchRecord 匹配成功记录
type MatchRecord struct {
	ID             uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MaleClientID   uint64    `json:"male_client_id" gorm:"column:male_client_id;index;comment:男方ID"`
	FemaleClientID uint64    `json:"female_client_id" gorm:"column:female_client_id;index;comment:女方ID"`
	MatchDate      time.Time `json:"match_date" gorm:"column:match_date;comment:匹配时间"`
	MatchScore     int       `json:"match_score" gorm:"column:match_score;comment:匹配得分"`
	Status         int8      `json:"status" gorm:"column:status;default:1;comment:状态 1已匹配 2已分手 3已结婚"`
	Remark         string    `json:"remark" gorm:"column:remark;type:text;comment:备注"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`

	// 关联对象 (查询时使用)
	MaleClient   *Client `json:"male_client" gorm:"foreignKey:MaleClientID"`
	FemaleClient *Client `json:"female_client" gorm:"foreignKey:FemaleClientID"`
}

func (t *MatchRecord) TableName() string {
	return "match_record"
}

// FollowUpRecord 情侣回访记录
type FollowUpRecord struct {
	ID             uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MatchRecordID  uint64    `json:"match_record_id" gorm:"column:match_record_id;index;comment:匹配记录ID"`
	FollowUpDate   time.Time `json:"follow_up_date" gorm:"column:follow_up_date;comment:回访时间"`
	Method         string    `json:"method" gorm:"column:method;size:32;comment:回访方式(电话/面谈/线上)"`
	Content        string    `json:"content" gorm:"column:content;type:text;comment:回访内容"`
	Feedback       string    `json:"feedback" gorm:"column:feedback;type:text;comment:客户反馈"`
	Satisfaction   int8      `json:"satisfaction" gorm:"column:satisfaction;comment:满意度 1-5"`
	Attachments    string    `json:"attachments" gorm:"column:attachments;type:text;comment:附件列表(JSON)"`
	NextFollowUpAt time.Time `json:"next_follow_up_at" gorm:"column:next_follow_up_at;comment:下次回访提醒时间"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (t *FollowUpRecord) TableName() string {
	return "follow_up_record"
}

type MatchInterface interface {
	Select(ctx context.Context, clause *biz.WhereClause, offset, limit int) ([]*MatchRecord, error)
	Create(ctx context.Context, record *MatchRecord) error
	Update(ctx context.Context, record *MatchRecord) error
	Get(ctx context.Context, id uint64) (*MatchRecord, error)
	Delete(ctx context.Context, id uint64) error

	// 回访相关
	CreateFollowUp(ctx context.Context, record *FollowUpRecord) error
	SelectFollowUps(ctx context.Context, matchRecordID uint64) ([]*FollowUpRecord, error)
	GetReminders(ctx context.Context) ([]*MatchRecord, error)
}
