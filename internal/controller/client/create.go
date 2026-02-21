package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

func (c *Controller) Create(ctx *gin.Context) {
	var req validates.ClientCreateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("Client Create validation failed: %v", err)
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	log.Infof("Creating client: %s, gender: %d", req.Name, req.Gender)
	client := &biz_omiai.Client{
		Name:                req.Name,
		Gender:              req.Gender,
		Phone:               req.Phone,
		Birthday:            req.Birthday,
		Avatar:              req.Avatar,
		Zodiac:              req.Zodiac,
		Height:              req.Height,
		Weight:              req.Weight,
		Education:           req.Education,
		MaritalStatus:       req.MaritalStatus,
		Address:             req.Address,
		FamilyDescription:   req.FamilyDescription,
		Income:              req.Income,
		Profession:          req.Profession,
		WorkUnit:            req.WorkUnit,
		WorkCity:            req.WorkCity,
		Position:            req.Position,
		ParentsProfession:   req.ParentsProfession,
		Tags:                req.Tags,
		HouseStatus:         req.HouseStatus,
		HouseAddress:        req.HouseAddress,
		CarStatus:           req.CarStatus,
		PartnerRequirements: req.PartnerRequirements,
		Remark:              req.Remark,
		Photos:              req.Photos,
	}

	// 自动计算年龄
	client.Age = client.RealAge()

	if err := c.client.Create(ctx, client); err != nil {
		log.WithContext(ctx).Errorf("Client Create failed: %v", err)
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建客户档案失败")
		return
	}

	response.SuccessResponse(ctx, "创建成功", client)
}
