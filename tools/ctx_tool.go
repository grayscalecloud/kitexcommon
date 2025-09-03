package tools

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/grayscalecloud/kitexcommon/consts"
)

// GetCtxValue 从上下文中获取指定键的值
//
// 如果键不存在，则返回默认值
//
// 参数:
//   - ctx: 上下文
//   - key: 键名
//   - defaultValue: 默认值，当键不存在时返回
//
// 返回:
//   - 上下文中的值或默认值
func GetCtxValue(ctx context.Context, key string, defaultValue string) string {
	lang, ok1 := metainfo.GetValue(ctx, key)
	if ok1 {
		return lang
	} else {
		lang := ctx.Value(key)
		if lang != nil {
			return lang.(string)
		}
		return defaultValue
	}
}

// SetCtxValue 在上下文中设置键值对
//
// 参数:
//   - ctx: 上下文
//   - key: 键名
//   - value: 值
//
// 返回:
//   - 新的上下文
func SetCtxValue(ctx context.Context, key string, value string) context.Context {
	ctx = metainfo.WithValue(ctx, key, value)
	return ctx
}

// GetAccountId 从上下文中获取账户ID
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - 账户ID，如果不存在则返回空字符串
func GetAccountId(ctx context.Context) string {
	return GetCtxValue(ctx, consts.AccountID, "")
}

// GetLanguage 从上下文中获取语言设置
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - 语言代码，如果不存在则返回默认语言
func GetLanguage(ctx context.Context) string {
	return GetCtxValue(ctx, consts.Language, consts.DefaultLanguage)
}

// GetUserID 从上下文中获取用户ID
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - 用户ID，如果不存在则返回空字符串
func GetUserID(ctx context.Context) string {
	return GetCtxValue(ctx, consts.UserID, "")
}

// WithTenant 向上下文中添加租户信息
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - 新的上下文
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return SetCtxValue(ctx, consts.TenantID, tenantID)
}

// GetTenant 从上下文中获取租户ID
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - 租户ID，如果不存在则返回空字符串
func GetTenant(ctx context.Context) string {
	return GetCtxValue(ctx, consts.TenantID, "")
}

// WithTraceID 向上下文中添加追踪ID
//
// 参数:
//   - ctx: 上下文
//   - traceID: 追踪ID
//
// 返回:
//   - 新的上下文
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return SetCtxValue(ctx, consts.TraceID, traceID)
}

// GetTraceID 从上下文中获取追踪ID
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - 追踪ID，如果不存在则返回空字符串
func GetTraceID(ctx context.Context) string {
	return GetCtxValue(ctx, consts.TraceID, "")
}
