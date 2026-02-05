package omiai

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"sort"
	"time"

	"github.com/iWuxc/go-wit/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ biz_omiai.MatchInterface = (*MatchRepo)(nil)

type MatchRepo struct {
	db *data.DB
}

func NewMatchRepo(db *data.DB) biz_omiai.MatchInterface {
	return &MatchRepo{db: db}
}

func (r *MatchRepo) Select(ctx context.Context, clause *biz.WhereClause, offset, limit int) ([]*biz_omiai.MatchRecord, error) {
	var list []*biz_omiai.MatchRecord
	db := r.db.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).
		Preload("MaleClient").Preload("FemaleClient")

	if clause.Where != "" {
		db = db.Where(clause.Where, clause.Args...)
	}

	err := db.Order(clause.OrderBy).Offset(offset).Limit(limit).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("MatchRepo:Select err:%w", err)
	}
	return list, nil
}

func (r *MatchRepo) Create(ctx context.Context, record *biz_omiai.MatchRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Match Record
		if err := tx.WithContext(ctx).Create(record).Error; err != nil {
			return err
		}
		// 2. Update Client Statuses to "Matched"
		if err := tx.WithContext(ctx).Model(&biz_omiai.Client{}).Where("id IN ?", []uint64{record.MaleClientID, record.FemaleClientID}).
			Update("status", biz_omiai.ClientStatusMatched).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *MatchRepo) Update(ctx context.Context, record *biz_omiai.MatchRecord) error {
	return r.db.WithContext(ctx).Model(record).Where("id = ?", record.ID).Updates(record).Error
}

func (r *MatchRepo) Get(ctx context.Context, id uint64) (*biz_omiai.MatchRecord, error) {
	var record biz_omiai.MatchRecord
	err := r.db.WithContext(ctx).Preload("MaleClient").Preload("FemaleClient").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *MatchRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var record biz_omiai.MatchRecord
		if err := tx.First(&record, id).Error; err != nil {
			return err
		}
		// 1. Delete Match Record
		if err := tx.Delete(&record).Error; err != nil {
			return err
		}
		// 2. Reset Client Statuses to "Single"
		if err := tx.Model(&biz_omiai.Client{}).Where("id IN ?", []uint64{record.MaleClientID, record.FemaleClientID}).
			Update("status", biz_omiai.ClientStatusSingle).Error; err != nil {
			return err
		}
		return nil
	})
}

// V2: GetCandidates 获取候选人列表
func (r *MatchRepo) GetCandidates(ctx context.Context, clientID uint64) ([]*biz_omiai.Candidate, error) {
	var client biz_omiai.Client
	if err := r.db.WithContext(ctx).First(&client, clientID).Error; err != nil {
		return nil, err
	}

	// 1. Try Cache
	if client.CandidateCacheJSON != "" {
		var candidates []*biz_omiai.Candidate
		if err := json.Unmarshal([]byte(client.CandidateCacheJSON), &candidates); err == nil && len(candidates) > 0 {
			return candidates, nil
		}
	}

	// 2. Fallback: Calculate Real-time
	targetGender := 1
	if client.Gender == 1 {
		targetGender = 2
	}

	var potentialMatches []*biz_omiai.Client
	if err := r.db.WithContext(ctx).Where("gender = ? AND status = ?", targetGender, biz_omiai.ClientStatusSingle).
		Limit(100).Find(&potentialMatches).Error; err != nil {
		return nil, err
	}

	var candidates []*biz_omiai.Candidate
	for _, match := range potentialMatches {
		score := 60 + rand.Intn(40) // Mock Algo

		tags := []string{}
		if match.Education == client.Education {
			tags = append(tags, "学历相当")
		}
		// Simple age gap check
		ageGap := client.RealAge() - match.RealAge()
		if ageGap < 0 {
			ageGap = -ageGap
		}
		if ageGap <= 3 {
			tags = append(tags, "年龄相仿")
		}
		if len(tags) == 0 {
			tags = append(tags, "推荐")
		}

		candidates = append(candidates, &biz_omiai.Candidate{
			CandidateID: match.ID,
			Name:        match.Name,
			Avatar:      match.Avatar,
			MatchScore:  score,
			Tags:        tags,
			Age:         match.RealAge(),
			Height:      match.Height,
			Education:   int(match.Education),
		})
	}

	// Sort by score desc
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MatchScore > candidates[j].MatchScore
	})
	if len(candidates) > 20 {
		candidates = candidates[:20]
	}

	// 3. Update Cache (Lazy Load)
	if bytes, err := json.Marshal(candidates); err == nil {
		// Ignore error on update
		r.db.WithContext(ctx).Model(&client).Update("candidate_cache_json", string(bytes))
	}

	return candidates, nil
}

// V2: Compare 比较详情
func (r *MatchRepo) Compare(ctx context.Context, clientID, candidateID uint64) (*biz_omiai.Comparison, error) {
	var c1, c2 biz_omiai.Client
	if err := r.db.WithContext(ctx).First(&c1, clientID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&c2, candidateID).Error; err != nil {
		return nil, err
	}

	// Build Comparison Data
	comp := &biz_omiai.Comparison{
		BasicInfo: map[string]map[string]interface{}{
			"age": {
				"client":    c1.RealAge(),
				"candidate": c2.RealAge(),
				"diff":      fmt.Sprintf("相差%d岁", c1.RealAge()-c2.RealAge()),
			},
			"height": {
				"client":    c1.Height,
				"candidate": c2.Height,
				"diff":      fmt.Sprintf("相差%dcm", c1.Height-c2.Height),
			},
			"education": {
				"client":    c1.Education, // Value 1-5
				"candidate": c2.Education,
				"match":     c1.Education == c2.Education,
			},
			"income": {
				"client":    c1.Income,
				"candidate": c2.Income,
				"match":     true, // Simple placeholder
			},
		},
		PersonalityRadar: map[string]map[string]int{
			"openness":          {"client": 80, "candidate": 75},
			"conscientiousness": {"client": 60, "candidate": 85},
			"extraversion":      {"client": 70, "candidate": 65},
			"agreeableness":     {"client": 90, "candidate": 90},
			"neuroticism":       {"client": 40, "candidate": 30},
		},
		Interests: map[string]interface{}{
			"overlap_percentage": 0.8,
			"common_list":        []string{"旅行", "摄影"},
		},
		Values: map[string]interface{}{
			"match_percentage": 0.9,
			"details":          []string{"家庭观念一致"},
		},
		RelationshipExpectations: map[string]map[string]int{
			"short_term": {"client": 2, "candidate": 1},
			"long_term":  {"client": 5, "candidate": 5},
		},
	}
	return comp, nil
}

// V2: ConfirmMatch 直接确认匹配
func (r *MatchRepo) ConfirmMatch(ctx context.Context, clientID, candidateID uint64, adminID string) (*biz_omiai.MatchRecord, error) {
	// 0. Distributed Lock using Redis
	lockKey := fmt.Sprintf("lock:match:client:%d:%d", clientID, candidateID)
	// Try to acquire lock for 10 seconds
	lock := redis.GetRedis().GetClient().SetNX(ctx, lockKey, 1, 10*time.Second)
	if err := lock.Err(); err != nil {
		return nil, fmt.Errorf("system busy, please try again")
	}
	if !lock.Val() {
		return nil, fmt.Errorf("matching in progress")
	}
	defer redis.GetRedis().GetClient().Del(ctx, lockKey)

	return r.confirmMatchDB(ctx, clientID, candidateID, adminID)
}

func (r *MatchRepo) confirmMatchDB(ctx context.Context, clientID, candidateID uint64, adminID string) (*biz_omiai.MatchRecord, error) {
	var matchRecord *biz_omiai.MatchRecord
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get Clients and Verify Status (Double Check)
		var c1, c2 biz_omiai.Client
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&c1, clientID).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&c2, candidateID).Error; err != nil {
			return err
		}

		if c1.Status != biz_omiai.ClientStatusSingle || c2.Status != biz_omiai.ClientStatusSingle {
			return fmt.Errorf("one or both clients are already matched")
		}

		maleID, femaleID := clientID, candidateID
		if c1.Gender == 2 { // If c1 is female
			maleID, femaleID = candidateID, clientID
		}

		// 2. Create Match Record
		matchRecord = &biz_omiai.MatchRecord{
			MaleClientID:   maleID,
			FemaleClientID: femaleID,
			MatchDate:      time.Now(),
			Status:         biz_omiai.MatchStatusAcquaintance,
			MatchScore:     85, // Mock score
			AdminID:        adminID,
		}
		if err := tx.WithContext(ctx).Create(matchRecord).Error; err != nil {
			return err
		}

		// 3. Update Client Statuses and Partner ID
		// Update Client 1
		if err := tx.WithContext(ctx).Model(&biz_omiai.Client{}).Where("id = ?", clientID).
			Updates(map[string]interface{}{
				"status":     biz_omiai.ClientStatusMatched,
				"partner_id": candidateID,
			}).Error; err != nil {
			return err
		}
		// Update Client 2
		if err := tx.WithContext(ctx).Model(&biz_omiai.Client{}).Where("id = ?", candidateID).
			Updates(map[string]interface{}{
				"status":     biz_omiai.ClientStatusMatched,
				"partner_id": clientID,
			}).Error; err != nil {
			return err
		}
		return nil
	})
	return matchRecord, err
}

// UpdateStatus 更新匹配状态并记录历史
func (r *MatchRepo) UpdateStatus(ctx context.Context, recordID uint64, oldStatus, newStatus int8, operator, reason string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update Match Status
		if err := tx.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).Where("id = ?", recordID).Update("status", newStatus).Error; err != nil {
			return err
		}

		// 2. Create History Record
		history := &biz_omiai.MatchStatusHistory{
			MatchRecordID: recordID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			ChangeTime:    time.Now(),
			Operator:      operator,
			Reason:        reason,
		}
		if err := tx.WithContext(ctx).Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

// DissolveMatch 解除匹配关系
func (r *MatchRepo) DissolveMatch(ctx context.Context, clientID uint64, operator, reason string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get Client and verify status
		var client biz_omiai.Client
		if err := tx.First(&client, clientID).Error; err != nil {
			return err
		}
		if client.Status != biz_omiai.ClientStatusMatched || client.PartnerID == nil {
			return fmt.Errorf("client is not currently matched")
		}

		partnerID := *client.PartnerID

		// 2. Find Active Match Record
		var matchRecord biz_omiai.MatchRecord
		if err := tx.Where("((male_client_id = ? AND female_client_id = ?) OR (male_client_id = ? AND female_client_id = ?)) AND status != ?",
			clientID, partnerID, partnerID, clientID, biz_omiai.MatchStatusBroken).
			First(&matchRecord).Error; err != nil {
			return fmt.Errorf("active match record not found: %v", err)
		}

		// 3. Update Match Record Status
		oldStatus := matchRecord.Status
		newStatus := biz_omiai.MatchStatusBroken
		if err := tx.Model(&matchRecord).Updates(map[string]interface{}{
			"status": newStatus,
		}).Error; err != nil {
			return err
		}

		// 4. Update Clients Status
		if err := tx.Model(&biz_omiai.Client{}).Where("id IN ?", []uint64{clientID, partnerID}).
			Updates(map[string]interface{}{
				"status":     biz_omiai.ClientStatusSingle,
				"partner_id": nil,
			}).Error; err != nil {
			return err
		}

		// 5. Create History Record
		history := &biz_omiai.MatchStatusHistory{
			MatchRecordID: matchRecord.ID,
			OldStatus:     oldStatus,
			NewStatus:     int8(newStatus),
			ChangeTime:    time.Now(),
			Operator:      operator,
			Reason:        reason,
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetStatusHistory 获取状态历史
func (r *MatchRepo) GetStatusHistory(ctx context.Context, recordID uint64) ([]*biz_omiai.MatchStatusHistory, error) {
	var list []*biz_omiai.MatchStatusHistory
	err := r.db.WithContext(ctx).Where("match_record_id = ?", recordID).Order("change_time desc").Find(&list).Error
	return list, err
}

func (r *MatchRepo) CreateFollowUp(ctx context.Context, record *biz_omiai.FollowUpRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *MatchRepo) SelectFollowUps(ctx context.Context, matchRecordID uint64) ([]*biz_omiai.FollowUpRecord, error) {
	var list []*biz_omiai.FollowUpRecord
	err := r.db.WithContext(ctx).Where("match_record_id = ?", matchRecordID).Order("follow_up_date desc").Find(&list).Error
	return list, err
}

func (r *MatchRepo) GetReminders(ctx context.Context) ([]*biz_omiai.MatchRecord, error) {
	var list []*biz_omiai.MatchRecord
	now := time.Now()
	err := r.db.WithContext(ctx).
		Joins("JOIN follow_up_record ON follow_up_record.match_record_id = match_record.id").
		Where("match_record.status NOT IN (?)", []int{biz_omiai.MatchStatusBroken}). // Not broken
		Where("follow_up_record.next_follow_up_at <= ?", now).
		Group("match_record.id").
		Preload("MaleClient").Preload("FemaleClient").
		Find(&list).Error

	return list, err
}

// Stats 统计分析
func (r *MatchRepo) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 1. Total Matches
	var total int64
	r.db.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).Count(&total)
	stats["total_matches"] = total

	// 2. Status Distribution
	type StatusCount struct {
		Status int8
		Count  int64
	}
	var statusCounts []StatusCount
	r.db.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).Select("status, count(*) as count").Group("status").Scan(&statusCounts)

	statusMap := make(map[int8]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}
	stats["status_distribution"] = statusMap

	// 3. Married Count (Success Rate)
	marriedCount := statusMap[biz_omiai.MatchStatusMarried]
	stats["married_count"] = marriedCount
	if total > 0 {
		stats["success_rate"] = float64(marriedCount) / float64(total)
	} else {
		stats["success_rate"] = 0
	}

	return stats, nil
}
