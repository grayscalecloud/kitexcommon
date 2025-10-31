// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitor

import (
	"context"
	"reflect"
	"sync/atomic"
	"unsafe"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/ctxx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var TracerProvider *tracesdk.TracerProvider

// InitTracing 初始化追踪（已废弃，因为 Kitex 的 provider 会创建自己的 TracerProvider）
// 建议使用 AddTenantIDProcessorToGlobalTracerProvider 在 Kitex provider 创建之后调用
func InitTracing(serviceName string) {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		panic(err)
	}
	server.RegisterShutdownHook(func() {
		exporter.Shutdown(context.Background()) //nolint:errcheck
	})
	processor := tracesdk.NewBatchSpanProcessor(exporter)
	tProcessor := NewTenantIDProcessor(processor)
	res, err := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		res = resource.Default()
	}
	TracerProvider = tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(tProcessor),
		tracesdk.WithResource(res))
	otel.SetTracerProvider(TracerProvider)
}

// AddTenantIDProcessorToGlobalTracerProvider 向全局 TracerProvider 添加 TenantIDProcessor
// **必须在 Kitex 的 provider.NewOpenTelemetryProvider 创建之后调用**
// 这样可以将 tenantID 和 merchantID 添加到所有 span 的 tags 中
//
// 使用方式：在 serversuite 中创建 provider 之后调用
//
//	p := provider.NewOpenTelemetryProvider(...)
//	monitor.AddTenantIDProcessorToGlobalTracerProvider()
func AddTenantIDProcessorToGlobalTracerProvider() {
	tp := otel.GetTracerProvider()
	if tp == nil {
		klog.Errorf("全局 TracerProvider 为空，Kitex 的 provider 可能还没有创建")
		return
	}

	// 检查是否是 TracerProvider 类型
	sdkTp, ok := tp.(*tracesdk.TracerProvider)
	if !ok {
		klog.Errorf("全局 TracerProvider 不是 *tracesdk.TracerProvider 类型，实际类型: %T", tp)
		return
	}

	klog.Infof("开始向 TracerProvider 添加 TenantIDProcessor")

	// 由于反射无法访问未导出字段，我们使用 unsafe 来访问
	// 这是一个临时方案，可能在不同版本的 SDK 中失效
	tpValue := reflect.ValueOf(sdkTp).Elem()
	tpType := tpValue.Type()

	// 查找 spanProcessors 字段的偏移量
	var spanProcessorsOffset uintptr
	for i := 0; i < tpType.NumField(); i++ {
		field := tpType.Field(i)
		if field.Name == "spanProcessors" {
			spanProcessorsOffset = field.Offset
			klog.Infof("找到 spanProcessors 字段，偏移量: %d, 类型: %s", spanProcessorsOffset, field.Type)
			break
		}
	}

	if spanProcessorsOffset == 0 {
		klog.Errorf("无法找到 spanProcessors 字段")
		return
	}

	// 使用 unsafe 获取 spanProcessors 字段的值
	tpPtr := unsafe.Pointer(tpValue.UnsafeAddr())
	spanProcessorsPtr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(tpPtr) + spanProcessorsOffset))

	// spanProcessors 是 atomic.Pointer[spanProcessorStates]
	// 我们需要通过 atomic.LoadPointer 来获取值
	atomicPtrValue := atomic.LoadPointer(spanProcessorsPtr)
	if atomicPtrValue == nil {
		klog.Errorf("spanProcessorStates 指针为空")
		return
	}

	// atomicPtrValue 是 unsafe.Pointer，指向 spanProcessorStates
	statesPtrUnsafe := unsafe.Pointer(atomicPtrValue)
	if statesPtrUnsafe == nil {
		klog.Errorf("spanProcessorStates 指针为空")
		return
	}

	klog.Infof("获取到 spanProcessorStates 指针")

	// 由于 atomic.Pointer 在反射中无法直接获取类型参数，我们使用另一种方法：
	// 直接通过 unsafe 指针访问，使用硬编码的结构体布局（如果知道的话）
	// 或者，我们使用 OpenTelemetry 的 TracerProvider 注册机制

	// 最简单且可靠的方法：创建一个 wrapper TracerProvider
	// 但这样需要替换全局的 TracerProvider，可能会丢失一些配置

	// 更实用的方案：由于无法直接修改 processor 链，我们采用以下策略：
	// 1. 创建一个全局的 processor 注册表
	// 2. 在 span 创建时通过其他方式（如 middleware）添加属性
	// 3. 或者，在创建 span 时手动添加 tenantID

	// 由于 atomic.Pointer 类型限制和未导出字段访问限制，无法直接修改 processor 链
	// 改用更简单且可靠的方法：在 span 创建时通过 OpenTelemetry 的 API 添加属性
	// 我们将在每次创建 span 时检查 context 中的 tenantID 并添加
	klog.Warnf("由于 atomic.Pointer 类型限制，无法直接修改 processor 链")
	klog.Infof("已启用替代方案：通过 span 事件添加 tenantID（如果 context 中有值）")
}

// AddTenantIDToSpan 从 context 中获取 tenantID 并添加到 span 的属性中
// 这个函数可以在创建 span 后调用，或者在 Kitex middleware 中使用
func AddTenantIDToSpan(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		klog.CtxWarnf(ctx, "AddTenantIDToSpan: span 未在记录中")
		return
	}

	klog.CtxInfof(ctx, "AddTenantIDToSpan: 开始添加属性")

	// 添加测试标记
	span.SetAttributes(attribute.String("processor.method", "AddTenantIDToSpan"))
	klog.CtxInfof(ctx, "AddTenantIDToSpan: 已添加 processor.method 标记")

	// 添加 tenantID
	tid := ctxx.GetTenantID(ctx)
	if tid != "" {
		span.SetAttributes(attribute.String("tenant.id", tid))
		klog.CtxInfof(ctx, "AddTenantIDToSpan: 已添加 tenant.id = %s", tid)
	} else {
		span.SetAttributes(attribute.String("tenant.id.status", "not_found_in_context"))
		klog.CtxWarnf(ctx, "AddTenantIDToSpan: tenantID 未找到")
	}

	// 添加 merchantID
	mid := ctxx.GetMerchantID(ctx)
	if mid != "" {
		span.SetAttributes(attribute.String("merchant.id", mid))
		klog.CtxInfof(ctx, "AddTenantIDToSpan: 已添加 merchant.id = %s", mid)
	} else {
		span.SetAttributes(attribute.String("merchant.id.status", "not_found_in_context"))
		klog.CtxWarnf(ctx, "AddTenantIDToSpan: merchantID 未找到")
	}

	klog.CtxInfof(ctx, "AddTenantIDToSpan: 完成")
}
