package scopes

import (
	"context"
	"gorm.io/gorm"
)

// TenantScope GORM 租户隔离拦截器
func TenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 从上下文中获取租户ID
		tenantID, ok := ctx.Value("tenant_id").(uint64)
		if !ok || tenantID == 0 {
			// 如果没有租户信息，或者为 0，且不是特殊超管上下文，则默认不限制 (视业务情况可改为报错)
			// 为了平滑过渡，如果未传递 tenant_id 暂不拦截
			return db
		}

		// 检查上下文是否标记了忽略租户隔离 (如：跨租户查询公海数据时)
		if ignore, ok := ctx.Value("ignore_tenant_scope").(bool); ok && ignore {
			return db
		}

		// 自动追加租户隔离条件
		return db.Where("tenant_id = ?", tenantID)
	}
}

// WithTenant 快速为 context 注入租户ID
func WithTenant(ctx context.Context, tenantID uint64) context.Context {
	return context.WithValue(ctx, "tenant_id", tenantID)
}

// IgnoreTenant 快速为 context 注入忽略租户隔离的标记
func IgnoreTenant(ctx context.Context) context.Context {
	return context.WithValue(ctx, "ignore_tenant_scope", true)
}
