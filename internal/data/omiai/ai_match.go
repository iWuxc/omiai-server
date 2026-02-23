package omiai

import (
	"time"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type AIMatchRepo struct {
	db *data.DB
}

func NewAIMatchRepo(db *data.DB) biz_omiai.AIMatchRepo {
	return &AIMatchRepo{db: db}
}

// GenerateDailyRecommendations 生成每日推荐
func (r *AIMatchRepo) GenerateDailyRecommendations() error {
	// 1. 获取所有 S/A 级活跃客户
	var activeClients []*biz_omiai.Client
	// 假设 Tags 包含 "S级" 或 "A级" 或者最近活跃
	// 这里简化为获取最近活跃的 50 个未婚客户
	if err := r.db.DB.Where("status = ? AND marital_status = ?", 1, 1).Order("updated_at desc").Limit(50).Find(&activeClients).Error; err != nil {
		return err
	}

	for _, client := range activeClients {
		// 2. 为每个客户寻找匹配对象 (简化逻辑：异性，年龄差5岁以内，同城)
		var candidates []*biz_omiai.Client
		targetGender := 1
		if client.Gender == 1 {
			targetGender = 2
		}

		// 简单的规则筛选
		query := r.db.DB.Where("gender = ? AND status = ? AND marital_status = ?", targetGender, 1, 1)
		
		// 年龄筛选
		if client.Age > 0 {
			minAge := client.Age - 5
			maxAge := client.Age + 5
			query = query.Where("age >= ? AND age <= ?", minAge, maxAge)
		}

		// 同城筛选 (如果已填写工作城市代码)
		if client.WorkCityCode != "" {
			query = query.Where("work_city_code = ?", client.WorkCityCode)
		}

		if err := query.Limit(5).Find(&candidates).Error; err != nil {
			continue
		}

		// 3. 保存推荐记录 (这里假设有一个 Recommendation 表，或者直接创建 Match 记录)
		// 暂时仅打印日志或更新某个状态，实际需创建 biz_omiai.Recommendation
		for _, candidate := range candidates {
			// 检查是否已存在匹配
			var count int64
			r.db.DB.Model(&biz_omiai.MatchRecord{}).Where(
				"(male_client_id = ? AND female_client_id = ?) OR (male_client_id = ? AND female_client_id = ?)", 
				client.ID, candidate.ID, candidate.ID, client.ID,
			).Count(&count)

			if count == 0 {
				// 创建待推荐记录
				match := &biz_omiai.MatchRecord{
					MaleClientID:   uint64(client.ID),
					FemaleClientID: uint64(candidate.ID),
					Status:         1, // 相识/待推荐
					Remark:         "AI每日推荐",
					AdminID:        "ai_bot",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				// 注意：这里需要根据性别正确分配 Male/Female ID
				if client.Gender == 2 { // 如果 client 是女性
					match.MaleClientID = uint64(candidate.ID)
					match.FemaleClientID = uint64(client.ID)
				}
				
				r.db.DB.Create(match)
			}
		}
	}
	return nil
}

func (r *AIMatchRepo) GetDailyStats() (map[string]interface{}, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	r.db.DB.Model(&biz_omiai.MatchRecord{}).Where("remark = ? AND DATE(created_at) = ?", "AI每日推荐", today).Count(&count)
	return map[string]interface{}{"daily_match_count": count}, nil
}
