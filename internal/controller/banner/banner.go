package banner

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/service/banner"
)

type Controller struct {
	db            *data.DB
	Banner        biz_omiai.BannerInterface
	bannerService *banner.Service
}

func NewController(db *data.DB, banner biz_omiai.BannerInterface, bannerService *banner.Service) *Controller {
	return &Controller{db: db, Banner: banner, bannerService: bannerService}
}
