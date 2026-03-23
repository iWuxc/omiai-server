package event

import (
	"context"
	biz_event "omiai-server/internal/biz/event"
	"omiai-server/internal/data"
	"omiai-server/internal/data/scopes"

	"gorm.io/gorm"
)

type EventRepo struct {
	db *gorm.DB
	m  *biz_event.Event
}

func NewEventRepo(db *data.DB) biz_event.EventInterface {
	return &EventRepo{db: db.DB, m: new(biz_event.Event)}
}

func (r *EventRepo) Create(ctx context.Context, event *biz_event.Event) error {
	if tenantID, ok := ctx.Value("tenant_id").(uint64); ok && tenantID > 0 {
		event.TenantID = tenantID
	}
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *EventRepo) Update(ctx context.Context, event *biz_event.Event) error {
	return r.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Save(event).Error
}

func (r *EventRepo) Get(ctx context.Context, id uint64) (*biz_event.Event, error) {
	var event biz_event.Event
	err := r.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).First(&event, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *EventRepo) List(ctx context.Context, status int8, offset, limit int) ([]*biz_event.Event, int64, error) {
	var list []*biz_event.Event
	var total int64
	
	query := r.db.WithContext(ctx).Model(r.m).Scopes(scopes.TenantScope(ctx))
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	query.Count(&total)
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *EventRepo) Register(ctx context.Context, eventID, clientID uint64, costCoins int) error {
	tx := r.db.WithContext(ctx).Begin()
	
	var event biz_event.Event
	if err := tx.Scopes(scopes.TenantScope(ctx)).Set("gorm:query_option", "FOR UPDATE").First(&event, eventID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if event.JoinedCount >= event.MaxQuota {
		tx.Rollback()
		return gorm.ErrInvalidData // 可以自定义满员错误
	}
	
	reg := &biz_event.EventRegistration{
		EventID:   eventID,
		ClientID:  clientID,
		CostCoins: costCoins,
		Status:    1,
	}
	if tenantID, ok := ctx.Value("tenant_id").(uint64); ok && tenantID > 0 {
		reg.TenantID = tenantID
	}
	
	if err := tx.Create(reg).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if err := tx.Model(&event).UpdateColumn("joined_count", gorm.Expr("joined_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

func (r *EventRepo) HasRegistered(ctx context.Context, eventID, clientID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&biz_event.EventRegistration{}).
		Where("event_id = ? AND client_id = ?", eventID, clientID).
		Count(&count).Error
	return count > 0, err
}
