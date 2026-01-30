package client

import (
	"fmt"
	"omiai-server/internal/biz"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Match(ctx *gin.Context) {
	var req validates.ClientDetailValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	// 1. Get the Source Client
	client, err := c.Client.Get(ctx, req.ID)
	if err != nil || client == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	// 2. Build Matching Criteria (Basic Rule-based)
	// Default Rule: Opposite Gender
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

	// TODO: Parse client.PartnerRequirements to add more filters (Age, Height, Income)
	// For now, we return the list of opposite gender candidates.
	// We can limit to 20 candidates.
	
	list, err := c.Client.Select(ctx, clause, nil, 0, 20)
	if err != nil {
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
