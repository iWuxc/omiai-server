package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Create(ctx *gin.Context) {
	var req validates.ClientCreateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	client := &biz_omiai.Client{
		Name:                req.Name,
		Gender:              req.Gender,
		Phone:               req.Phone,
		Birthday:            req.Birthday,
		Zodiac:              req.Zodiac,
		Height:              req.Height,
		Weight:              req.Weight,
		Education:           req.Education,
		MaritalStatus:       req.MaritalStatus,
		Address:             req.Address,
		FamilyDescription:   req.FamilyDescription,
		Income:              req.Income,
		Profession:          req.Profession,
		HouseStatus:         req.HouseStatus,
		CarStatus:           req.CarStatus,
		PartnerRequirements: req.PartnerRequirements,
		Remark:              req.Remark,
		Photos:              req.Photos,
	}

	if err := c.Client.Create(ctx, client); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建客户档案失败")
		return
	}

	response.SuccessResponse(ctx, "创建成功", client)
}
