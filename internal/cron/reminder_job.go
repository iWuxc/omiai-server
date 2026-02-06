package cron

import (
	"context"
	"time"
)

// ReminderCronJob 提醒定时任务
type ReminderCronJob struct {
	reminderService *ReminderService
}

func NewReminderCronJob(reminderService *ReminderService) *ReminderCronJob {
	return &ReminderCronJob{reminderService: reminderService}
}

// JobName 任务名称
func (j *ReminderCronJob) JobName() string {
	return "reminder_generator"
}

// Schedule 执行计划 (每天早上8点执行)
func (j *ReminderCronJob) Schedule() string {
	return "0 0 8 * * *" // 秒 分 时 日 月 周
}

// Run 执行任务
func (j *ReminderCronJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Info("【提醒生成任务】开始执行...")
	if err := j.reminderService.GenerateDailyReminders(ctx); err != nil {
		log.Errorf("【提醒生成任务】执行失败: %v", err)
	} else {
		log.Info("【提醒生成任务】执行完成")
	}
}
