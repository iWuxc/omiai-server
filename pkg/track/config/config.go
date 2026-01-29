package config

import (
	"errors"
	"github.com/iWuxc/go-wit/utils"
	"omiai-server/pkg/track"
)

type TrackConf struct {
	Path   string `json:"path"`
	MaxAge uint32 `json:"max_age"`
}

func TrackInit(conf *TrackConf) error {

	if conf == nil || conf.Path == "" || conf.MaxAge == 0 {
		panic(errors.New("track conf error"))
	}
	managerConf := track.ManagerConf{
		Path:   conf.Path,
		MaxAge: conf.MaxAge,
	}
	utils.DebugInfo("track init", managerConf)
	track.NewManager(managerConf)
	return nil
}
