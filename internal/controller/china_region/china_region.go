package china_region

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/data/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo biz_omiai.ChinaRegionInterface
}

func NewController(db *data.DB) *Controller {
	return &Controller{
		repo: omiai.NewChinaRegionRepo(db),
	}
}

// GetProvinces 获取所有省份
func (c *Controller) GetProvinces(ctx *gin.Context) {
	regions, err := c.repo.GetProvinces()
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取省份失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", regions)
}

// GetCities 获取某省下的城市
func (c *Controller) GetCities(ctx *gin.Context) {
	provinceCode := ctx.Query("province_code")
	if provinceCode == "" {
		response.ErrorResponse(ctx, response.ValidateCommonError, "省份代码不能为空")
		return
	}

	regions, err := c.repo.GetCitiesByProvince(provinceCode)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取城市失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", regions)
}

// GetDistricts 获取某城市下的区县
func (c *Controller) GetDistricts(ctx *gin.Context) {
	cityCode := ctx.Query("city_code")
	if cityCode == "" {
		response.ErrorResponse(ctx, response.ValidateCommonError, "城市代码不能为空")
		return
	}

	regions, err := c.repo.GetDistrictsByCity(cityCode)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取区县失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", regions)
}

// GetHotCities 获取热门城市
func (c *Controller) GetHotCities(ctx *gin.Context) {
	regions, err := c.repo.GetHotCities()
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取热门城市失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", regions)
}

// Search 搜索地区 (关键词或拼音)
func (c *Controller) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		response.ErrorResponse(ctx, response.ValidateCommonError, "搜索关键词不能为空")
		return
	}

	// 简单的判断：如果包含字母则认为是拼音，否则是中文
	var regions []*biz_omiai.ChinaRegion
	var err error

	if isAlpha(keyword) {
		regions, err = c.repo.SearchByPinyin(keyword)
	} else {
		regions, err = c.repo.SearchByKeyword(keyword)
	}

	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "搜索失败")
		return
	}
	response.SuccessResponse(ctx, "搜索成功", regions)
}

// isAlpha 判断是否纯字母
func isAlpha(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
