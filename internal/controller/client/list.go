package client

import (
	"omiai-server/internal/biz"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) List(ctx *gin.Context) {
	var req validates.ClientListValidate
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
		OrderBy: "created_at desc",
		Where:   "1=1",
		Args:    []interface{}{},
	}

	if req.Name != "" {
		clause.Where += " AND name LIKE ?"
		clause.Args = append(clause.Args, "%"+req.Name+"%")
	}
	if req.Phone != "" {
		clause.Where += " AND phone LIKE ?"
		clause.Args = append(clause.Args, "%"+req.Phone+"%")
	}
	if req.Gender != 0 {
		clause.Where += " AND gender = ?"
		clause.Args = append(clause.Args, req.Gender)
	}

	// Range Filters
	if req.MinHeight > 0 {
		clause.Where += " AND height >= ?"
		clause.Args = append(clause.Args, req.MinHeight)
	}
	if req.MaxHeight > 0 {
		clause.Where += " AND height <= ?"
		clause.Args = append(clause.Args, req.MaxHeight)
	}
	if req.MinIncome > 0 {
		clause.Where += " AND income >= ?"
		clause.Args = append(clause.Args, req.MinIncome)
	}
	if req.Education > 0 {
		clause.Where += " AND education >= ?" // Assuming higher value = higher education
		clause.Args = append(clause.Args, req.Education)
	}
	if req.Address != "" {
		clause.Where += " AND address LIKE ?"
		clause.Args = append(clause.Args, "%"+req.Address+"%")
	}
	if req.Profession != "" {
		clause.Where += " AND profession LIKE ?"
		clause.Args = append(clause.Args, "%"+req.Profession+"%")
	}
	if req.MaritalStatus > 0 {
		clause.Where += " AND marital_status = ?"
		clause.Args = append(clause.Args, req.MaritalStatus)
	}
	if req.HouseStatus > 0 {
		clause.Where += " AND house_status = ?"
		clause.Args = append(clause.Args, req.HouseStatus)
	}
	if req.CarStatus > 0 {
		clause.Where += " AND car_status = ?"
		clause.Args = append(clause.Args, req.CarStatus)
	}

	// Phase 1: 状态筛选
	if req.Status > 0 {
		clause.Where += " AND status = ?"
		clause.Args = append(clause.Args, req.Status)
	}

	// 单人模式：移除 Scope 权限过滤，默认返回所有客户
	// 原公海池逻辑废弃，所有录入数据均可见
	/*
		currentUserID := ctx.GetUint64("current_user_id")
		switch req.Scope {
		case "my":
			clause.Where += " AND manager_id = ?"
			clause.Args = append(clause.Args, currentUserID)
		case "public":
			clause.Where += " AND is_public = 1"
		default:
		}
	*/

	// Phase 1: 标签筛选 (JSON 数组包含)
	// MySQL 5.7+ 支持 JSON_CONTAINS(tags, '"tag_name"')
	// 这里假设 tags 存的是 ["tag1", "tag2"] 字符串
	if req.Tags != "" {
		// 简单实现：LIKE
		clause.Where += " AND tags LIKE ?"
		clause.Args = append(clause.Args, "%"+req.Tags+"%")
	}

	// Age Filter (Birthday based)
	now := time.Now()
	if req.MinAge > 0 {
		// MinAge 25 means born BEFORE (Now - 25 years)
		// e.g. 2023 - 25 = 1998. Born in 1998 is 25. Born in 1997 is 26.
		// So birthday <= 1998-MM-DD
		targetDate := now.AddDate(-req.MinAge, 0, 0).Format("2006-01-02")
		clause.Where += " AND birthday <= ?"
		clause.Args = append(clause.Args, targetDate)
	}
	if req.MaxAge > 0 {
		// MaxAge 30 means born AFTER (Now - 30 years - 1 year?)
		// e.g. 2023 - 30 = 1993. Born in 1993 is 30. Born in 1992 is 31.
		// So birthday >= 1993-01-01 (approx)
		// Actually, simpler: Age = Year(Now) - Year(Birth).
		// MaxAge 30 -> Year(Birth) >= Year(Now) - 30
		targetDate := now.AddDate(-req.MaxAge-1, 0, 0).Format("2006-01-02")
		clause.Where += " AND birthday > ?"
		clause.Args = append(clause.Args, targetDate)
	}

	list, err := c.client.Select(ctx, clause, nil, offset, req.PageSize)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户列表失败")
		return
	}

	respList := make([]*ClientResponse, 0)
	for _, v := range list {
		respList = append(respList, &ClientResponse{
			ID:                  v.ID,
			Name:                v.Name,
			Gender:              v.Gender,
			Phone:               v.Phone,
			Birthday:            v.Birthday,
			Age:                 CalculateAge(v.Birthday),
			Avatar:              "https://api.dicebear.com/7.x/avataaars/svg?seed=" + v.Name, // Mock avatar
			Zodiac:              v.Zodiac,
			Height:              v.Height,
			Weight:              v.Weight,
			Education:           v.Education,
			MaritalStatus:       v.MaritalStatus,
			Address:             v.Address,
			FamilyDescription:   v.FamilyDescription,
			Income:              v.Income,
			Profession:          v.Profession,
			HouseStatus:         v.HouseStatus,
			HouseAddress:        v.HouseAddress,
			CarStatus:           v.CarStatus,
			PartnerRequirements: v.PartnerRequirements,
			Remark:              v.Remark,
			Photos:              v.Photos,
			// Phase 1 Response
			ManagerID: v.ManagerID,
			IsPublic:  v.IsPublic,
			Tags:      v.Tags,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	response.SuccessResponse(ctx, "ok", map[string]interface{}{
		"list": respList,
	})
}
