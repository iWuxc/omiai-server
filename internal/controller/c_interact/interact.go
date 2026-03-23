package c_interact

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
	Match  biz_omiai.MatchInterface
}

func NewController(db *data.DB, client biz_omiai.ClientInterface, match biz_omiai.MatchInterface) *Controller {
	return &Controller{
		db:     db,
		Client: client,
		Match:  match,
	}
}

type LikeRequest struct {
	TargetClientID uint64 `json:"target_client_id" binding:"required"`
}

// Like C端用户心动(右滑)
func (c *Controller) Like(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}
	fromID := clientID.(uint64)

	var req LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}
	toID := req.TargetClientID

	if fromID == toID {
		response.ErrorResponse(ctx, response.ParamsCommonError, "不能对自己心动")
		return
	}

	// 1. 记录我喜欢对方
	interaction := &biz_omiai.ClientInteraction{
		FromClientID: fromID,
		ToClientID:   toID,
		ActionType:   2, // 2 = 单向心动
	}
	if err := c.Client.SaveInteraction(ctx, interaction); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "操作失败")
		return
	}

	// 2. 检查对方是否也喜欢我
	reverseInteraction, err := c.Client.GetInteraction(ctx, toID, fromID)
	if err != nil {
		log.Errorf("Check reverse interaction failed: %v", err)
	}

	isMatched := false
	if reverseInteraction != nil && reverseInteraction.ActionType >= 2 {
		// 互相心动！
		isMatched = true

		// 升级双方互动状态为互相心动
		interaction.ActionType = 3
		_ = c.Client.SaveInteraction(ctx, interaction)

		reverseInteraction.ActionType = 3
		_ = c.Client.SaveInteraction(ctx, reverseInteraction)

		// 触发业务核心逻辑：在后台创建一条“相识”状态的 MatchRecord
		me, _ := c.Client.Get(ctx, fromID)
		target, _ := c.Client.Get(ctx, toID)

		if me != nil && target != nil {
			var maleID, femaleID uint64
			if me.Gender == 1 {
				maleID = me.ID
				femaleID = target.ID
			} else {
				maleID = target.ID
				femaleID = me.ID
			}

			// 创建相识记录 (赋能红娘)
			matchRecord := &biz_omiai.MatchRecord{
				MaleClientID:   maleID,
				FemaleClientID: femaleID,
				Status:         biz_omiai.MatchStatusAcquaintance,
				MatchDate:      time.Now(),
				Remark:         "【系统自动生成】双方在小程序互相心动，请红娘尽快介入跟进！",
				AdminID:        "system",
			}

			if err := c.Match.Create(ctx, matchRecord); err != nil {
				log.Errorf("Auto create match record failed: %v", err)
			}
		}
	}

	response.SuccessResponse(ctx, "操作成功", map[string]interface{}{
		"is_matched": isMatched,
	})
}

// GetReceivedLikes 获取谁喜欢了我 (ActionType = 2)
func (c *Controller) GetReceivedLikes(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	list, err := c.Client.GetClientInteractions(ctx, clientID.(uint64), 2, 0, 50)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取列表失败")
		return
	}

	var result []map[string]interface{}
	for _, item := range list {
		fromClient, _ := c.Client.Get(ctx, item.FromClientID)
		if fromClient != nil {
			result = append(result, map[string]interface{}{
				"interaction_id": item.ID,
				"user_id":        fromClient.ID,
				"name":           string([]rune(fromClient.Name)[0]) + "***", // 脱敏
				"avatar":         fromClient.Avatar,                          // 可以在前端做高斯模糊
				"age":            fromClient.RealAge(),
				"work_city":      fromClient.WorkCity,
				"created_at":     item.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	response.SuccessResponse(ctx, "ok", result)
}

// GetMutualMatches 获取互相心动的列表 (ActionType = 3)
func (c *Controller) GetMutualMatches(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}
	myID := clientID.(uint64)

	list, err := c.Client.GetClientInteractions(ctx, myID, 3, 0, 50)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取列表失败")
		return
	}

	// 互相心动需要去重（因为双方的记录都会被查出来）
	seen := make(map[uint64]bool)
	var result []map[string]interface{}

	for _, item := range list {
		// 找出对方是谁
		var targetID uint64
		if item.FromClientID == myID {
			targetID = item.ToClientID
		} else {
			targetID = item.FromClientID
		}

		if seen[targetID] {
			continue
		}
		seen[targetID] = true

		targetClient, _ := c.Client.Get(ctx, targetID)
		if targetClient != nil {
			result = append(result, map[string]interface{}{
				"interaction_id": item.ID,
				"user_id":        targetClient.ID,
				"name":           targetClient.Name, // 互相心动了，可以展示全名
				"avatar":         targetClient.Avatar,
				"age":            targetClient.RealAge(),
				"phone":          targetClient.Phone, // 也可以直接给联系方式
				"created_at":     item.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	response.SuccessResponse(ctx, "ok", result)
}
