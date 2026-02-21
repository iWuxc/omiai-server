package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

func (c *Controller) Update(ctx *gin.Context) {
	var req validates.ClientUpdateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("Client Update validation failed: %v", err)
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	log.Infof("Updating client ID: %d", req.ID)

	// 先获取现有数据，或者直接更新字段
	// 这里我们构造一个 Client 对象，只包含需要更新的字段
	// 注意：GORM 的 Update 行为取决于实现，这里假设传入的 struct 字段会被更新
	client := &biz_omiai.Client{
		ID:                  req.ID,
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

	// 重新计算年龄
	client.Age = client.RealAge()

	if err := c.client.Update(ctx, client); err != nil {
		log.Errorf("Failed to update client: %v", err)
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新客户档案失败")
		return
	}

	response.SuccessResponse(ctx, "更新成功", client)
}
