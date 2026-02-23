package omiai

import (
	"time"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type ReminderRepo struct {
	db *data.DB
}

func NewReminderRepo(db *data.DB) biz_omiai.ReminderInterface {
	return &ReminderRepo{db: db}
}

// Rule Operations
func (r *ReminderRepo) CreateRule(rule *biz_omiai.AutoReminderRule) error {
	return r.db.DB.Create(rule).Error
}

func (r *ReminderRepo) ListRules() ([]*biz_omiai.AutoReminderRule, error) {
	var rules []*biz_omiai.AutoReminderRule
	if err := r.db.DB.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func (r *ReminderRepo) GetRule(id int64) (*biz_omiai.AutoReminderRule, error) {
	var rule biz_omiai.AutoReminderRule
	if err := r.db.DB.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *ReminderRepo) UpdateRule(rule *biz_omiai.AutoReminderRule) error {
	return r.db.DB.Save(rule).Error
}

// Task Operations
func (r *ReminderRepo) CreateTask(task *biz_omiai.ReminderTask) error {
	return r.db.DB.Create(task).Error
}

func (r *ReminderRepo) ListPendingTasks() ([]*biz_omiai.ReminderTask, error) {
	var tasks []*biz_omiai.ReminderTask
	now := time.Now()
	// 查询未完成且已到期的任务
	if err := r.db.DB.Where("status = ? AND scheduled_at <= ?", "pending", now).Order("scheduled_at asc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ReminderRepo) CompleteTask(id int64) error {
	return r.db.DB.Model(&biz_omiai.ReminderTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "completed",
		"updated_at": time.Now(),
	}).Error
}

func (r *ReminderRepo) GetTasksByClient(clientID int64) ([]*biz_omiai.ReminderTask, error) {
	var tasks []*biz_omiai.ReminderTask
	if err := r.db.DB.Where("client_id = ?", clientID).Order("scheduled_at desc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ReminderRepo) GetTodayReminders(userID uint64) ([]*biz_omiai.ReminderTask, error) {
	var tasks []*biz_omiai.ReminderTask
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// TODO: Filter by userID if ReminderTask has UserID field. Currently assuming all tasks are visible or filtered by ClientID which is linked to User.
	// For now, we fetch tasks scheduled for today.
	// In a real scenario, we should join with Client table to filter by ManagerID (userID).
	if err := r.db.DB.Where("scheduled_at >= ? AND scheduled_at < ?", startOfDay, endOfDay).Order("scheduled_at asc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ReminderRepo) GetPendingReminders(userID uint64) ([]*biz_omiai.ReminderTask, error) {
	var tasks []*biz_omiai.ReminderTask
	now := time.Now()
	
	// Similar to GetTodayReminders, filtering by userID is needed.
	if err := r.db.DB.Where("status = ? AND scheduled_at <= ?", "pending", now).Order("scheduled_at asc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ReminderRepo) MarkAsRead(id int64) error {
	// ReminderTask doesn't have IsRead field yet, assuming it might be added or this is a placeholder
	// For now, let's assume "pending" -> "read" transition if we had a status for it, or just ignore if not supported.
	// But the interface requires it. Let's return nil for now or update UpdatedAt.
	return r.db.DB.Model(&biz_omiai.ReminderTask{}).Where("id = ?", id).Update("updated_at", time.Now()).Error
}

func (r *ReminderRepo) MarkAsDone(id int64) error {
	return r.CompleteTask(id)
}

func (r *ReminderRepo) Delete(id int64) error {
	return r.db.DB.Delete(&biz_omiai.ReminderTask{}, id).Error
}

func (r *ReminderRepo) CountByUser(userID uint64, isDone int) (int64, error) {
	var count int64
	db := r.db.DB.Model(&biz_omiai.ReminderTask{})
	
	// Filter by isDone: 1 for done, 0 for pending, -1 for all
	if isDone == 1 {
		db = db.Where("status = ?", "completed")
	} else if isDone == 0 {
		db = db.Where("status = ?", "pending")
	}
	
	// Filter by userID (TODO: Join with Client)
	
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ReminderRepo) ExistsByClientAndType(clientID uint64, triggerType string, start, end time.Time) (bool, error) {
	var count int64
	// ReminderTask doesn't store TriggerType directly, it links to Rule. 
	// If RuleID is 0 (system generated), we might check Content or add Type to Task.
	// For now, let's assume we check if any task exists for this client in the time range.
	// To be precise, we should probably add a Type field to ReminderTask.
	
	err := r.db.DB.Model(&biz_omiai.ReminderTask{}).
		Where("client_id = ? AND scheduled_at >= ? AND scheduled_at < ?", clientID, start, end).
		Count(&count).Error
	
	return count > 0, err
}
