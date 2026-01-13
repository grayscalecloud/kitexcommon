package ctxx

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

// Define key types for context name values

// SetMetaInfo 设置 metainfo 值，同时设置到 context 和 metainfo 中
func SetMetaInfo(ctx context.Context, key string, value string) context.Context {
	return metainfo.WithValue(ctx, key, value)
}

// GetMetaInfo 获取 metainfo 值
func GetMetaInfo(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	if value, ok := metainfo.GetValue(ctx, key); ok && value != "" {
		return value
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

// WithAppType adds app type to the context
func WithAppType(ctx context.Context, appType string) context.Context {
	return SetMetaInfo(ctx, AppTypeKey, appType)
}

// GetAppType retrieves app type from the context
func GetAppType(ctx context.Context) string {
	return GetMetaInfo(ctx, AppTypeKey)
}

// WithTenantIsolation enables or disables tenant isolation for the context
func WithTenantIsolation(ctx context.Context, enabled bool) context.Context {
	enabledStr := "true"
	if !enabled {
		enabledStr = "false"
	}
	return metainfo.WithValue(ctx, TenantIsolationKey, enabledStr)
}

// IsTenantIsolationEnabled checks if tenant isolation is enabled for the context
func IsTenantIsolationEnabled(ctx context.Context) bool {
	if ctx == nil {
		return true // 默认启用租户隔离
	}

	value := GetMetaInfo(ctx, TenantIsolationKey)
	if value == "false" {
		return false
	}

	return true // 默认启用租户隔离
}

// WithMerchantIsolation enables or disables merchant isolation for the context
func WithMerchantIsolation(ctx context.Context, enabled bool) context.Context {
	enabledStr := "true"
	if !enabled {
		enabledStr = "false"
	}
	return metainfo.WithValue(ctx, MerchantIsolationKey, enabledStr)
}

// IsMerchantIsolationEnabled checks if merchant isolation is enabled for the context
func IsMerchantIsolationEnabled(ctx context.Context) bool {
	if ctx == nil {
		return true // 默认启用商户隔离
	}

	value := GetMetaInfo(ctx, MerchantIsolationKey)
	if value == "false" {
		return false
	}

	return true // 默认启用商户隔离
}

// GetContextInfo retrieves all context information
func GetContextInfo(ctx context.Context) *ContextInfo {
	return &ContextInfo{
		TenantID:            GetTenantID(ctx),
		UserID:              GetUserID(ctx),
		RequestID:           GetRequestID(ctx),
		MerchantID:          GetMerchantID(ctx),
		MemberID:            GetMemberID(ctx),
		DonorID:             GetDonorID(ctx),
		AppType:             GetAppType(ctx),
		TenantName:          GetTenantName(ctx),
		UserName:            GetUserName(ctx),
		MerchantName:        GetMerchantName(ctx),
		MemberName:          GetMemberName(ctx),
		DonorName:           GetDonorName(ctx),
		AppName:             GetAppName(ctx),
		ExpandedInfo:        GetMetaInfo(ctx, ExpandedKey),
		AppId:               GetMetaInfo(ctx, AppIdKey),
		Ip:                  GetMetaInfo(ctx, IpKey),
		SkipDesensitization: IsSkipDesensitizationEnabled(ctx),
	}
}

// GetAllMetaInfo 获取所有 metainfo 信息
func GetAllMetaInfo(ctx context.Context) map[string]string {
	result := make(map[string]string)

	// 获取所有预定义的键（包括 ID、Name 和 UserAgent）
	keys := []string{
		TenantKey, UserKey, RequestKey, MerchantKey, MemberKey, DonorKey, AppTypeKey,
		TenantNameKey, UserNameKey, MerchantNameKey, MemberNameKey, DonorNameKey, AppNameKey,
		UserAgentKey,
	}
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
	keys := []string{
		TenantKey, UserKey, RequestKey, MerchantKey, MemberKey, DonorKey, AppTypeKey,
		TenantNameKey, UserNameKey, MerchantNameKey, MemberNameKey, DonorNameKey, AppNameKey,
		UserAgentKey,
	}
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

// WithTenantName adds tenant name to the context
func WithTenantName(ctx context.Context, tenantName string) context.Context {
	return SetMetaInfo(ctx, TenantNameKey, tenantName)
}

// GetTenantName retrieves tenant name from the context
func GetTenantName(ctx context.Context) string {
	return GetMetaInfo(ctx, TenantNameKey)
}

// WithMemberName adds member name to the context
func WithMemberName(ctx context.Context, memberName string) context.Context {
	return SetMetaInfo(ctx, MemberNameKey, memberName)
}

// GetMemberName retrieves member name from the context
func GetMemberName(ctx context.Context) string {
	return GetMetaInfo(ctx, MemberNameKey)
}

// WithMerchantName adds merchant name to the context
func WithMerchantName(ctx context.Context, merchantName string) context.Context {
	return SetMetaInfo(ctx, MerchantNameKey, merchantName)
}

// GetMerchantName retrieves merchant name from the context
func GetMerchantName(ctx context.Context) string {
	return GetMetaInfo(ctx, MerchantNameKey)
}

// WithUserName adds user name to the context
func WithUserName(ctx context.Context, userName string) context.Context {
	return SetMetaInfo(ctx, UserNameKey, userName)
}

// GetUserName retrieves user name from the context
func GetUserName(ctx context.Context) string {
	return GetMetaInfo(ctx, UserNameKey)
}

// WithDonorName adds donor name to the context
func WithDonorName(ctx context.Context, donorName string) context.Context {
	return SetMetaInfo(ctx, DonorNameKey, donorName)
}

// GetDonorName retrieves donor name from the context
func GetDonorName(ctx context.Context) string {
	return GetMetaInfo(ctx, DonorNameKey)
}

// WithAppName adds app name to the context
func WithAppName(ctx context.Context, appName string) context.Context {
	return SetMetaInfo(ctx, AppNameKey, appName)
}

// GetAppName retrieves app name from the context
func GetAppName(ctx context.Context) string {
	return GetMetaInfo(ctx, AppNameKey)
}

// WithUserAgent adds user agent to the context
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return SetMetaInfo(ctx, UserAgentKey, userAgent)
}

// GetUserAgent retrieves user agent from the context
func GetUserAgent(ctx context.Context) string {
	return GetMetaInfo(ctx, UserAgentKey)
}

// WithExpandedInfo adds expanded info to the context
func WithExpandedInfo(ctx context.Context, key string, expandedInfo string) context.Context {
	return SetMetaInfo(ctx, key, expandedInfo)
}

// GetExpandedInfo retrieves expanded info from the context
func GetExpandedInfo(ctx context.Context, key string) string {
	return GetMetaInfo(ctx, key)
}

// WithAppId adds app id to the context
func WithAppId(ctx context.Context, appId string) context.Context {
	return SetMetaInfo(ctx, AppIdKey, appId)
}

// GetAppId retrieves app id from the context
func GetAppId(ctx context.Context) string {
	return GetMetaInfo(ctx, AppIdKey)
}

// WithIp adds ip to the context
func WithIp(ctx context.Context, ip string) context.Context {
	return SetMetaInfo(ctx, IpKey, ip)
}

// GetIp retrieves ip from the context
func GetIp(ctx context.Context) string {
	return GetMetaInfo(ctx, IpKey)
}

// WithTenantType adds tenant type to the context
func WithTenantType(ctx context.Context, tenantType string) context.Context {
	return SetMetaInfo(ctx, TenantTypeKey, tenantType)
}

// GetTenantType 获取租户类型
func GetTenantType(ctx context.Context) string {
	return GetMetaInfo(ctx, TenantTypeKey)
}

// WithSkipDesensitization enables or disables skipping desensitization for the context
func WithSkipDesensitization(ctx context.Context, skip bool) context.Context {
	skipStr := "true"
	if !skip {
		skipStr = "false"
	}
	return SetMetaInfo(ctx, SkipDesensitizationKey, skipStr)
}

// IsSkipDesensitizationEnabled checks if desensitization should be skipped for the context
func IsSkipDesensitizationEnabled(ctx context.Context) bool {
	if ctx == nil {
		return false // 默认不跳过脱敏（即启用脱敏）
	}

	value := GetMetaInfo(ctx, SkipDesensitizationKey)
	if value == "true" {
		return true
	}

	return false // 默认不跳过脱敏（即启用脱敏）
}
