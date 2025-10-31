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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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

	// 获取 spanProcessorStates 的类型（tpType 已经在前面定义）
	spanProcessorsField, found := tpType.FieldByName("spanProcessors")
	if !found {
		klog.Errorf("无法找到 spanProcessors 字段类型")
		return
	}

	// spanProcessors 是 atomic.Pointer[spanProcessorStates]
	// 需要获取内部类型（spanProcessorStates）
	spanProcessorStatesType := spanProcessorsField.Type.Elem().Elem()

	// 将 unsafe.Pointer 转换为 reflect.Value
	statesValue := reflect.NewAt(spanProcessorStatesType, statesPtrUnsafe).Elem()

	statesType := statesValue.Type()
	klog.Infof("spanProcessorStates 类型: %s", statesType.Name())

	// 查找 processors slice 字段
	var processorsOffset uintptr
	for i := 0; i < statesType.NumField(); i++ {
		field := statesType.Field(i)
		if field.Type.Kind() == reflect.Slice {
			klog.Infof("找到 slice 字段: %s, 类型: %s, 偏移量: %d", field.Name, field.Type, field.Offset)
			processorsOffset = field.Offset

			// 检查是否是 []SpanProcessor 类型
			if field.Type.Elem().Name() == "SpanProcessor" || field.Type.String() == "[]trace.SpanProcessor" {
				klog.Infof("这是 processors slice 字段")

				// 使用 unsafe 获取 slice
				statesPtr := unsafe.Pointer(statesValue.UnsafeAddr())
				slicePtr := (*reflect.SliceHeader)(unsafe.Pointer(uintptr(statesPtr) + processorsOffset))

				// 将 SliceHeader 转换为 []SpanProcessor
				slice := *(*[]tracesdk.SpanProcessor)(unsafe.Pointer(slicePtr))
				klog.Infof("找到 %d 个 processors", len(slice))

				if len(slice) > 0 {
					// 包装第一个 processor
					firstProcessor := slice[0]
					klog.Infof("包装第一个 processor: %T", firstProcessor)
					tenantProcessor := NewTenantIDProcessor(firstProcessor)

					// 创建新的 slice，包含我们的 processor
					newSlice := append([]tracesdk.SpanProcessor{tenantProcessor}, slice[1:]...)

					// 使用 unsafe 更新 slice
					newSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&newSlice))
					slicePtr.Data = newSliceHeader.Data
					slicePtr.Len = newSliceHeader.Len
					slicePtr.Cap = newSliceHeader.Cap

					klog.Infof("成功更新 processors slice，新长度: %d", len(newSlice))
				}
				return
			}
		}
	}

	klog.Errorf("无法在 spanProcessorStates 中找到 processors slice 字段")
}
