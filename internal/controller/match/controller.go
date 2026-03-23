package match

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/service/notification"
)

type Controller struct {
	db       *data.DB
	match    biz_omiai.MatchInterface
	client   biz_omiai.ClientInterface
	user     biz_omiai.UserInterface
	Notifier notification.Service
}

func NewController(db *data.DB, match biz_omiai.MatchInterface, client biz_omiai.ClientInterface, user biz_omiai.UserInterface, notifier notification.Service) *Controller {
	return &Controller{db: db, match: match, client: client, user: user, Notifier: notifier}
}
