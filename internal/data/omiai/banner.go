package omiai

import (
	"context"
	"fmt"
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"gorm.io/gorm"
)

var _ biz_omiai.BannerInterface = (*BannerRepo)(nil)

type BannerRepo struct {
	db *data.DB
	m  *biz_omiai.Banner
}

func NewBannerRepo(db *data.DB) biz_omiai.BannerInterface {
	return &BannerRepo{db: db, m: new(biz_omiai.Banner)}
}

func (b *BannerRepo) Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*biz_omiai.Banner, error) {
	var bannerList []*biz_omiai.Banner
	err := b.db.Model(b.m).WithContext(ctx).Select(fields).Where(clause.Where, clause.Args...).Order(clause.OrderBy).Offset(offset).Limit(limit).Find(&bannerList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("BannerRepo:Select where:%v err:%w", clause, err)
	}
	return bannerList, nil
}

func (b *BannerRepo) Create(ctx context.Context, banner *biz_omiai.Banner) error {
	return b.db.WithContext(ctx).Model(b.m).Create(banner).Error
}

func (b *BannerRepo) Update(ctx context.Context, banner *biz_omiai.Banner) error {
	return b.db.WithContext(ctx).Model(b.m).Where("id = ?", banner.ID).Updates(banner).Error
}

func (b *BannerRepo) Delete(ctx context.Context, id uint64) error {
	return b.db.WithContext(ctx).Model(b.m).Delete(&biz_omiai.Banner{}, id).Error
}

func (b *BannerRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Banner, error) {
	var banner biz_omiai.Banner
	err := b.db.WithContext(ctx).Model(b.m).First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}
