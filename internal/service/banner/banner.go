package banner

import (
	"github.com/iWuxc/go-wit/redis"
)

type Service struct {
	redis *redis.Redis
}

func NewService(redis *redis.Redis) *Service {
	return &Service{redis: redis}
}
