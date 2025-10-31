package monitor

import (
	"context"

	"github.com/grayscalecloud/kitexcommon/ctxx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// 实现 SpanProcessor
type tenantIDProcessor struct{ next trace.SpanProcessor }

func NewTenantIDProcessor(next trace.SpanProcessor) trace.SpanProcessor {
	return &tenantIDProcessor{next: next}
}

func (p *tenantIDProcessor) OnStart(ctx context.Context, s trace.ReadWriteSpan) {
	// 测试用固定值，验证 processor 是否被执行
	s.SetAttributes(attribute.String("processor.test", "executed"))

	// 只有在 tenantID 不为空时才设置属性，避免设置空字符串
	if tid := ctxx.GetTenantID(ctx); tid != "" {
		s.SetAttributes(attribute.String("tenant.id", tid))
	} else {
		// 测试用：即使 tenantID 为空，也设置一个标记
		s.SetAttributes(attribute.String("tenant.id.status", "not_found_in_context"))
	}

	// 只有在 merchantID 不为空时才设置属性，避免设置空字符串
	if mid := ctxx.GetMerchantID(ctx); mid != "" {
		s.SetAttributes(attribute.String("merchant.id", mid))
	} else {
		// 测试用：即使 merchantID 为空，也设置一个标记
		s.SetAttributes(attribute.String("merchant.id.status", "not_found_in_context"))
	}
}

func (p *tenantIDProcessor) Shutdown(ctx context.Context) error   { return p.next.Shutdown(ctx) }
func (p *tenantIDProcessor) ForceFlush(ctx context.Context) error { return p.next.ForceFlush(ctx) }
func (p *tenantIDProcessor) OnEnd(s trace.ReadOnlySpan)           { p.next.OnEnd(s) }
