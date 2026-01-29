package queues

import (
	"omiai-server/internal/conf"
	"omiai-server/pkg/trace"
	"time"

	"github.com/google/wire"
	logger "github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/queue"
	"github.com/iWuxc/go-wit/queue/server"
	"github.com/sirupsen/logrus"
)

var (
	log         logger.Logger
	ProviderSet = wire.NewSet(
		wire.Struct(new(InitQueue), "*"),
		NewQueue,
		NewOutfitRatingQueue,
	)
)

const (
	QueueOutfitRating = "omiai-server:outfit_rating"
)

type InitQueue struct {
	OutfitRatingQueue *OutfitRatingQueue
}

func queueHandle(q *InitQueue) *queue.ServeMux {
	mux := queue.NewServeMux()
	mux.Handle(OutfitRatingTask, q.OutfitRatingQueue)

	return mux
}

func NewQueue(q *InitQueue) *server.Server {
	log = logger.NewLogger("queues",
		func(log *logrus.Logger) {
			log.AddHook(trace.NewLogCtx(trace.GroupQueue))
		},
		logger.SetOutPath(conf.GetConfig().Log.Path),
		logger.SetOutputWithRotationTime("omiai-server", conf.GetConfig().Log.MaxAge, time.Duration(conf.GetConfig().Log.RotationTimeHour)*time.Hour),
		logger.SetOutPutLevel(conf.GetConfig().Log.Level),
		logger.SetOutFormat(logger.JsonFormat()),
	)

	mux := queueHandle(q)
	var mainSrv *server.Server

	// 指定每个队列并发数量
	concurrencyQueues := map[string]int{
		QueueOutfitRating: 3, //搭配评分队列
	}

	// 为每个指定并发的队列，启动一个独立实例
	if len(concurrencyQueues) > 0 {
		for qname, conc := range concurrencyQueues {
			s := server.NewServer(server.Config{
				Concurrency: conc,                     //启动队列数量
				Queues:      map[string]int{qname: 1}, //队列名称
			})
			go func(srv *server.Server, name string) {
				if err := srv.Run(mux); err != nil {
					log.Errorf("queues(%s) error: %s", name, err.Error())
				}
			}(s, qname)

			// 若没有默认实例，则返回第一个专用实例
			if mainSrv == nil {
				mainSrv = s
			}
		}
	}

	return mainSrv
}
