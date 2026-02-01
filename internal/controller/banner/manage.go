package banner

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Create(ctx *gin.Context) {
	var req validates.BannerCreateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	banner := &biz_omiai.Banner{
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		LinkUrl:   req.LinkUrl,
	}

	if err := c.Banner.Create(ctx, banner); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建轮播图失败")
		return
	}
	response.SuccessResponse(ctx, "创建成功", banner)
}

func (c *Controller) Update(ctx *gin.Context) {
	var req validates.BannerUpdateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	banner := &biz_omiai.Banner{
		ID:        req.ID,
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		LinkUrl:   req.LinkUrl,
	}

	if err := c.Banner.Update(ctx, banner); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新轮播图失败")
		return
	}
	response.SuccessResponse(ctx, "更新成功", banner)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var req validates.BannerDeleteValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	if err := c.Banner.Delete(ctx, req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBDeleteCommonError, "删除轮播图失败")
		return
	}
	response.SuccessResponse(ctx, "删除成功", nil)
}
