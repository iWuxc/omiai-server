package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
}

func NewController(db *data.DB, client biz_omiai.ClientInterface) *Controller {
	return &Controller{db: db, Client: client}
}
