package tenant

import (
	biz_tenant "omiai-server/internal/biz/tenant"
	"omiai-server/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Tenant biz_tenant.TenantInterface
}

func NewController(tenant biz_tenant.TenantInterface) *Controller {
	return &Controller{Tenant: tenant}
}

type CreateTenantRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	SubscribeMonths int `json:"subscribe_months" binding:"required"`
}

// Create 平台超管创建新入驻的机构(租户)
func (c *Controller) Create(ctx *gin.Context) {
	var req CreateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	tenant := &biz_tenant.Tenant{
		Name:         req.Name,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       1,
		ExpireAt:     time.Now().AddDate(0, req.SubscribeMonths, 0),
	}

	if err := c.Tenant.Create(ctx, tenant); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建机构失败, 可能名称已存在")
		return
	}

	response.SuccessResponse(ctx, "机构入驻成功", tenant)
}

// List 获取所有入驻机构列表
func (c *Controller) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	list, total, err := c.Tenant.List(ctx, offset, size)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取机构列表失败")
		return
	}

	response.SuccessResponse(ctx, "success", map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
