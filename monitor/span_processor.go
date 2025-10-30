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
	if tid := ctxx.GetTenantID(ctx); tid != "" {
		s.SetAttributes(attribute.String("tenant.id", tid))
	}
	if mid := ctxx.GetMerchantID(ctx); mid != "" {
		s.SetAttributes(attribute.String("merchant.id", mid))
	}
}

func (p *tenantIDProcessor) Shutdown(ctx context.Context) error   { return p.next.Shutdown(ctx) }
func (p *tenantIDProcessor) ForceFlush(ctx context.Context) error { return p.next.ForceFlush(ctx) }
func (p *tenantIDProcessor) OnEnd(s trace.ReadOnlySpan)           { p.next.OnEnd(s) }
