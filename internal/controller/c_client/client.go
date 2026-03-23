package c_client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
}

func NewController(db *data.DB, client biz_omiai.ClientInterface) *Controller {
	return &Controller{
		db:     db,
		Client: client,
	}
}

// GetMine 获取当前用户的个人资料
func (c *Controller) GetMine(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	client, err := c.Client.Get(ctx, clientID.(uint64))
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取资料失败")
		return
	}

	// TODO: 可以使用 DTO 进行脱敏，防止返回 ManagerID 或内部 Remark
	response.SuccessResponse(ctx, "success", client)
}

type ProfileUpdateRequest struct {
	Name                string `json:"name"`
	Gender              int8   `json:"gender"`
	Phone               string `json:"phone"`
	Birthday            string `json:"birthday"`
	Avatar              string `json:"avatar"`
	Height              int    `json:"height"`
	Weight              int    `json:"weight"`
	Education           int8   `json:"education"`
	MaritalStatus       int8   `json:"marital_status"`
	Address             string `json:"address"`
	Income              int    `json:"income"`
	Profession          string `json:"profession"`
	WorkCity            string `json:"work_city"`
	PartnerRequirements string `json:"partner_requirements"`
	Photos              string `json:"photos"`
}

// UpdateMine 更新个人资料
func (c *Controller) UpdateMine(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	var req ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	client, err := c.Client.Get(ctx, clientID.(uint64))
	if err != nil || client == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "用户不存在")
		return
	}

	// 增量更新允许修改的字段
	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Gender > 0 {
		client.Gender = req.Gender
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Birthday != "" {
		client.Birthday = req.Birthday
	}
	if req.Avatar != "" {
		client.Avatar = req.Avatar
	}
	if req.Height > 0 {
		client.Height = req.Height
	}
	if req.Weight > 0 {
		client.Weight = req.Weight
	}
	if req.Education > 0 {
		client.Education = req.Education
	}
	if req.MaritalStatus > 0 {
		client.MaritalStatus = req.MaritalStatus
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.Income > 0 {
		client.Income = req.Income
	}
	if req.Profession != "" {
		client.Profession = req.Profession
	}
	if req.WorkCity != "" {
		client.WorkCity = req.WorkCity
	}
	if req.PartnerRequirements != "" {
		client.PartnerRequirements = req.PartnerRequirements
	}
	if req.Photos != "" {
		client.Photos = req.Photos
	}

	if err := c.Client.Update(ctx, client); err != nil {
		log.Errorf("Update CClient failed: %v", err)
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新资料失败")
		return
	}

	response.SuccessResponse(ctx, "更新成功", client)
}
