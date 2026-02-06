package omiai

import (
	"math"

	biz_omiai "omiai-server/internal/biz/omiai"
)

// MatchCalculator 匹配度计算器
type MatchCalculator struct {
	client    *biz_omiai.Client
	candidate *biz_omiai.Client
}

// NewMatchCalculator 创建匹配计算器
func NewMatchCalculator(client, candidate *biz_omiai.Client) *MatchCalculator {
	return &MatchCalculator{
		client:    client,
		candidate: candidate,
	}
}

// Calculate 计算匹配度总分（0-100）
func (m *MatchCalculator) Calculate() int {
	score := 50.0 // 基础分50分

	// 1. 基础条件匹配度（权重40%）
	score += m.basicMatchScore() * 0.40

	// 2. 年龄匹配度（权重20%）
	score += m.ageMatchScore() * 0.20

	// 3. 学历匹配度（权重15%）
	score += m.educationMatchScore() * 0.15

	// 4. 收入匹配度（权重15%）
	score += m.incomeMatchScore() * 0.15

	// 5. 房车条件匹配度（权重10%）
	score += m.assetMatchScore() * 0.10

	// 确保在0-100范围内
	return int(math.Max(0, math.Min(100, score)))
}

// basicMatchScore 基础条件匹配度（身高、婚姻状况）
func (m *MatchCalculator) basicMatchScore() float64 {
	score := 50.0

	// 身高匹配（传统观念：男高女低合适）
	if m.client.Gender == 1 && m.candidate.Gender == 2 {
		// 男找女：男方身高应 > 女方身高
		if m.client.Height > m.candidate.Height {
			score += 20
			diff := m.client.Height - m.candidate.Height
			if diff >= 10 && diff <= 20 {
				score += 10 // 最佳身高差10-20cm
			}
		} else {
			score -= 10
		}
	} else if m.client.Gender == 2 && m.candidate.Gender == 1 {
		// 女找男：男方身高应 > 女方身高
		if m.candidate.Height > m.client.Height {
			score += 20
			diff := m.candidate.Height - m.client.Height
			if diff >= 10 && diff <= 20 {
				score += 10 // 最佳身高差10-20cm
			}
		} else {
			score -= 10
		}
	}

	// 婚姻状况匹配
	if m.client.MaritalStatus == m.candidate.MaritalStatus {
		score += 15 // 相同婚姻状况加分（未婚配未婚，离异配离异）
	} else if m.client.MaritalStatus == 1 || m.candidate.MaritalStatus == 1 {
		// 其中一方未婚，稍微减分（观念差异）
		score -= 5
	}

	return score
}

// ageMatchScore 年龄匹配度
func (m *MatchCalculator) ageMatchScore() float64 {
	clientAge := m.client.RealAge()
	candidateAge := m.candidate.RealAge()

	if clientAge == 0 || candidateAge == 0 {
		return 50.0 // 无法计算年龄，给中等分
	}

	ageDiff := math.Abs(float64(clientAge - candidateAge))
	score := 50.0

	// 传统观念：男大女小
	if m.client.Gender == 1 && m.candidate.Gender == 2 {
		// 男找女：男方应比女方大
		if clientAge > candidateAge {
			score += 20
			if ageDiff >= 2 && ageDiff <= 5 {
				score += 20 // 最佳年龄差2-5岁
			} else if ageDiff > 8 {
				score -= 10 // 年龄差太大
			}
		} else {
			score -= 15 // 女大男小，传统观念减分
		}
	} else if m.client.Gender == 2 && m.candidate.Gender == 1 {
		// 女找男：男方应比女方大
		if candidateAge > clientAge {
			score += 20
			if ageDiff >= 2 && ageDiff <= 8 {
				score += 20 // 最佳年龄差2-8岁
			} else if ageDiff > 12 {
				score -= 10 // 年龄差太大
			}
		} else {
			score -= 10
		}
	}

	return score
}

// educationMatchScore 学历匹配度
func (m *MatchCalculator) educationMatchScore() float64 {
	score := 50.0

	if m.client.Education == 0 || m.candidate.Education == 0 {
		return score
	}

	eduDiff := math.Abs(float64(m.client.Education - m.candidate.Education))

	if eduDiff == 0 {
		score += 30 // 学历相同，最佳匹配
	} else if eduDiff == 1 {
		score += 15 // 学历差一级，可接受
	} else if eduDiff >= 3 {
		score -= 20 // 学历差距大
	}

	return score
}

// incomeMatchScore 收入匹配度
func (m *MatchCalculator) incomeMatchScore() float64 {
	score := 50.0

	if m.client.Income == 0 || m.candidate.Income == 0 {
		return score
	}

	// 收入差距比例
	maxIncome := math.Max(float64(m.client.Income), float64(m.candidate.Income))
	minIncome := math.Min(float64(m.client.Income), float64(m.candidate.Income))
	gapRatio := (maxIncome - minIncome) / maxIncome

	if gapRatio <= 0.3 {
		score += 30 // 收入差距30%以内，很匹配
	} else if gapRatio <= 0.5 {
		score += 15 // 收入差距50%以内，可接受
	} else if gapRatio <= 0.8 {
		score -= 10 // 收入差距较大
	} else {
		score -= 25 // 收入差距太大
	}

	return score
}

// assetMatchScore 房车条件匹配度
func (m *MatchCalculator) assetMatchScore() float64 {
	score := 50.0

	// 房产情况匹配
	if m.client.HouseStatus >= 2 && m.candidate.HouseStatus >= 2 {
		score += 15 // 双方都有房，很稳定
	} else if m.client.HouseStatus >= 2 || m.candidate.HouseStatus >= 2 {
		score += 8 // 一方有房，可接受
	}

	// 车辆情况匹配
	if m.client.CarStatus == 2 && m.candidate.CarStatus == 2 {
		score += 10 // 双方都有车
	} else if m.client.CarStatus == 2 || m.candidate.CarStatus == 2 {
		score += 5 // 一方有车
	}

	return score
}

// GetMatchTags 获取匹配标签
func (m *MatchCalculator) GetMatchTags() []string {
	var tags []string

	// 年龄标签
	ageDiff := math.Abs(float64(m.client.RealAge() - m.candidate.RealAge()))
	if ageDiff <= 3 {
		tags = append(tags, "年龄相仿")
	}

	// 学历标签
	if m.client.Education == m.candidate.Education {
		tags = append(tags, "学历相当")
	}

	// 收入标签
	if m.client.Income > 0 && m.candidate.Income > 0 {
		maxIncome := math.Max(float64(m.client.Income), float64(m.candidate.Income))
		minIncome := math.Min(float64(m.client.Income), float64(m.candidate.Income))
		if (maxIncome-minIncome)/maxIncome <= 0.3 {
			tags = append(tags, "收入匹配")
		}
	}

	// 房车标签
	if m.client.HouseStatus >= 2 && m.candidate.HouseStatus >= 2 {
		tags = append(tags, "都有房产")
	}

	// 婚姻状况标签
	if m.client.MaritalStatus == m.candidate.MaritalStatus {
		tags = append(tags, "婚史相同")
	}

	// 默认标签
	if len(tags) == 0 {
		tags = append(tags, "可以尝试")
	}

	return tags
}

// GetMatchLevel 获取匹配等级
func GetMatchLevel(score int) string {
	if score >= 85 {
		return "perfect" // 完美匹配
	} else if score >= 70 {
		return "good" // 良好匹配
	} else if score >= 55 {
		return "average" // 一般匹配
	} else {
		return "poor" // 不太匹配
	}
}

// GetLevelText 获取匹配等级文本
func GetLevelText(level string) string {
	switch level {
	case "perfect":
		return "完美匹配"
	case "good":
		return "非常合适"
	case "average":
		return "可以尝试"
	default:
		return "不太合适"
	}
}

// GetLevelColor 获取匹配等级颜色
func GetLevelColor(level string) string {
	switch level {
	case "perfect":
		return "#52c41a" // 绿色
	case "good":
		return "#1890ff" // 蓝色
	case "average":
		return "#faad14" // 橙色
	default:
		return "#ff4d4f" // 红色
	}
}
