package ctxx

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

// SetMetaInfo 设置 metainfo 值，同时设置到 context 和 metainfo 中
func SetMetaInfo(ctx context.Context, key string, value string) context.Context {
	// 设置到 metainfo 中
	ctx = metainfo.WithValue(ctx, key, value)
	// 根据 key 类型设置到对应的 context 中
	switch key {
	case TenantKey:
		ctx = metainfo.WithValue(ctx, TenantKey, value)
		return context.WithValue(ctx, tenantKey{}, value)
	case UserKey:
		ctx = metainfo.WithValue(ctx, UserKey, value)
		return context.WithValue(ctx, userIDKey{}, value)
	case RequestKey:
		ctx = metainfo.WithValue(ctx, RequestKey, value)
		return context.WithValue(ctx, requestIDKey{}, value)
	case MerchantKey:
		ctx = metainfo.WithValue(ctx, MerchantKey, value)
		return context.WithValue(ctx, merchantIDKey{}, value)
	case MemberKey:
		ctx = metainfo.WithValue(ctx, MemberKey, value)
		return context.WithValue(ctx, memberIDKey{}, value)
	case DonorKey:
		ctx = metainfo.WithValue(ctx, DonorKey, value)
		return context.WithValue(ctx, donorIDKey{}, value)
	default:
		// 对于其他 key，只设置到 metainfo 中
		return ctx
	}
}

// GetMetaInfo 获取 metainfo 值，支持 fallback 机制
func GetMetaInfo(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	// 首先尝试从 metainfo 中获取
	if value, ok := metainfo.GetValue(ctx, key); ok && value != "" {
		return value
	}

	// 最后尝试从 context 中获取
	switch key {
	case TenantKey:
		if value, ok := ctx.Value(tenantKey{}).(string); ok {
			if value == "" {
				return ""
			}
			return value
		}
	case UserKey:
		if value, ok := ctx.Value(userIDKey{}).(string); ok {
			return value
		}
	case RequestKey:
		if value, ok := ctx.Value(requestIDKey{}).(string); ok {
			return value
		}
	case MerchantKey:
		if value, ok := ctx.Value(merchantIDKey{}).(string); ok {
			if value == "" {
				return ""
			}
			return value
		}
	case MemberKey:
		if value, ok := ctx.Value(memberIDKey{}).(string); ok {
			if value == "" {
				return ""
			}
			return value
		}
	case DonorKey:
		if value, ok := ctx.Value(donorIDKey{}).(string); ok {
			if value == "" {
				return ""
			}
			return value
		}
	}

	return ""
}

// GetMetaInfoWithFallback 获取 metainfo 值，支持自定义 fallback 键名
func GetMetaInfoWithFallback(ctx context.Context, primaryKey string, fallbackKeys ...string) string {
	if ctx == nil {
		return ""
	}

	// 首先尝试主键
	if value, ok := metainfo.GetValue(ctx, primaryKey); ok && value != "" {
		return value
	}

	// 尝试 fallback 键名
	for _, fallbackKey := range fallbackKeys {
		if value, ok := metainfo.GetValue(ctx, fallbackKey); ok && value != "" {
			return value
		}
	}

	return ""
}

// WithTenant adds tenant ID to the context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return SetMetaInfo(ctx, TenantKey, tenantID)
}

// GetTenantID retrieves tenant ID from the context
func GetTenantID(ctx context.Context) string {
	return GetMetaInfo(ctx, TenantKey)
}

// WithUserID adds user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return SetMetaInfo(ctx, UserKey, userID)
}

// GetUserID retrieves user ID from the context
func GetUserID(ctx context.Context) string {
	return GetMetaInfo(ctx, UserKey)
}

// WithRequestID adds request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return SetMetaInfo(ctx, RequestKey, requestID)
}

// GetRequestID retrieves request ID from the context
func GetRequestID(ctx context.Context) string {
	return GetMetaInfo(ctx, RequestKey)
}

// WithMerchantID adds merchant ID to the context
func WithMerchantID(ctx context.Context, merchantID string) context.Context {
	return SetMetaInfo(ctx, MerchantKey, merchantID)
}

// GetMerchantID retrieves merchant ID from the context
func GetMerchantID(ctx context.Context) string {
	return GetMetaInfo(ctx, MerchantKey)
}

// WithMemberID adds member ID to the context
func WithMemberID(ctx context.Context, memberID string) context.Context {
	return SetMetaInfo(ctx, MemberKey, memberID)
}

// GetMemberID retrieves member ID from the context
func GetMemberID(ctx context.Context) string {
	return GetMetaInfo(ctx, MemberKey)
}

// WithDonorID adds donor ID to the context
func WithDonorID(ctx context.Context, donorID string) context.Context {
	return SetMetaInfo(ctx, DonorKey, donorID)
}

// GetDonorID retrieves donor ID from the context
func GetDonorID(ctx context.Context) string {
	return GetMetaInfo(ctx, DonorKey)
}

// WithTenantIsolation enables or disables tenant isolation for the context
func WithTenantIsolation(ctx context.Context, enabled bool) context.Context {
	enabledStr := "true"
	if !enabled {
		enabledStr = "false"
	}
	ctx = metainfo.WithValue(ctx, TenantIsolationKey, enabledStr)
	return context.WithValue(ctx, tenantIsolationKey{}, enabled)
}

// IsTenantIsolationEnabled checks if tenant isolation is enabled for the context
func IsTenantIsolationEnabled(ctx context.Context) bool {
	if ctx == nil {
		return true // 默认启用租户隔离
	}

	if enabled, ok := ctx.Value(tenantIsolationKey{}).(bool); ok {
		return enabled
	}

	return true // 默认启用租户隔离
}

// WithMerchantIsolation enables or disables merchant isolation for the context
func WithMerchantIsolation(ctx context.Context, enabled bool) context.Context {
	enabledStr := "true"
	if !enabled {
		enabledStr = "false"
	}
	ctx = metainfo.WithValue(ctx, MerchantIsolationKey, enabledStr)
	return context.WithValue(ctx, merchantIsolationKey{}, enabled)
}

// IsMerchantIsolationEnabled checks if merchant isolation is enabled for the context
func IsMerchantIsolationEnabled(ctx context.Context) bool {
	if ctx == nil {
		return true // 默认启用商户隔离
	}

	if enabled, ok := ctx.Value(merchantIsolationKey{}).(bool); ok {
		return enabled
	}

	return true // 默认启用商户隔离
}

// GetContextInfo retrieves all context information
func GetContextInfo(ctx context.Context) *ContextInfo {
	return &ContextInfo{
		TenantID:   GetTenantID(ctx),
		UserID:     GetUserID(ctx),
		RequestID:  GetRequestID(ctx),
		MerchantID: GetMerchantID(ctx),
		MemberID:   GetMemberID(ctx),
		DonorID:    GetDonorID(ctx),
	}
}

// GetAllMetaInfo 获取所有 metainfo 信息
func GetAllMetaInfo(ctx context.Context) map[string]string {
	result := make(map[string]string)

	// 获取所有预定义的键
	keys := []string{TenantKey, UserKey, RequestKey, MerchantKey, MemberKey, DonorKey}
	for _, key := range keys {
		if value := GetMetaInfo(ctx, key); value != "" {
			result[key] = value
		}
	}

	return result
}

// SetMultipleMetaInfo 批量设置 metainfo 值
func SetMultipleMetaInfo(ctx context.Context, values map[string]string) context.Context {
	for key, value := range values {
		ctx = SetMetaInfo(ctx, key, value)
	}
	return ctx
}

// CopyMetaInfo 从源 context 复制 metainfo 到目标 context
func CopyMetaInfo(fromCtx, toCtx context.Context) context.Context {
	keys := []string{TenantKey, UserKey, RequestKey, MerchantKey, MemberKey, DonorKey}
	for _, key := range keys {
		if value := GetMetaInfo(fromCtx, key); value != "" {
			toCtx = SetMetaInfo(toCtx, key, value)
		}
	}
	return toCtx
}

// HasMetaInfo 检查 context 中是否存在指定的 metainfo
func HasMetaInfo(ctx context.Context, key string) bool {
	return GetMetaInfo(ctx, key) != ""
}

// GetMetaInfoOrDefault 获取 metainfo 值，如果不存在则返回默认值
func GetMetaInfoOrDefault(ctx context.Context, key, defaultValue string) string {
	if value := GetMetaInfo(ctx, key); value != "" {
		return value
	}
	return defaultValue
}
