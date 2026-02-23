package template

import (
	"strconv"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo biz_omiai.TemplateRepo
}

func NewController(repo biz_omiai.TemplateRepo) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) Create(ctx *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}

	template := &biz_omiai.CommunicationTemplate{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
	}

	if err := c.repo.Create(template); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建失败")
		return
	}
	response.SuccessResponse(ctx, "创建成功", template)
}

func (c *Controller) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	category := ctx.Query("category")

	list, total, err := c.repo.List(category, page, pageSize)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取列表失败")
		return
	}

	response.SuccessResponse(ctx, "获取成功", map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

func (c *Controller) Update(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}

	template, err := c.repo.Get(id)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "模板不存在")
		return
	}

	if req.Title != "" {
		template.Title = req.Title
	}
	if req.Content != "" {
		template.Content = req.Content
	}
	if req.Category != "" {
		template.Category = req.Category
	}

	if err := c.repo.Update(template); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新失败")
		return
	}
	response.SuccessResponse(ctx, "更新成功", template)
}

func (c *Controller) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := c.repo.Delete(id); err != nil {
		response.ErrorResponse(ctx, response.DBDeleteCommonError, "删除失败")
		return
	}
	response.SuccessResponse(ctx, "删除成功", nil)
}

// Use 记录使用并增加计数
func (c *Controller) Use(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := c.repo.IncrementUsage(id); err != nil {
		// 记录失败不影响主流程
	}
	response.SuccessResponse(ctx, "记录成功", nil)
}
