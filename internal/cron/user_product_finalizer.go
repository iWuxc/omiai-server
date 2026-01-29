package cron

import (
	"context"
	"omiai-server/internal/data"
	"time"

	"github.com/google/uuid"
	"github.com/iWuxc/go-wit/redis"
)

type UserProductFinalizer struct {
	db *data.DB
}

func NewUserProductFinalizer(db *data.DB) *UserProductFinalizer {
	return &UserProductFinalizer{db: db}
}

func (s *UserProductFinalizer) JobName() string {
	return "UserProductFinalizer"
}

func (s *UserProductFinalizer) Schedule() string {
	// 秒 分 时 日 月 周
	return "*/2 * * * * *"
}

func (s *UserProductFinalizer) Run() {
	ctx := context.WithValue(context.Background(), "request_id", uuid.NewString())

	lockKey := "lock:UserProductFinalizer"
	lockRet := redis.GetRedis().GetClient().SetNX(ctx, lockKey, 1, time.Minute*5)
	if lockRet.Err() != nil {
		log.WithContext(ctx).Errorf("【定时任务-%s】 lockRet err:%s", s.JobName(), lockRet.Err().Error())
		return
	}
	if !lockRet.Val() {
		return
	}
	defer func() {
		_ = redis.GetRedis().Delete(ctx, lockKey)
		log.WithContext(ctx).Infof("%s end", s.JobName())
	}()
	log.WithContext(ctx).Infof("%s start", s.JobName())

	//TODO:: 这里写具体的逻辑
}
