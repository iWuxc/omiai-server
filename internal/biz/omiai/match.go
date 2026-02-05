package biz_omiai

import (
	"context"
	"omiai-server/internal/biz"
	"time"
)

const (
	MatchStatusAcquaintance = 1 // 相识
	MatchStatusDating       = 2 // 交往
	MatchStatusStable       = 3 // 稳定
	MatchStatusEngagement   = 4 // 订婚
	MatchStatusMarried      = 5 // 结婚
	MatchStatusBroken       = 6 // 分手
)

const (
	ClientStatusSingle   = 1 // 单身
	ClientStatusMatching = 2 // 匹配中
	ClientStatusMatched  = 3 // 已匹配
	ClientStatusStopped  = 4 // 停止服务
)

// MatchRecord 匹配成功记录 (情侣档案)
type MatchRecord struct {
	ID             uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MaleClientID   uint64    `json:"male_client_id" gorm:"column:male_client_id;index;comment:男方ID"`
	FemaleClientID uint64    `json:"female_client_id" gorm:"column:female_client_id;index;comment:女方ID"`
	MatchDate      time.Time `json:"match_date" gorm:"column:match_date;comment:匹配确认时间"`
	MatchScore     int       `json:"match_score" gorm:"column:match_score;comment:匹配得分"`
	Status         int8      `json:"status" gorm:"column:status;default:1;comment:状态 1相识 2交往 3稳定 4订婚 5结婚 6分手"`
	Remark         string    `json:"remark" gorm:"column:remark;type:text;comment:备注"`
	AdminID        string    `json:"admin_id" gorm:"column:admin_id;size:64;comment:操作管理员ID"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`

	// 关联对象 (查询时使用)
	MaleClient   *Client `json:"male_client" gorm:"foreignKey:MaleClientID"`
	FemaleClient *Client `json:"female_client" gorm:"foreignKey:FemaleClientID"`
}

func (t *MatchRecord) TableName() string {
	return "match_record"
}

// Candidate 匹配候选人
type Candidate struct {
	CandidateID uint64   `json:"candidate_id"`
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	MatchScore  int      `json:"match_score"`
	Tags        []string `json:"tags"`
	Age         int      `json:"age"`
	Height      int      `json:"height"`
	Education   int      `json:"education"`
}

// Comparison 匹配对比详情
type Comparison struct {
	BasicInfo                map[string]map[string]interface{} `json:"basic_info"`
	PersonalityRadar         map[string]map[string]int         `json:"personality_radar"`
	Interests                map[string]interface{}            `json:"interests"`
	Values                   map[string]interface{}            `json:"values"`
	RelationshipExpectations map[string]map[string]int         `json:"relationship_expectations"`
}

// MatchStatusHistory 状态变更记录
type MatchStatusHistory struct {
	ID            uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MatchRecordID uint64    `json:"match_record_id" gorm:"column:match_record_id;index;comment:匹配记录ID"`
	OldStatus     int8      `json:"old_status" gorm:"column:previous_status;comment:旧状态"`
	NewStatus     int8      `json:"new_status" gorm:"column:current_status;comment:新状态"`
	ChangeTime    time.Time `json:"change_time" gorm:"column:change_time;comment:变更时间"`
	Operator      string    `json:"operator" gorm:"column:operator;size:64;comment:操作人"`
	Reason        string    `json:"reason" gorm:"column:reason;size:255;comment:变更原因"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
}

func (t *MatchStatusHistory) TableName() string {
	return "match_status_history"
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

	// V2: 新增候选人与对比接口
	GetCandidates(ctx context.Context, clientID uint64) ([]*Candidate, error)
	Compare(ctx context.Context, clientID, candidateID uint64) (*Comparison, error)

	// V2: 直接确认匹配 (替换 ConfirmRequest)
	ConfirmMatch(ctx context.Context, clientID, candidateID uint64, adminID, remark string) (*MatchRecord, error)

	// 状态管理
	UpdateStatus(ctx context.Context, recordID uint64, oldStatus, newStatus int8, operator, reason string) error
	DissolveMatch(ctx context.Context, clientID uint64, operator, reason string) error
	GetStatusHistory(ctx context.Context, recordID uint64) ([]*MatchStatusHistory, error)

	// 回访相关
	CreateFollowUp(ctx context.Context, record *FollowUpRecord) error
	SelectFollowUps(ctx context.Context, matchRecordID uint64) ([]*FollowUpRecord, error)
	GetReminders(ctx context.Context) ([]*MatchRecord, error)

	// 统计分析
	Stats(ctx context.Context) (map[string]interface{}, error)
}
