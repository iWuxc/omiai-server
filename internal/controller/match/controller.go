package match

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type Controller struct {
	db     *data.DB
	match  biz_omiai.MatchInterface
	client biz_omiai.ClientInterface
	user   biz_omiai.UserInterface
}

func NewController(db *data.DB, match biz_omiai.MatchInterface, client biz_omiai.ClientInterface, user biz_omiai.UserInterface) *Controller {
	return &Controller{db: db, match: match, client: client, user: user}
}
