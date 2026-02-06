package omiai

import (
	"context"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"time"
)

var _ biz_omiai.ReminderInterface = (*ReminderRepo)(nil)

type ReminderRepo struct {
	db *data.DB
}

func NewReminderRepo(db *data.DB) biz_omiai.ReminderInterface {
	return &ReminderRepo{db: db}
}

func (r *ReminderRepo) Create(ctx context.Context, reminder *biz_omiai.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *ReminderRepo) Update(ctx context.Context, reminder *biz_omiai.Reminder) error {
	return r.db.WithContext(ctx).Model(reminder).Updates(reminder).Error
}

func (r *ReminderRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Reminder, error) {
	var reminder biz_omiai.Reminder
	err := r.db.WithContext(ctx).Preload("Client").Preload("MatchRecord").First(&reminder, id).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *ReminderRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&biz_omiai.Reminder{}, id).Error
}

func (r *ReminderRepo) SelectByUser(ctx context.Context, userID uint64, isDone int8, offset, limit int) ([]*biz_omiai.Reminder, error) {
	var list []*biz_omiai.Reminder
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if isDone >= 0 {
		db = db.Where("is_done = ?", isDone)
	}
	err := db.Order("priority desc, remind_at asc").
		Preload("Client").
		Offset(offset).
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *ReminderRepo) CountByUser(ctx context.Context, userID uint64, isDone int8) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&biz_omiai.Reminder{}).Where("user_id = ?", userID)
	if isDone >= 0 {
		db = db.Where("is_done = ?", isDone)
	}
	err := db.Count(&count).Error
	return count, err
}

func (r *ReminderRepo) GetTodayReminders(ctx context.Context, userID uint64) ([]*biz_omiai.Reminder, error) {
	var list []*biz_omiai.Reminder
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_done = ?", userID, 0).
		Where("remind_at >= ? AND remind_at < ?", todayStart, todayEnd).
		Order("priority desc, remind_at asc").
		Preload("Client").
		Find(&list).Error
	return list, err
}

func (r *ReminderRepo) GetPendingReminders(ctx context.Context, userID uint64) ([]*biz_omiai.Reminder, error) {
	var list []*biz_omiai.Reminder
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_done = ?", userID, 0).
		Where("remind_at <= ?", now).
		Order("priority desc, remind_at asc").
		Preload("Client").
		Limit(20).
		Find(&list).Error
	return list, err
}

func (r *ReminderRepo) MarkAsRead(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&biz_omiai.Reminder{}).Where("id = ?", id).Update("is_read", 1).Error
}

func (r *ReminderRepo) MarkAsDone(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&biz_omiai.Reminder{}).Where("id = ?", id).Update("is_done", 1).Error
}

func (r *ReminderRepo) ExistsByClientAndType(ctx context.Context, clientID uint64, reminderType int8, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&biz_omiai.Reminder{}).
		Where("client_id = ? AND type = ?", clientID, reminderType).
		Where("remind_at >= ? AND remind_at < ?", startTime, endTime).
		Count(&count).Error
	return count > 0, err
}
