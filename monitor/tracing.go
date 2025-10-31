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
	"sync"

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

	// 使用反射获取 TracerProvider 的 processors 字段
	tpValue := reflect.ValueOf(sdkTp).Elem()
	tpType := tpValue.Type()

	// 打印所有字段用于调试
	klog.Infof("TracerProvider 类型: %s", tpType.Name())
	for i := 0; i < tpType.NumField(); i++ {
		field := tpType.Field(i)
		klog.Infof("  字段: %s, 类型: %s, 可设置: %v", field.Name, field.Type, tpValue.Field(i).CanSet())
	}

	// 尝试不同的字段名
	processorsField := tpValue.FieldByName("processors")
	if !processorsField.IsValid() {
		// 尝试其他可能的字段名
		processorsField = tpValue.FieldByName("mu")
		if processorsField.IsValid() {
			klog.Infof("找到 mu 字段，尝试查找 processors")
			// OpenTelemetry SDK v1.0+ 可能使用不同的结构
			// 尝试通过其他方式添加 processor
		}
		klog.Errorf("无法找到 processors 字段")
		return
	}

	if !processorsField.CanSet() {
		klog.Errorf("processors 字段不可设置，尝试其他方法")
		// 尝试使用 RegisterSpanProcessor 方法（如果存在）
		registerMethod := reflect.ValueOf(sdkTp).MethodByName("RegisterSpanProcessor")
		if registerMethod.IsValid() {
			klog.Infof("找到 RegisterSpanProcessor 方法，尝试调用")
			tenantProcessor := NewTenantIDProcessor(nil)
			registerMethod.Call([]reflect.Value{reflect.ValueOf(tenantProcessor)})
			klog.Infof("成功通过 RegisterSpanProcessor 方法添加 processor")
			return
		}
		return
	}

	// 获取当前的 processors
	processors := processorsField.Interface()
	klog.Infof("当前 processors 类型: %T", processors)

	// 尝试将其转换为 sync.Map（OpenTelemetry SDK v1.0+ 使用 sync.Map）
	if syncMap, ok := processors.(*sync.Map); ok {
		klog.Infof("processors 是 sync.Map 类型，尝试添加 processor")
		// 遍历 sync.Map，找到第一个 processor 并包装它
		var firstProcessor tracesdk.SpanProcessor
		syncMap.Range(func(key, value interface{}) bool {
			if p, ok := value.(tracesdk.SpanProcessor); ok {
				// 检查是否已经是我们的 processor
				if _, isTenantProcessor := p.(*tenantIDProcessor); !isTenantProcessor {
					firstProcessor = p
					klog.Infof("找到第一个 processor: %T", p)
					return false // 停止遍历
				}
			}
			return true
		})

		if firstProcessor != nil {
			// 包装第一个 processor
			tenantProcessor := NewTenantIDProcessor(firstProcessor)
			// 注意：sync.Map 的 key 不是简单的字符串，需要获取正确的 key
			// 这里我们需要删除旧的 processor 并添加新的
			// 但这可能比较复杂，让我们先尝试直接添加
			syncMap.Range(func(key, value interface{}) bool {
				if value == firstProcessor {
					syncMap.Delete(key)
					syncMap.Store(key, tenantProcessor)
					klog.Infof("成功替换 processor")
					return false
				}
				return true
			})
		} else {
			klog.Errorf("未找到可包装的 processor")
		}
		return
	}

	// 如果不是 sync.Map，尝试作为 slice
	processorsSlice, ok := processors.([]tracesdk.SpanProcessor)
	if !ok {
		klog.Errorf("processors 既不是 sync.Map 也不是 slice，类型: %T", processors)
		return
	}

	klog.Infof("找到 %d 个 processors", len(processorsSlice))

	// 检查是否已经添加了我们的 processor
	hasTenantProcessor := false
	for _, p := range processorsSlice {
		if _, ok := p.(*tenantIDProcessor); ok {
			hasTenantProcessor = true
			break
		}
	}

	// 如果没有，添加我们的 processor（包装第一个 processor）
	if !hasTenantProcessor && len(processorsSlice) > 0 {
		// 包装第一个 processor，在其之前执行我们的逻辑
		firstProcessor := processorsSlice[0]
		klog.Infof("包装第一个 processor: %T", firstProcessor)
		tenantProcessor := NewTenantIDProcessor(firstProcessor)
		// 通过反射设置 processors 字段
		newProcessors := append([]tracesdk.SpanProcessor{tenantProcessor}, processorsSlice[1:]...)
		processorsField.Set(reflect.ValueOf(newProcessors))
		klog.Infof("成功添加 TenantIDProcessor 到 TracerProvider")
	} else if hasTenantProcessor {
		klog.Infof("TenantIDProcessor 已经存在")
	} else {
		klog.Errorf("没有找到任何 processor 可以包装")
	}
}
