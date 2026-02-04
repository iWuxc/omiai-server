package client

import (
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/service/chat_parser"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

type ImportAnalyzeRequest struct {
	Content string `json:"content" binding:"required"`
}

type ImportBatchRequest struct {
	List []chat_parser.ImportRecord `json:"list" binding:"required"`
}

// ImportAnalyze 接收文本，返回解析结果预览
func (c *Controller) ImportAnalyze(ctx *gin.Context) {
	var req ImportAnalyzeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	records, err := c.chatParserService.Parse(req.Content)
	if err != nil {
		response.ErrorResponse(ctx, response.ServiceCommonError, "解析失败: "+err.Error())
		return
	}

	response.SuccessResponse(ctx, "解析完成", map[string]interface{}{
		"total":   len(records),
		"records": records,
	})
}

// ImportBatch 批量入库
func (c *Controller) ImportBatch(ctx *gin.Context) {
	var req ImportBatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// currentUserID := ctx.GetUint64("current_user_id")

	successCount := 0
	failCount := 0
	errors := []string{}

	// Transaction support ideally, but for simplicity looping
	// In real world, use BatchInsert
	for _, record := range req.List {
		if record.ParseStatus == "error" {
			failCount++
			continue
		}

		// Deduplication check: Phone
		clause := &biz.WhereClause{
			Where: "phone = ?",
			Args:  []interface{}{record.Phone},
		}
		exists, _ := c.client.Select(ctx, clause, nil, 0, 1)
		if len(exists) > 0 {
			failCount++
			errors = append(errors, "手机号重复: "+record.Phone)
			continue
		}

		client := &biz_omiai.Client{
			Name:                record.Name,
			Gender:              record.Gender,
			Phone:               record.Phone,
			Birthday:            record.Birthday,
			Height:              record.Height,
			Weight:              record.Weight,
			Education:           record.Education,
			MaritalStatus:       record.MaritalStatus,
			Income:              record.Income,
			Address:             record.Address,
			Profession:          record.Profession,
			HouseStatus:         record.HouseStatus,
			CarStatus:           record.CarStatus,
			PartnerRequirements: record.PartnerRequirements,
			// ManagerID:           currentUserID, // Import to current user's private pool
			// IsPublic:            true, // Phase 1 adjustment: Default to public since single user mode
			Status: 1, // Default single
			Remark: "批量导入数据",
		}

		// Check if manager_id column exists (Dynamic adaptation for different DB versions)
		// Or just ignore it if we are in single user mode and DB hasn't been migrated
		// But error log shows "Unknown column manager_id", so we must fix the DB or the code.
		// Since user mentioned "Single User Mode" and we removed migration requirement,
		// we should remove ManagerID and IsPublic assignment here if they cause error,
		// OR ensure DB has columns.

		// Based on user's error log: "Unknown column 'manager_id'"
		// It means the Phase 1 migration SQL was NOT executed.
		// To fix this immediately without forcing user to run SQL, we should comment out these fields
		// or use a version of Client struct that matches current DB.

		// However, Client struct HAS ManagerID field (we added it). GORM tries to insert it.
		// We should revert Client struct changes OR force user to migrate DB.
		// Given the user instruction "I need to correct... resources belong only to me",
		// maybe we should revert the DB dependency.

		// BUT, the best way is to apply the migration SQL.
		// Since I cannot run SQL directly on running instance easily without knowing credentials/env fully,
		// I will try to "Hide" these fields from GORM if possible, OR tell user to run migration.

		// Wait, I have "Write" tool. I can create a migration tool or just update the code to NOT use these fields
		// if they are not critical.

		// Let's modify the code to NOT set ManagerID/IsPublic for now,
		// AND more importantly, we need to modify the Client struct tags to `-` if we want to skip them,
		// but that affects Read.

		// Correct approach: The error is "Unknown column". The column is missing in DB.
		// I should check internal/biz/omiai/client.go again.

		if err := c.client.Create(ctx, client); err != nil {
			log.Errorf("Import create failed: %v", err)
			failCount++
			errors = append(errors, "写入失败: "+record.Name)
		} else {
			successCount++
		}
	}

	response.SuccessResponse(ctx, "导入完成", map[string]interface{}{
		"success_count": successCount,
		"fail_count":    failCount,
		"errors":        errors,
	})
}
