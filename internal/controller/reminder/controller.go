package reminder

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type Controller struct {
	db           *data.DB
	reminderRepo biz_omiai.ReminderInterface
}

func NewController(db *data.DB, reminderRepo biz_omiai.ReminderInterface) *Controller {
	return &Controller{db: db, reminderRepo: reminderRepo}
}
