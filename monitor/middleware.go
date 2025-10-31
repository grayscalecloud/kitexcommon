package monitor

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/grayscalecloud/kitexcommon/ctxx"
	"go.opentelemetry.io/otel/trace"
)

// TenantIDMiddleware Kitex middleware，自动将 tenantID 和 merchantID 添加到 span 的属性中
// 使用方式：
//
//	opts = append(opts, server.WithMiddleware(monitor.TenantIDMiddleware))
func TenantIDMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) error {
		// 添加调试日志
		klog.CtxInfof(ctx, "TenantIDMiddleware 开始执行")
		
		// 从 context 中获取当前 span
		span := trace.SpanFromContext(ctx)
		klog.CtxInfof(ctx, "TenantIDMiddleware: span.IsRecording() = %v", span.IsRecording())
		
		// 检查 context 中是否有 tenantID
		tid := ctxx.GetTenantID(ctx)
		mid := ctxx.GetMerchantID(ctx)
		klog.CtxInfof(ctx, "TenantIDMiddleware: tenantID = %s, merchantID = %s", tid, mid)
		
		if span.IsRecording() {
			// 调用 AddTenantIDToSpan 添加 tenantID 和 merchantID
			AddTenantIDToSpan(ctx)
			klog.CtxInfof(ctx, "TenantIDMiddleware: 已调用 AddTenantIDToSpan")
		} else {
			klog.CtxWarnf(ctx, "TenantIDMiddleware: span 未在记录中，无法添加属性")
		}

		// 继续处理请求
		return next(ctx, req, resp)
	}
}
