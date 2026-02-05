package client

import (
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Detail(ctx *gin.Context) {
	var req validates.ClientDetailValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	client, err := c.client.Get(ctx, req.ID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户详情失败")
		return
	}
	if client == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	var partnerID uint64
	if client.PartnerID != nil {
		partnerID = *client.PartnerID
	}

	resp := &ClientResponse{
		ID:                  client.ID,
		Name:                client.Name,
		Gender:              client.Gender,
		Phone:               client.Phone,
		Birthday:            client.Birthday,
		Age:                 CalculateAge(client.Birthday),
		Avatar:              client.Avatar,
		Zodiac:              client.Zodiac,
		Height:              client.Height,
		Weight:              client.Weight,
		Education:           client.Education,
		MaritalStatus:       client.MaritalStatus,
		Address:             client.Address,
		FamilyDescription:   client.FamilyDescription,
		Income:              client.Income,
		Profession:          client.Profession,
		HouseStatus:         client.HouseStatus,
		HouseAddress:        client.HouseAddress,
		CarStatus:           client.CarStatus,
		Status:              client.Status,
		PartnerID:           partnerID,
		PartnerRequirements: client.PartnerRequirements,
		Remark:              client.Remark,
		Photos:              client.Photos,
		CreatedAt:           client.CreatedAt,
		UpdatedAt:           client.UpdatedAt,
	}

	if resp.Avatar == "" {
		resp.Avatar = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + resp.Name
	}

	if partnerID > 0 {
		if partner, err := c.client.Get(ctx, partnerID); err == nil && partner != nil {
			resp.PartnerName = partner.Name
			resp.PartnerAvatar = partner.Avatar
			if resp.PartnerAvatar == "" {
				resp.PartnerAvatar = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + partner.Name
			}
		}
	}

	response.SuccessResponse(ctx, "ok", resp)
}
