package common

import "omiai-server/pkg/storage"

type Controller struct {
	storage storage.Driver
}

func NewController(storage storage.Driver) *Controller {
	return &Controller{storage: storage}
}
