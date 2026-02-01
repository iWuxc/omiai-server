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
	return c.db.WithContext(ctx).Model(c.m).Updates(client).Error
}

func (c *ClientRepo) Delete(ctx context.Context, id uint64) error {
	return c.db.WithContext(ctx).Model(c.m).Delete(&biz_omiai.Client{}, id).Error
}

func (c *ClientRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Model(c.m).First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepo) Stats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	var total int64
	if err := c.db.WithContext(ctx).Model(c.m).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	var today int64
	now := time.Now()
	todayStartObj := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if err := c.db.WithContext(ctx).Model(c.m).Where("created_at >= ?", todayStartObj).Count(&today).Error; err != nil {
		return nil, err
	}
	stats["today"] = today

	// Pending matches: customers with no matched records or specific status
	// For now, let's just mock one more stat
	stats["pending"] = 5 // Mock for now as requested
	return stats, nil
}
