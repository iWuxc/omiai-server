package track

import (
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/utils"
)

var Mg Manager

type Manager struct {
	Log log.Logger
}

type ManagerConf struct {
	Path   string `json:"path"`
	MaxAge uint32 `json:"max_age"`
}

func NewManager(conf ManagerConf) {
	utils.DebugInfo("track New", conf)
	var opts []log.Option
	opts = append(opts, log.SetOutPutLevel("debug"))
	if conf.Path != "" {
		opts = append(opts, log.SetOutPath(conf.Path))
	}
	opts = append(opts, log.SetOutput("track-log.log", conf.MaxAge))
	logger := log.NewLogger("track-log", opts...)
	Mg = Manager{
		Log: logger,
	}
}
