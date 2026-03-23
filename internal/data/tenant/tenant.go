package tenant

import (
	"context"
	biz_tenant "omiai-server/internal/biz/tenant"
	"omiai-server/internal/data"

	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
	m  *biz_tenant.Tenant
}

func NewTenantRepo(db *data.DB) biz_tenant.TenantInterface {
	return &TenantRepo{db: db.DB, m: new(biz_tenant.Tenant)}
}

func (r *TenantRepo) Create(ctx context.Context, tenant *biz_tenant.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepo) Get(ctx context.Context, id uint64) (*biz_tenant.Tenant, error) {
	var tenant biz_tenant.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) List(ctx context.Context, offset, limit int) ([]*biz_tenant.Tenant, int64, error) {
	var list []*biz_tenant.Tenant
	var total int64
	
	query := r.db.WithContext(ctx).Model(r.m)
	query.Count(&total)
	
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *TenantRepo) Update(ctx context.Context, tenant *biz_tenant.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}
