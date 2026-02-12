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

func (c *ClientRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Client, error) {
	var client biz_omiai.Client
	err := c.db.WithContext(ctx).Model(c.m).Preload("Partner").First(&client, id).Error
	if err != nil {
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

	return stats, nil
}

func (c *ClientRepo) GetClientTrend(ctx context.Context, days int) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	// 按日期统计每日新增
	type TrendResult struct {
		Date  string `gorm:"column:date"`
		Count int64  `gorm:"column:count"`
	}
	
	var trends []TrendResult
	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	
	err := c.db.WithContext(ctx).Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count 
		FROM client 
		WHERE created_at >= ? 
		GROUP BY DATE(created_at) 
		ORDER BY date ASC
	`, startDate).Scan(&trends).Error
	
	if err != nil {
		return nil, err
	}
	
	// 构建日期和数值数组
	dates := make([]string, 0, days)
	values := make([]int64, 0, days)
	
	// 补全所有日期
	dateMap := make(map[string]int64)
	for _, t := range trends {
		dateMap[t.Date] = t.Count
	}
	
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -days+i+1)
		dateStr := date.Format("2006-01-02")
		dates = append(dates, dateStr)
		values = append(values, dateMap[dateStr])
	}
	
	result["dates"] = dates
	result["values"] = values
	
	return result, nil
}
