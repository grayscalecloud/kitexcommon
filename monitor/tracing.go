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
		// 如果没有全局 TracerProvider，说明 Kitex 的 provider 还没有创建
		return
	}

	// 检查是否是 TracerProvider 类型
	sdkTp, ok := tp.(*tracesdk.TracerProvider)
	if !ok {
		// 如果不是 SDK 的 TracerProvider，无法添加 processor
		return
	}

	// 使用反射获取 TracerProvider 的 processors 字段
	// 注意：这是通过反射访问私有字段，可能在不同版本的 SDK 中失效
	tpValue := reflect.ValueOf(sdkTp).Elem()
	processorsField := tpValue.FieldByName("processors")
	if !processorsField.IsValid() || !processorsField.CanSet() {
		// 如果无法访问或设置 processors 字段，尝试使用 syncMap 字段
		// OpenTelemetry SDK 可能使用 sync.Map 来存储 processors
		syncMapField := tpValue.FieldByName("processors")
		if syncMapField.IsValid() {
			// 尝试通过 sync.Map 添加 processor（这更复杂，需要根据实际 SDK 版本调整）
			return
		}
		return
	}

	// 获取当前的 processors
	processors := processorsField.Interface()

	// 尝试将其转换为 sync.Map（OpenTelemetry SDK v1.0+ 使用 sync.Map）
	if syncMap, ok := processors.(*sync.Map); ok {
		// 遍历 sync.Map，找到第一个 processor 并包装它
		var firstProcessor tracesdk.SpanProcessor
		syncMap.Range(func(key, value interface{}) bool {
			if p, ok := value.(tracesdk.SpanProcessor); ok {
				// 检查是否已经是我们的 processor
				if _, isTenantProcessor := p.(*tenantIDProcessor); !isTenantProcessor {
					firstProcessor = p
					return false // 停止遍历
				}
			}
			return true
		})

		if firstProcessor != nil {
			// 包装第一个 processor
			tenantProcessor := NewTenantIDProcessor(firstProcessor)
			syncMap.Store("tenantIDProcessor", tenantProcessor)
		}
		return
	}

	// 如果不是 sync.Map，尝试作为 slice
	processorsSlice, ok := processors.([]tracesdk.SpanProcessor)
	if !ok {
		return
	}

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
		tenantProcessor := NewTenantIDProcessor(firstProcessor)
		// 通过反射设置 processors 字段
		newProcessors := append([]tracesdk.SpanProcessor{tenantProcessor}, processorsSlice[1:]...)
		processorsField.Set(reflect.ValueOf(newProcessors))
	}
}
