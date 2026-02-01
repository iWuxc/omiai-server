package client

import (
	"fmt"
	"omiai-server/internal/biz"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

func (c *Controller) Match(ctx *gin.Context) {
	var req validates.ClientDetailValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Errorf("Match binding failed: %v", err)
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	log.Infof("Smart Match requested for client ID: %d", req.ID)

	// 1. Get the Source Client
	client, err := c.Client.Get(ctx, req.ID)
	if err != nil || client == nil {
		log.Errorf("Client not found: %d", req.ID)
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	// 2. Build Matching Criteria
	targetGender := 0
	if client.Gender == 1 {
		targetGender = 2
	} else if client.Gender == 2 {
		targetGender = 1
	}

	clause := &biz.WhereClause{
		OrderBy: "created_at desc",
		Where:   "gender = ?",
		Args:    []interface{}{targetGender},
	}

	// Basic Intelligent Filtering: Age (+/- 5 years)
	clientAge := CalculateAge(client.Birthday)
	if clientAge > 0 {
		now := time.Now()
		// If client is male (1), usually looks for same age or younger
		// If client is female (2), usually looks for same age or older
		// Simple logic for demo: +/- 10 years
		minBirthYear := now.Year() - (clientAge + 10)
		maxBirthYear := now.Year() - (clientAge - 10)
		clause.Where += " AND birthday >= ? AND birthday <= ?"
		clause.Args = append(clause.Args, fmt.Sprintf("%d-01-01", minBirthYear), fmt.Sprintf("%d-12-31", maxBirthYear))
	}

	// Height Filter (+/- 20cm)
	if client.Height > 0 {
		clause.Where += " AND height >= ? AND height <= ?"
		clause.Args = append(clause.Args, client.Height-20, client.Height+20)
	}

	log.Infof("Matching clause: %s, args: %v", clause.Where, clause.Args)
	
	list, err := c.Client.Select(ctx, clause, nil, 0, 20)
	if err != nil {
		log.Errorf("Select candidates failed: %v", err)
		response.ErrorResponse(ctx, response.DBSelectCommonError, "匹配失败")
		return
	}

	respList := make([]*ClientResponse, 0)
	for _, v := range list {
		respList = append(respList, &ClientResponse{
			ID:                  v.ID,
			Name:                v.Name,
			Gender:              v.Gender,
			Phone:               v.Phone,
			Birthday:            v.Birthday,
			Age:                 CalculateAge(v.Birthday),
			Avatar:              "https://api.dicebear.com/7.x/avataaars/svg?seed=" + v.Name,
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
			CarStatus:           v.CarStatus,
			PartnerRequirements: v.PartnerRequirements,
			Remark:              v.Remark,
			Photos:              v.Photos,
			CreatedAt:           v.CreatedAt,
			UpdatedAt:           v.UpdatedAt,
		})
	}

	response.SuccessResponse(ctx, fmt.Sprintf("为您匹配到 %d 位嘉宾", len(respList)), map[string]interface{}{
		"list": respList,
	})
}
