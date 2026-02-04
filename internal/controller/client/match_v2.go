package client

import (
	"encoding/json"
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"sort"

	"github.com/gin-gonic/gin"
)

// PartnerRequirements defines the structure of the JSON field
type PartnerRequirements struct {
	MinAge        int    `json:"min_age"`
	MaxAge        int    `json:"max_age"`
	MinHeight     int    `json:"min_height"`
	MaxHeight     int    `json:"max_height"`
	MinIncome     int    `json:"min_income"`
	Education     int8   `json:"education"`
	MaritalStatus []int8 `json:"marital_status"` // Allow multiple: [1, 3] (Unmarried, Divorced)
	HouseStatus   int8   `json:"house_status"`   // Minimum requirement
}

// ScoredCandidate wraps the client with a match score
type ScoredCandidate struct {
	Client *ClientResponse
	Score  int
	Reason []string // Why it matched (or penalty reasons)
}

// MatchV2 implements the Smart Match V2.0 logic
// It uses both hard filtering (SQL) and soft scoring (Go)
func (c *Controller) MatchV2(ctx *gin.Context) {
	var req validates.ClientDetailValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	// 1. Get Source Client
	source, err := c.client.Get(ctx, req.ID)
	if err != nil || source == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户档案不存在")
		return
	}

	// 2. Parse Source Requirements
	var reqs PartnerRequirements
	if source.PartnerRequirements != "" {
		_ = json.Unmarshal([]byte(source.PartnerRequirements), &reqs)
	}
	// Set defaults if empty
	if reqs.MinAge == 0 {
		reqs.MinAge = 18
	}
	if reqs.MaxAge == 0 {
		reqs.MaxAge = 99
	}

	// 3. Build SQL Query (Hard Filters)
	// We only filter by Gender and Status initially to get a candidate pool
	// Detailed filtering happens in memory for better flexibility and scoring
	targetGender := 1
	if source.Gender == 1 {
		targetGender = 2
	}

	clause := &biz.WhereClause{
		Where: "gender = ? AND status = 1", // Only single candidates
		Args:  []interface{}{targetGender},
	}

	// Fetch candidates (limit 100 for performance, then score them)
	candidates, err := c.client.Select(ctx, clause, nil, 0, 100)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "匹配库查询失败")
		return
	}

	// 4. Scoring & Filtering
	var scoredList []ScoredCandidate
	sourceAge := CalculateAge(source.Birthday)

	for _, target := range candidates {
		score := 0
		reasons := []string{}
		isHardPass := true

		targetAge := CalculateAge(target.Birthday)

		// --- Hard Filter Check (Based on Source's Requirements) ---

		// Age Check
		if targetAge < reqs.MinAge || targetAge > reqs.MaxAge {
			// If strict, set isHardPass = false. For now, we penalize heavily instead of hiding
			score -= 50
			reasons = append(reasons, fmt.Sprintf("年龄不符(%d岁)", targetAge))
		} else {
			score += 20
		}

		// Height Check
		if reqs.MinHeight > 0 && target.Height < reqs.MinHeight {
			score -= 30
			reasons = append(reasons, fmt.Sprintf("身高不符(%dcm)", target.Height))
		} else if reqs.MaxHeight > 0 && target.Height > reqs.MaxHeight {
			score -= 10
		} else {
			score += 15
		}

		// Education Check
		if reqs.Education > 0 && target.Education < reqs.Education {
			score -= 20
			reasons = append(reasons, "学历未达标")
		} else {
			score += 10
		}

		// --- Reverse Match Check (Does Source meet Target's reqs?) ---
		var targetReqs PartnerRequirements
		if target.PartnerRequirements != "" {
			_ = json.Unmarshal([]byte(target.PartnerRequirements), &targetReqs)

			// Reverse Age
			if targetReqs.MinAge > 0 && (sourceAge < targetReqs.MinAge || sourceAge > targetReqs.MaxAge) {
				score -= 40
				reasons = append(reasons, "对方觉得你年龄不合适")
			} else {
				score += 10 // Mutual match bonus
			}

			// Reverse Height
			if targetReqs.MinHeight > 0 && source.Height < targetReqs.MinHeight {
				score -= 30
				reasons = append(reasons, "对方觉得你身高不合适")
			}
		}

		// --- Base Score ---
		score += 50 // Base points

		// Final threshold
		if score < 0 {
			score = 0
		}
		if isHardPass {
			scoredList = append(scoredList, ScoredCandidate{
				Client: convertToResponse(target),
				Score:  score,
				Reason: reasons,
			})
		}
	}

	// 5. Sort by Score DESC
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].Score > scoredList[j].Score
	})

	// 6. Return Top 20
	limit := 20
	if len(scoredList) < limit {
		limit = len(scoredList)
	}

	finalList := make([]map[string]interface{}, limit)
	for i := 0; i < limit; i++ {
		finalList[i] = map[string]interface{}{
			"client":     scoredList[i].Client,
			"score":      scoredList[i].Score,
			"match_tags": scoredList[i].Reason,
		}
	}

	response.SuccessResponse(ctx, "匹配成功", map[string]interface{}{
		"list":       finalList,
		"source_req": reqs, // Return used requirements for UI display
	})
}

func convertToResponse(v *biz_omiai.Client) *ClientResponse {
	return &ClientResponse{
		ID:                  v.ID,
		Name:                v.Name,
		Gender:              v.Gender,
		Phone:               v.Phone,
		Birthday:            v.Birthday,
		Age:                 CalculateAge(v.Birthday),
		Avatar:              v.Avatar,
		Zodiac:              v.Zodiac,
		Height:              v.Height,
		Weight:              v.Weight,
		Education:           v.Education,
		MaritalStatus:       v.MaritalStatus,
		Address:             v.Address,
		FamilyDescription:   v.FamilyDescription,
		Income:              v.Income,
		Profession:          v.Profession,
		HouseStatus:         v.HouseStatus,
		HouseAddress:        v.HouseAddress,
		CarStatus:           v.CarStatus,
		PartnerRequirements: v.PartnerRequirements,
		Remark:              v.Remark,
		Photos:              v.Photos,
		ManagerID:           v.ManagerID,
		IsPublic:            v.IsPublic,
		Tags:                v.Tags,
		CreatedAt:           v.CreatedAt,
		UpdatedAt:           v.UpdatedAt,
	}
}
