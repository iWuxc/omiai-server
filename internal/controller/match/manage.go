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
	male, err := c.Client.Get(ctx, req.MaleClientID)
	if err != nil || male == nil || male.Gender != 1 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "男方信息错误或性别不符")
		return
	}

	female, err := c.Client.Get(ctx, req.FemaleClientID)
	if err != nil || female == nil || female.Gender != 2 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "女方信息错误或性别不符")
		return
	}

	// 2. Status Validation (Optional: can they be rematched?)
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
		Status:         biz_omiai.MatchStatusMatched,
		Remark:         req.Remark,
	}

	if err := c.Match.Create(ctx, record); err != nil {
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

	// Note: Complex filtering by names might need Joins or specific SQL
	// For simplicity, we filter by status and pagination first
	list, err := c.Match.Select(ctx, &biz.WhereClause{
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

func (c *Controller) UpdateStatus(ctx *gin.Context) {
	var req validates.MatchUpdateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	record, err := c.Match.Get(ctx, req.ID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "记录不存在")
		return
	}

	record.Status = req.Status
	if req.Remark != "" {
		record.Remark = req.Remark
	}

	if err := c.Match.Update(ctx, record); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新失败")
		return
	}

	// If status changed to Broken, reset clients to single
	if req.Status == biz_omiai.MatchStatusBroken {
		// This should ideally be in a transaction in the Repo
		// For now, keep it simple
		fmt.Println("Handle broken status logic...")
	}

	response.SuccessResponse(ctx, "更新成功", nil)
}
