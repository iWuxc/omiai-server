package match

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type Controller struct {
	db     *data.DB
	match  biz_omiai.MatchInterface
	client biz_omiai.ClientInterface
}

func NewController(db *data.DB, match biz_omiai.MatchInterface, client biz_omiai.ClientInterface) *Controller {
	return &Controller{db: db, match: match, client: client}
}
