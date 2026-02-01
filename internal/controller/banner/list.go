package banner

import (
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) List(ctx *gin.Context) {
	var req validates.BannerListValidate
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	clause := &biz.WhereClause{
		OrderBy: "sort_order desc",
		Where:   "status = ?",
		Args:    []interface{}{biz_omiai.BannerStatusEnable},
	}

	bannerList, err := c.Banner.Select(ctx, clause, []string{"id", "title", "image_url", "status", "link_url"}, offset, req.PageSize)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取轮播图列表失败")
		return
	}
	bannerResponseList := make([]*BannerResponse, 0)
	for _, banner := range bannerList {
		bannerResponseList = append(bannerResponseList, &BannerResponse{
			ID:       banner.ID,
			Title:    banner.Title,
			ImageURL: banner.ImageURL,
			LinkUrl:  banner.LinkUrl,
		})
	}
	response.SuccessResponse(ctx, "ok", map[string]interface{}{
		"list": bannerResponseList,
	})

}
