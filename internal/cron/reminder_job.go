package cron

import (
	"context"
)

type ReminderCronJob struct {
	reminderService *ReminderService
}

func NewReminderCronJob(reminderService *ReminderService) *ReminderCronJob {
	return &ReminderCronJob{
		reminderService: reminderService,
	}
}

func (j *ReminderCronJob) JobName() string {
	return "GenerateDailyReminders"
}

func (j *ReminderCronJob) Schedule() string {
	// Every day at 9:00 AM
	return "0 0 9 * * *"
}

func (j *ReminderCronJob) Run() {
	ctx := context.Background()
	// Using the package level logger defined in cron.go if available, or just standard logging
	// Since we commented out logger in ReminderService to fix build, we can just run it here.
	if log != nil {
		log.Infof("Starting daily reminder generation job")
	}
	
	if err := j.reminderService.GenerateDailyReminders(ctx); err != nil {
		if log != nil {
			log.Errorf("Daily reminder generation job failed: %v", err)
		}
	} else {
		if log != nil {
			log.Infof("Daily reminder generation job completed successfully")
		}
	}
}
