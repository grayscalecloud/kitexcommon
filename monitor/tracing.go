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

	// 使用反射获取 TracerProvider 的 spanProcessors 字段
	tpValue := reflect.ValueOf(sdkTp).Elem()
	spanProcessorsField := tpValue.FieldByName("spanProcessors")

	if !spanProcessorsField.IsValid() {
		klog.Errorf("无法找到 spanProcessors 字段")
		return
	}

	klog.Infof("找到 spanProcessors 字段，类型: %s", spanProcessorsField.Type())

	// spanProcessors 是 atomic.Pointer[spanProcessorStates]
	// 我们需要获取原子指针的值
	atomicPtr := spanProcessorsField.Interface()

	// 使用反射调用 Load 方法获取指针值
	loadMethod := reflect.ValueOf(atomicPtr).MethodByName("Load")
	if !loadMethod.IsValid() {
		klog.Errorf("atomic.Pointer 没有 Load 方法")
		return
	}

	statesPtr := loadMethod.Call(nil)[0]
	if statesPtr.IsNil() {
		klog.Errorf("spanProcessorStates 指针为空")
		return
	}

	klog.Infof("获取到 spanProcessorStates 指针")

	// 获取 spanProcessorStates 结构的值
	statesValue := statesPtr.Elem()
	statesType := statesValue.Type()

	// 打印 spanProcessorStates 的所有字段
	klog.Infof("spanProcessorStates 类型: %s", statesType.Name())
	for i := 0; i < statesType.NumField(); i++ {
		field := statesType.Field(i)
		fieldValue := statesValue.Field(i)
		klog.Infof("  字段: %s, 类型: %s, 值类型: %T", field.Name, field.Type, fieldValue.Interface())
	}

	// 尝试找到 processors 字段（可能是 slice 或其他结构）
	processorsField := statesValue.FieldByName("processors")
	if !processorsField.IsValid() {
		// 可能字段名不同，尝试查找 slice 类型的字段
		for i := 0; i < statesType.NumField(); i++ {
			field := statesType.Field(i)
			fieldValue := statesValue.Field(i)
			fieldType := field.Type

			// 检查是否是 slice 类型
			if fieldType.Kind() == reflect.Slice {
				klog.Infof("找到 slice 字段: %s, 类型: %s", field.Name, fieldType)
				processorsSlice := fieldValue.Interface()

				if slice, ok := processorsSlice.([]tracesdk.SpanProcessor); ok {
					klog.Infof("找到 %d 个 processors", len(slice))

					if len(slice) > 0 {
						// 包装第一个 processor
						firstProcessor := slice[0]
						klog.Infof("包装第一个 processor: %T", firstProcessor)
						tenantProcessor := NewTenantIDProcessor(firstProcessor)

						// 创建新的 processors slice，包含我们的 processor
						newSlice := append([]tracesdk.SpanProcessor{tenantProcessor}, slice[1:]...)

						// 使用反射设置字段值
						fieldValue.Set(reflect.ValueOf(newSlice))

						// 使用原子操作更新指针（这需要创建新的 spanProcessorStates）
						// 由于这是复杂的原子操作，我们可以尝试直接修改 slice
						// 但更好的方法是创建一个新的 spanProcessorStates 并原子地更新
						klog.Infof("成功修改 processors slice")
					}
				}
			}
		}
		klog.Errorf("无法在 spanProcessorStates 中找到 processors 字段")
		return
	}

	// 如果找到了 processors 字段
	processors := processorsField.Interface()
	klog.Infof("当前 processors 类型: %T", processors)

	if processorsSlice, ok := processors.([]tracesdk.SpanProcessor); ok {
		klog.Infof("找到 %d 个 processors", len(processorsSlice))

		if len(processorsSlice) > 0 {
			// 包装第一个 processor
			firstProcessor := processorsSlice[0]
			klog.Infof("包装第一个 processor: %T", firstProcessor)
			tenantProcessor := NewTenantIDProcessor(firstProcessor)

			// 创建新的 processors slice
			newSlice := append([]tracesdk.SpanProcessor{tenantProcessor}, processorsSlice[1:]...)

			// 设置新的 slice
			processorsField.Set(reflect.ValueOf(newSlice))

			klog.Infof("成功修改 processors，现在需要原子地更新 spanProcessors")

			// 使用原子操作更新 spanProcessors（Store 方法）
			storeMethod := reflect.ValueOf(atomicPtr).MethodByName("Store")
			if storeMethod.IsValid() {
				storeMethod.Call([]reflect.Value{statesPtr})
				klog.Infof("成功通过原子操作更新 spanProcessors")
			} else {
				klog.Warnf("atomic.Pointer 没有 Store 方法，可能需要其他方式更新")
			}
		}
	} else {
		klog.Errorf("processors 不是 []SpanProcessor 类型，实际类型: %T", processors)
	}
}
