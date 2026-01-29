package api

import "github.com/libi/dcron"

// CronJobInterface .
type CronJobInterface interface {
	dcron.Job
	JobName() string
	Schedule() string
}
