package china_region

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 地区控制器
type Controller struct {
	regionRepo biz_omiai.ChinaRegionInterface
}

func NewController(regionRepo biz_omiai.ChinaRegionInterface) *Controller {
	return &Controller{regionRepo: regionRepo}
}

// GetProvinces 获取所有省份
func (c *Controller) GetProvinces(ctx *gin.Context) {
	provinces, err := c.regionRepo.GetProvinces()
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取省份列表失败")
		return
	}
	response.SuccessResponse(ctx, "ok", provinces)
}

// GetCitiesByProvince 获取某省下的城市
func (c *Controller) GetCitiesByProvince(ctx *gin.Context) {
	provinceCode := ctx.Param("code")
	if provinceCode == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "省份编码不能为空")
		return
	}

	cities, err := c.regionRepo.GetCitiesByProvince(provinceCode)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取城市列表失败")
		return
	}
	response.SuccessResponse(ctx, "ok", cities)
}

// GetDistrictsByCity 获取某城市下的区县
func (c *Controller) GetDistrictsByCity(ctx *gin.Context) {
	cityCode := ctx.Param("code")
	if cityCode == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "城市编码不能为空")
		return
	}

	districts, err := c.regionRepo.GetDistrictsByCity(cityCode)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取区县列表失败")
		return
	}
	response.SuccessResponse(ctx, "ok", districts)
}

// GetHotCities 获取热门城市
func (c *Controller) GetHotCities(ctx *gin.Context) {
	cities, err := c.regionRepo.GetHotCities()
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取热门城市失败")
		return
	}
	response.SuccessResponse(ctx, "ok", cities)
}

// SearchRegions 搜索地区
func (c *Controller) SearchRegions(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "搜索关键词不能为空")
		return
	}

	// 优先按名称搜索
	regions, err := c.regionRepo.SearchByKeyword(keyword)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "搜索失败")
		return
	}

	// 如果名称搜索无结果，尝试拼音搜索
	if len(regions) == 0 {
		regions, err = c.regionRepo.SearchByPinyin(keyword)
		if err != nil {
			response.ErrorResponse(ctx, response.DBSelectCommonError, "搜索失败")
			return
		}
	}

	response.SuccessResponse(ctx, "ok", regions)
}

// GetRegionDetail 获取地区详情（包含完整路径）
func (c *Controller) GetRegionDetail(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "地区编码不能为空")
		return
	}

	// 获取当前地区
	region, err := c.regionRepo.GetByCode(code)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取地区详情失败")
		return
	}

	// 获取完整路径
	fullPath, err := c.regionRepo.GetFullPath(code)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取地区路径失败")
		return
	}

	response.SuccessResponse(ctx, "ok", gin.H{
		"region":    region,
		"full_path": fullPath,
	})
}
