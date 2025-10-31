package monitor

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"go.opentelemetry.io/otel/trace"
)

// TenantIDMiddleware Kitex middleware，自动将 tenantID 和 merchantID 添加到 span 的属性中
// 使用方式：
//   opts = append(opts, server.WithMiddleware(monitor.TenantIDMiddleware))
func TenantIDMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) error {
		// 从 context 中获取当前 span
		span := trace.SpanFromContext(ctx)
		if span.IsRecording() {
			// 调用 AddTenantIDToSpan 添加 tenantID 和 merchantID
			AddTenantIDToSpan(ctx)
		}
		
		// 继续处理请求
		return next(ctx, req, resp)
	}
}

