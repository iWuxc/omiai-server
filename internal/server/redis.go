package server

import (
	"fmt"
	"omiai-server/internal/conf"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/iWuxc/go-wit/redis"
)

func NewRedisX() (*redis.Redis, error) {
	return redis.NewRedis(conf.GetConfig().Cache.URL, "aicloset")
}

// NewRedisSync redis分布式锁
func NewRedisSync() *redsync.Redsync {
	c := conf.GetConfig().Redis
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Default.Host, c.Default.Port),
		Password: c.Default.Password,
	})
	pool := goredis.NewPool(client)
	return redsync.New(pool)
}
