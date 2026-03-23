package omiai

import (
	"context"
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/data/scopes"
	"time"

	"gorm.io/gorm"
)

var _ biz_omiai.ClientInterface = (*ClientRepo)(nil)

type ClientRepo struct {
	db *gorm.DB
	m  *biz_omiai.Client
}

func NewClientRepo(db *data.DB) biz_omiai.ClientInterface {
	return &ClientRepo{db: db.DB, m: new(biz_omiai.Client)}
}

func (c *ClientRepo) Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*biz_omiai.Client, error) {
	var clientList []*biz_omiai.Client
	err := c.db.Model(c.m).WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Select(fields).Where(clause.Where, clause.Args...).Order(clause.OrderBy).Offset(offset).Limit(limit).Find(&clientList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("ClientRepo:Select where:%v err:%w", clause, err)
	}
	return clientList, nil
}

func (c *ClientRepo) Create(ctx context.Context, client *biz_omiai.Client) error {
	// 如果上下文中带有租户ID且客户端未显式设置，则自动注入
	if tenantID, ok := ctx.Value("tenant_id").(uint64); ok && tenantID > 0 && client.TenantID == 0 {
		client.TenantID = tenantID
	}
	err := c.db.WithContext(ctx).Create(client).Error
	if err != nil {
		return fmt.Errorf("ClientRepo:Create client:%v err:%w", client, err)
	}
	return nil
}

func (c *ClientRepo) Update(ctx context.Context, client *biz_omiai.Client) error {
	err := c.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Save(client).Error
	if err != nil {
		return fmt.Errorf("ClientRepo:Update client:%v err:%w", client, err)
	}
	return nil
}

func (c *ClientRepo) Delete(ctx context.Context, id uint64) error {
	err := c.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Delete(c.m, id).Error
	if err != nil {
		return fmt.Errorf("ClientRepo:Delete id:%d err:%w", id, err)
	}
	return nil
}

func (c *ClientRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Preload("Partner").First(&client, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("ClientRepo:Get id:%d err:%w", id, err)
	}
	return &client, nil
}

func (c *ClientRepo) Stats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 客户总数
	var total int64
	if err := c.db.WithContext(ctx).Model(c.m).Scopes(scopes.TenantScope(ctx)).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 今日新增
	var today int64
	now := time.Now()
	todayStartObj := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if err := c.db.WithContext(ctx).Model(c.m).Scopes(scopes.TenantScope(ctx)).Where("created_at >= ?", todayStartObj).Count(&today).Error; err != nil {
		return nil, err
	}
	stats["today"] = today

	// 单身待匹配客户
	var pending int64
	if err := c.db.WithContext(ctx).Model(c.m).Scopes(scopes.TenantScope(ctx)).Where("status = ?", biz_omiai.ClientStatusSingle).Count(&pending).Error; err != nil {
		return nil, err
	}
	stats["pending"] = pending

	// 已匹配客户数
	var matched int64
	if err := c.db.WithContext(ctx).Model(c.m).Scopes(scopes.TenantScope(ctx)).Where("status = ?", biz_omiai.ClientStatusMatched).Count(&matched).Error; err != nil {
		return nil, err
	}
	stats["matched"] = matched

	return stats, nil
}

// GetByWxOpenID 根据微信OpenID查询客户
func (c *ClientRepo) GetByWxOpenID(ctx context.Context, openID string) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Model(c.m).Scopes(scopes.TenantScope(ctx)).Where("wx_openid = ?", openID).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// SaveInteraction 保存或更新互动记录
func (c *ClientRepo) SaveInteraction(ctx context.Context, interaction *biz_omiai.ClientInteraction) error {
	if tenantID, ok := ctx.Value("tenant_id").(uint64); ok && tenantID > 0 && interaction.TenantID == 0 {
		interaction.TenantID = tenantID
	}
	if interaction.ID > 0 {
		return c.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Model(interaction).Updates(interaction).Error
	}
	return c.db.WithContext(ctx).Model(interaction).Create(interaction).Error
}

// GetInteraction 获取互动记录
func (c *ClientRepo) GetInteraction(ctx context.Context, fromID, toID uint64) (*biz_omiai.ClientInteraction, error) {
	var interaction biz_omiai.ClientInteraction
	err := c.db.WithContext(ctx).Scopes(scopes.TenantScope(ctx)).Model(&biz_omiai.ClientInteraction{}).
		Where("from_client_id = ? AND to_client_id = ?", fromID, toID).
		First(&interaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &interaction, nil
}

// GetInteractionLeads 获取B端红娘的高意向线索 (ActionType=3 且未处理)
func (c *ClientRepo) GetInteractionLeads(ctx context.Context, managerID uint64, offset, limit int) ([]*biz_omiai.ClientInteraction, error) {
	var list []*biz_omiai.ClientInteraction

	query := c.db.WithContext(ctx).Model(&biz_omiai.ClientInteraction{}).Scopes(scopes.TenantScope(ctx)).
		Joins("JOIN client as c1 ON c1.id = client_interaction.from_client_id").
		Joins("JOIN client as c2 ON c2.id = client_interaction.to_client_id").
		Where("client_interaction.action_type = 3 AND client_interaction.status = 0").
		Where("(c1.manager_id = ? OR c2.manager_id = ?)", managerID, managerID)

	err := query.Order("client_interaction.created_at DESC").
		Offset(offset).Limit(limit).Find(&list).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

func (c *ClientRepo) GetClientInteractions(ctx context.Context, clientID uint64, actionType int8, offset, limit int) ([]*biz_omiai.ClientInteraction, error) {
	var list []*biz_omiai.ClientInteraction

	query := c.db.WithContext(ctx).Model(&biz_omiai.ClientInteraction{}).Scopes(scopes.TenantScope(ctx))

	if actionType == 2 {
		// 谁喜欢了我 (to_client_id = me, action_type = 2)
		query = query.Where("to_client_id = ? AND action_type = ?", clientID, actionType)
	} else if actionType == 3 {
		// 互相心动
		query = query.Where("(from_client_id = ? OR to_client_id = ?) AND action_type = ?", clientID, clientID, actionType)
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

// AddCoins 增加/扣除虚拟币，并记录流水
func (c *ClientRepo) AddCoins(ctx context.Context, clientID uint64, amount int, recordType int8, remark string) error {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. 更新余额
	if err := tx.Model(c.m).Where("id = ?", clientID).UpdateColumn("coins", gorm.Expr("coins + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 插入流水
	record := &biz_omiai.ClientCoinRecord{
		ClientID: clientID,
		Amount:   amount,
		Type:     recordType,
		Remark:   remark,
	}
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *ClientRepo) IsVip(ctx context.Context, clientID uint64) bool {
	client, err := c.Get(ctx, clientID)
	if err != nil || client == nil {
		return false
	}
	return client.VipExpireAt.After(time.Now())
}

func (c *ClientRepo) GetDashboardStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	now := time.Now()

	// 客户总数
	var total int64
	if err := c.db.WithContext(ctx).Model(c.m).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["client_total"] = total

	// 今日新增
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var today int64
	if err := c.db.WithContext(ctx).Model(c.m).Where("created_at >= ?", todayStart).Count(&today).Error; err != nil {
		return nil, err
	}
	stats["client_today"] = today

	// 本月新增
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var month int64
	if err := c.db.WithContext(ctx).Model(c.m).Where("created_at >= ?", monthStart).Count(&month).Error; err != nil {
		return nil, err
	}
	stats["client_month"] = month

	return stats, nil
}
