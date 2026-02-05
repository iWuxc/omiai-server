package match

import (
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Create(ctx *gin.Context) {
	var req validates.MatchCreateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. Gender Validation
	male, err := c.client.Get(ctx, req.MaleClientID)
	if err != nil || male == nil || male.Gender != 1 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "男方信息错误或性别不符")
		return
	}

	female, err := c.client.Get(ctx, req.FemaleClientID)
	if err != nil || female == nil || female.Gender != 2 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "女方信息错误或性别不符")
		return
	}

	// 2. Status Validation
	if male.Status == biz_omiai.ClientStatusMatched || female.Status == biz_omiai.ClientStatusMatched {
		response.ErrorResponse(ctx, response.ParamsCommonError, "其中一方已有匹配对象")
		return
	}

	matchDate := req.MatchDate
	if matchDate.IsZero() {
		matchDate = time.Now()
	}

	record := &biz_omiai.MatchRecord{
		MaleClientID:   req.MaleClientID,
		FemaleClientID: req.FemaleClientID,
		MatchDate:      matchDate,
		MatchScore:     req.MatchScore,
		Status:         biz_omiai.MatchStatusAcquaintance, // Default to Acquaintance
		Remark:         req.Remark,
	}

	if err := c.match.Create(ctx, record); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "保存匹配记录失败")
		return
	}

	response.SuccessResponse(ctx, "匹配成功", record)
}

func (c *Controller) List(ctx *gin.Context) {
	var req validates.MatchListValidate
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	where := "1=1"
	args := []interface{}{}

	if req.Status > 0 {
		where += " AND status = ?"
		args = append(args, req.Status)
	}

	list, err := c.match.Select(ctx, &biz.WhereClause{
		Where:   where,
		Args:    args,
		OrderBy: "match_date desc",
	}, req.Offset(), req.Limit())

	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "查询失败")
		return
	}

	response.SuccessResponse(ctx, "ok", list)
}

func (c *Controller) Dissolve(ctx *gin.Context) {
	var req validates.DissolveMatchValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// Get Operator Name from Context
	operator := "Admin"
	if userID, ok := ctx.Get("user_id"); ok {
		// If user_id is stored as uint64 or similar, we might want to fetch user details.
		// For simplicity, we can use the ID or try to fetch user if needed.
		// Assuming we want to record the actual user name if possible.
		if id, ok := userID.(uint64); ok {
			if user, err := c.user.GetByID(ctx, id); err == nil && user != nil {
				operator = user.Nickname
			} else {
				operator = fmt.Sprintf("User:%d", id)
			}
		} else {
			operator = fmt.Sprintf("%v", userID)
		}
	}

	if err := c.match.DissolveMatch(ctx, req.ClientID, operator, req.Reason); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, err.Error())
		return
	}

	response.SuccessResponse(ctx, "解除匹配成功", nil)
}

func (c *Controller) UpdateStatus(ctx *gin.Context) {
	var req validates.MatchUpdateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	record, err := c.match.Get(ctx, req.ID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "记录不存在")
		return
	}

	// Permission Check
	role, exists := ctx.Get("role")
	if !exists {
		response.ErrorResponse(ctx, response.ParamsCommonError, "未授权操作")
		return
	}
	roleStr, ok := role.(string)
	if !ok || (roleStr != biz_omiai.RoleAdmin && roleStr != biz_omiai.RoleOperator) {
		response.ErrorResponse(ctx, response.ParamsCommonError, "权限不足")
		return
	}

	// Get Operator Name
	operator := "System"
	if userID, ok := ctx.Get("user_id"); ok {
		if id, ok := userID.(uint64); ok {
			user, err := c.user.GetByID(ctx, id)
			if err == nil && user != nil {
				operator = user.Nickname
			}
		}
	}
	// Fallback to request operator if system/unknown (optional, but let's prefer real user)
	if operator == "System" && req.Operator != "" {
		operator = req.Operator
	}

	if err := c.match.UpdateStatus(ctx, record.ID, record.Status, req.Status, operator, req.Reason); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新状态失败")
		return
	}

	response.SuccessResponse(ctx, "更新成功", nil)
}

func (c *Controller) GetStatusHistory(ctx *gin.Context) {
	idStr := ctx.Query("match_record_id")
	if idStr == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "参数错误")
		return
	}
	var id uint64
	fmt.Sscanf(idStr, "%d", &id)

	list, err := c.match.GetStatusHistory(ctx, id)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "查询失败")
		return
	}
	response.SuccessResponse(ctx, "ok", list)
}

func (c *Controller) Stats(ctx *gin.Context) {
	stats, err := c.match.Stats(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "统计失败")
		return
	}
	response.SuccessResponse(ctx, "ok", stats)
}
