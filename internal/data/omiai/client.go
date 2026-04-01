package omiai

import (
	"context"
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"time"

	"gorm.io/gorm"
)

var _ biz_omiai.ClientInterface = (*ClientRepo)(nil)

type ClientRepo struct {
	db *data.DB
	m  *biz_omiai.Client
}

func NewClientRepo(db *data.DB) biz_omiai.ClientInterface {
	return &ClientRepo{db: db, m: new(biz_omiai.Client)}
}

func (c *ClientRepo) Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*biz_omiai.Client, error) {
	var clientList []*biz_omiai.Client
	err := c.db.Model(c.m).WithContext(ctx).Select(fields).Where(clause.Where, clause.Args...).Order(clause.OrderBy).Offset(offset).Limit(limit).Find(&clientList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("ClientRepo:Select where:%v err:%w", clause, err)
	}
	return clientList, nil
}

func (c *ClientRepo) Create(ctx context.Context, client *biz_omiai.Client) error {
	return c.db.WithContext(ctx).Model(c.m).Create(client).Error
}

func (c *ClientRepo) Update(ctx context.Context, client *biz_omiai.Client) error {
	return c.db.WithContext(ctx).Model(client).Updates(client).Error
}

func (c *ClientRepo) Delete(ctx context.Context, id uint64) error {
	return c.db.WithContext(ctx).Model(c.m).Delete(&biz_omiai.Client{}, id).Error
}

// HasActiveMatch 检查客户是否有未解除的匹配关系
func (c *ClientRepo) HasActiveMatch(ctx context.Context, clientID uint64) (bool, error) {
	var count int64
	err := c.db.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).
		Where("((male_client_id = ? OR female_client_id = ?) AND status != ?)",
			clientID, clientID, biz_omiai.MatchStatusBroken).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteWithTx 使用事务删除客户，并处理关联数据
func (c *ClientRepo) DeleteWithTx(ctx context.Context, id uint64) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除客户的跟进记录（如果有）
		if err := tx.WithContext(ctx).Where("match_record_id IN (SELECT id FROM match_record WHERE male_client_id = ? OR female_client_id = ?)", id, id).
			Delete(&biz_omiai.FollowUpRecord{}).Error; err != nil {
			return err
		}

		// 2. 删除客户的匹配状态历史记录
		if err := tx.WithContext(ctx).Where("match_record_id IN (SELECT id FROM match_record WHERE male_client_id = ? OR female_client_id = ?)", id, id).
			Delete(&biz_omiai.MatchStatusHistory{}).Error; err != nil {
			return err
		}

		// 3. 删除客户的匹配记录
		if err := tx.WithContext(ctx).Where("male_client_id = ? OR female_client_id = ?", id, id).
			Delete(&biz_omiai.MatchRecord{}).Error; err != nil {
			return err
		}

		// 4. 如果客户有 partner，清除 partner 的 partner_id
		if err := tx.WithContext(ctx).Model(&biz_omiai.Client{}).
			Where("partner_id = ?", id).
			Update("partner_id", nil).Error; err != nil {
			return err
		}

		// 5. 最后删除客户
		if err := tx.WithContext(ctx).Model(c.m).Delete(&biz_omiai.Client{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *ClientRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Model(c.m).Preload("Partner").First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepo) GetByPhone(ctx context.Context, phone string) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Model(c.m).Where("phone = ?", phone).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepo) Stats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 客户总数
	var total int64
	if err := c.db.WithContext(ctx).Model(c.m).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 今日新增
	var today int64
	now := time.Now()
	todayStartObj := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if err := c.db.WithContext(ctx).Model(c.m).Where("created_at >= ?", todayStartObj).Count(&today).Error; err != nil {
		return nil, err
	}
	stats["today"] = today

	// 单身待匹配客户
	var pending int64
	if err := c.db.WithContext(ctx).Model(c.m).Where("status = ?", biz_omiai.ClientStatusSingle).Count(&pending).Error; err != nil {
		return nil, err
	}
	stats["pending"] = pending

	// 已匹配客户数
	var matched int64
	if err := c.db.WithContext(ctx).Model(c.m).Where("status = ?", biz_omiai.ClientStatusMatched).Count(&matched).Error; err != nil {
		return nil, err
	}
	stats["matched"] = matched

	return stats, nil
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

	// BI 漏斗转化指标: 单身(未匹配), 匹配中, 已匹配
	var status1, status2, status3 int64
	c.db.WithContext(ctx).Model(c.m).Where("status = ?", 1).Count(&status1)
	c.db.WithContext(ctx).Model(c.m).Where("status = ?", 2).Count(&status2)
	c.db.WithContext(ctx).Model(c.m).Where("status = ?", 3).Count(&status3)
	stats["client_single"] = status1
	stats["client_matching"] = status2
	stats["client_matched"] = status3

	return stats, nil
}
