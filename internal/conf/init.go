package conf

import (
	"omiai-server/pkg/trace"
	"time"

	"github.com/iWuxc/go-wit/cache"
	"github.com/iWuxc/go-wit/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// cacheInit 缓存初始化 .
func cacheInit(c *Cache) error {
	if c == nil {
		return nil
	}

	if _, err := cache.NewCache(c.Driver, c.URL); err != nil {
		return errors.Wrap(err, "cache init")
	}

	return nil
}

// logInit 日志初始化 .
func logInit(logger *Logger) error {
	if logger != nil {
		log.ReplaceLogger(log.NewLogger("omiai-server",
			func(log *logrus.Logger) {
				log.AddHook(trace.NewLogCtx(trace.GroupHttp))
			},
			log.SetOutPath(logger.Path),
			// log.SetOutput(logger.Path+"omiai-server", logger.MaxAge),
			log.SetOutputWithRotationTime("omiai-server.log", logger.MaxAge, time.Duration(logger.RotationTimeHour)*time.Hour),
			log.SetOutPutLevel(logger.Level),
			// log.SetOutFormat(log.JsonFormat()),
		))
	}

	return nil
}
