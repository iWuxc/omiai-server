package queues

import (
	"context"
	"encoding/json"
	"omiai-server/internal/data"

	logger "github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/queue"
	kitContext "github.com/iWuxc/go-wit/queue/context"
	"github.com/iWuxc/go-wit/redis"
)

const (
	OutfitRatingQueueName = "omiai-server:outfit_rating"
	OutfitRatingTask      = "outfit_rating"
)

type OutfitRatingQueueParams struct {
	RatingID      uint   `json:"rating_id"`      // 评分记录ID
	UserID        uint   `json:"user_id"`        // 用户ID
	OriginalImage string `json:"original_image"` // 原图URL（客户端上传的OSS地址）
}

type OutfitRatingQueue struct {
	db    *data.DB
	redis *redis.Redis
}

func NewOutfitRatingQueue(
	db *data.DB,
	redis *redis.Redis,
) *OutfitRatingQueue {
	return &OutfitRatingQueue{
		db:    db,
		redis: redis,
	}
}

// PushOutfitRatingQueue 评分任务到队列
func PushOutfitRatingQueue(params *OutfitRatingQueueParams) *queue.Task {
	jsonData, err := json.Marshal(params)
	if err != nil {
		logger.Errorf("评分队列任务参数序列化失败: %v", err)
		return nil
	}
	task := queue.NewTask(OutfitRatingTask, jsonData)
	return task
}

func (q *OutfitRatingQueue) ProcessTask(ctx context.Context, task *queue.Task) error {
	taskID, ok := kitContext.GetTaskID(ctx)
	if !ok {
		logger.WithContext(ctx).Errorf("评分队列获取任务ID失败")
		return nil
	}
	logger.WithContext(ctx).Infof("评分队列任务开始执行: %s", taskID)
	return nil
}
