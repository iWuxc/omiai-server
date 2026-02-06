package cron

import (
	"fmt"
	"omiai-server/internal/api"
	"omiai-server/internal/conf"
	"omiai-server/pkg/trace"
	"time"

	driver "github.com/dcron-contrib/redisdriver"
	"github.com/google/wire"
	logger "github.com/iWuxc/go-wit/log"
	"github.com/libi/dcron"
	"github.com/libi/dcron/cron"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	log             logger.Logger
	ProviderCronSet = wire.NewSet(
		wire.Struct(new(InitCron), "*"),
		NewCron,
		NewUserProductFinalizer,
		NewCandidatePreFilterService,
		NewReminderService,
		NewReminderCronJob,
	)
)

type InitCron struct {
	*UserProductFinalizer
	*CandidatePreFilterService
	*ReminderCronJob
}

func jobs(cron *InitCron) []api.CronJobInterface {
	return []api.CronJobInterface{
		cron.UserProductFinalizer,
		cron.CandidatePreFilterService,
		cron.ReminderCronJob,
	}
}
func NewCron(initCron *InitCron) (*dcron.Dcron, error) {
	log = logger.NewLogger("cron",
		func(log *logrus.Logger) {
			log.AddHook(trace.NewLogCtx(trace.GroupCron))
		},
		logger.SetOutPath(conf.GetConfig().Log.Path),
		logger.SetOutPutLevel("debug"),
		logger.SetOutputWithRotationTime("omiai-server", 24, time.Duration(1)*time.Hour),
		logger.SetOutFormat(logger.JsonFormat()),
	)
	drv := driver.NewDriver(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.GetConfig().CronConf.Host, conf.GetConfig().CronConf.Port),
		DB:       conf.GetConfig().CronConf.DB,
		Password: conf.GetConfig().CronConf.Password,
	}))

	//newDcron := dcron.NewDcron("aicloset", drv, cron.WithSeconds(), cron.WithLogger(log))
	newDcron := dcron.NewDcron(conf.GetConfig().CronConf.CronName, drv, cron.WithSeconds())
	newDcron.SetLogger(log)

	for _, job := range jobs(initCron) {
		if err := newDcron.AddJob(job.JobName(), job.Schedule(), job); err != nil {
			log.Errorf("【定时任务-%s】错误: %s", job.JobName(), err)
			return nil, err
		}
	}
	return newDcron, nil
}
