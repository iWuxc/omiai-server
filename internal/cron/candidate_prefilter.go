package cron

import (
	"context"
	"encoding/json"
	"math/rand"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/iWuxc/go-wit/redis"
)

type CandidatePreFilterService struct {
	db *data.DB
}

func NewCandidatePreFilterService(db *data.DB) *CandidatePreFilterService {
	return &CandidatePreFilterService{db: db}
}

func (s *CandidatePreFilterService) JobName() string {
	return "CandidatePreFilterService"
}

func (s *CandidatePreFilterService) Schedule() string {
	// Daily at 2:00 AM
	return "0 0 2 * * *"
}

func (s *CandidatePreFilterService) Run() {
	ctx := context.WithValue(context.Background(), "request_id", uuid.NewString())

	lockKey := "lock:CandidatePreFilterService"
	lockRet := redis.GetRedis().GetClient().SetNX(ctx, lockKey, 1, time.Minute*30)
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
	s.Execute(ctx)
}

func (s *CandidatePreFilterService) Execute(ctx context.Context) {
	// 1. Get all single clients
	var clients []*biz_omiai.Client
	if err := s.db.WithContext(ctx).Where("status = ?", biz_omiai.ClientStatusSingle).Find(&clients).Error; err != nil {
		log.WithContext(ctx).Errorf("Failed to fetch clients: %v", err)
		return
	}

	// 2. Pre-calculate candidates for each client
	for _, client := range clients {
		if client.Age == 0 {
			client.Age = client.RealAge()
		}

		targetGender := 1
		if client.Gender == 1 {
			targetGender = 2
		}

		// Find potential matches (Status Single, Opposite Gender)
		var potentialMatches []*biz_omiai.Client
		if err := s.db.WithContext(ctx).Where("gender = ? AND status = ?", targetGender, biz_omiai.ClientStatusSingle).
			Limit(100).Find(&potentialMatches).Error; err != nil {
			log.WithContext(ctx).Errorf("Failed to fetch matches for client %d: %v", client.ID, err)
			continue
		}

		var candidates []*biz_omiai.Candidate
		for _, match := range potentialMatches {
			if match.Age == 0 {
				match.Age = match.RealAge()
			}

			// Scoring Logic
			score := 60
			tags := []string{}

			// 1. Age Score (Max 20)
			ageGap := abs(match.Age - client.Age)
			if ageGap <= 3 {
				score += 20
				tags = append(tags, "年龄相仿")
			} else if ageGap <= 5 {
				score += 10
			} else if ageGap > 10 {
				score -= 10
			}

			// 2. Education Score (Max 15)
			if match.Education >= client.Education {
				score += 15
				tags = append(tags, "学历相当")
			} else if match.Education < client.Education-1 {
				score -= 5
			}

			// 3. Height Score (Max 5)
			if client.Gender == 1 { // Male Client prefers Female < Male
				if match.Height > 0 && client.Height > 0 && match.Height < client.Height && match.Height > client.Height-20 {
					score += 5
				}
			} else { // Female Client prefers Male > Female
				if match.Height > 0 && client.Height > 0 && match.Height > client.Height {
					score += 5
				}
			}

			// 4. Random Jitter (0-5)
			score += rand.Intn(6)

			// Cap score
			if score > 99 {
				score = 99
			}
			if score < 40 {
				score = 40
			}

			if len(tags) == 0 {
				tags = append(tags, "缘分推荐")
			}

			candidates = append(candidates, &biz_omiai.Candidate{
				CandidateID: match.ID,
				Name:        match.Name,
				Avatar:      match.Avatar,
				MatchScore:  score,
				Tags:        tags,
				Age:         match.Age,
				Height:      match.Height,
				Education:   int(match.Education),
			})
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].MatchScore > candidates[j].MatchScore
		})
		if len(candidates) > 20 {
			candidates = candidates[:20]
		}

		if bytes, err := json.Marshal(candidates); err == nil {
			if err := s.db.WithContext(ctx).Model(client).Update("candidate_cache_json", string(bytes)).Error; err != nil {
				log.WithContext(ctx).Errorf("Failed to update cache for client %d: %v", client.ID, err)
			}
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
