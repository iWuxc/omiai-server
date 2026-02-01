package omiai

import (
	"context"
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"time"

	"gorm.io/gorm"
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

func (r *MatchRepo) CreateFollowUp(ctx context.Context, record *biz_omiai.FollowUpRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *MatchRepo) SelectFollowUps(ctx context.Context, matchRecordID uint64) ([]*biz_omiai.FollowUpRecord, error) {
	var list []*biz_omiai.FollowUpRecord
	err := r.db.WithContext(ctx).Where("match_record_id = ?", matchRecordID).Order("follow_up_date desc").Find(&list).Error
	return list, err
}

func (r *MatchRepo) GetReminders(ctx context.Context) ([]*biz_omiai.MatchRecord, error) {
	// Simple reminder: matched for more than 30 days but no followup in last 30 days
	// Or based on next_follow_up_at
	var list []*biz_omiai.MatchRecord
	now := time.Now()
	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN follow_up_record ON follow_up_record.match_record_id = match_record.id").
		Where("match_record.status = ?", biz_omiai.MatchStatusMatched).
		Where("follow_up_record.next_follow_up_at <= ? OR follow_up_record.id IS NULL", now).
		Preload("MaleClient").Preload("FemaleClient").
		Find(&list).Error
	return list, err
}
