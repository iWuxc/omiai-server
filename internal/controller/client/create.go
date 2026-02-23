package client

import (
	"strings"
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
		WorkProvinceCode:    req.WorkProvinceCode,
		WorkCityCode:        req.WorkCityCode,
		WorkDistrictCode:    req.WorkDistrictCode,
		Position:            req.Position,
		ParentsProfession:   req.ParentsProfession,
		Tags:                req.Tags,
		HouseStatus:         req.HouseStatus,
		HouseAddress:        req.HouseAddress,
		HouseProvinceCode:   req.HouseProvinceCode,
		HouseCityCode:       req.HouseCityCode,
		HouseDistrictCode:   req.HouseDistrictCode,
		CarStatus:           req.CarStatus,
		PartnerRequirements: req.PartnerRequirements,
		Remark:              req.Remark,
		Photos:              req.Photos,
	}

	// 自动计算年龄
	client.Age = client.RealAge()

	// 自动生成标签 (Auto-Tagging)
	autoTags := []string{}
	// 收入标签
	if client.Income >= 30000 {
		autoTags = append(autoTags, "高收入")
	} else if client.Income >= 15000 {
		autoTags = append(autoTags, "中高收入")
	}
	// 学历标签
	if client.Education >= 3 { // 本科及以上
		autoTags = append(autoTags, "本科及以上")
	}
	if client.Education >= 4 { // 硕士及以上
		autoTags = append(autoTags, "硕博")
	}
	// 房车标签
	if client.HouseStatus == 2 || client.HouseStatus == 3 {
		autoTags = append(autoTags, "有房")
	}
	if client.CarStatus == 2 {
		autoTags = append(autoTags, "有车")
	}
	// 职业标签
	if client.WorkUnit != "" {
		// 简单关键词匹配
		if contains(client.WorkUnit, "公务员", "政府", "局", "委") {
			autoTags = append(autoTags, "体制内")
		} else if contains(client.WorkUnit, "银行", "证券", "金融") {
			autoTags = append(autoTags, "金融圈")
		} else if contains(client.WorkUnit, "学校", "大学", "中学", "小学") {
			autoTags = append(autoTags, "教师")
		} else if contains(client.WorkUnit, "医院") {
			autoTags = append(autoTags, "医生/护士")
		}
	}
	// 简单的海归判断 (需更复杂的逻辑，这里仅示例)
	if contains(client.Education, "海外", "留学") { // 假设 Education 字段能体现
		autoTags = append(autoTags, "海归")
	}

	// 将自动标签追加到原有 Tags 字段中 (假设 Tags 是 JSON 字符串或逗号分隔)
	// 这里简化处理，直接拼接到 Remark 或者专门的 Tags 字段
	// 由于 Tags 目前是 string，建议存储 JSON 数组
	// TODO: 需要引入 encoding/json 处理 Tags 字段的合并

	if err := c.client.Create(ctx, client); err != nil {
		log.WithContext(ctx).Errorf("Client Create failed: %v", err)
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建客户档案失败")
		return
	}

	response.SuccessResponse(ctx, "创建成功", client)
}

func contains(target interface{}, keywords ...string) bool {
	s := ""
	switch v := target.(type) {
	case string:
		s = v
	case int, int8, int64:
		// 数值类型无法包含关键词，返回 false
		return false
	}
	
	for _, k := range keywords {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}
