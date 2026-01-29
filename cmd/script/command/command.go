package command

import (
	"omiai-server/internal/data"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewScript,
)

type Script struct {
	db *data.DB
}

func NewScript(db *data.DB) *Script {
	return &Script{
		db: db,
	}
}
