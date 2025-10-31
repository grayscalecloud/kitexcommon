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
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var TracerProvider *tracesdk.TracerProvider

// SetupTracerProviderWithTenantID 在 Kitex provider 创建之前设置包含 TenantIDProcessor 的 TracerProvider
func SetupTracerProviderWithTenantID(serviceName, otelEndpoint string) {
	// 处理 endpoint 格式：提取 host:port
	endpoint := otelEndpoint
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}
	if idx := strings.Index(endpoint, "/"); idx != -1 {
		endpoint = endpoint[:idx]
	}

	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		klog.Errorf("创建 OTLP exporter 失败: %v", err)
		return
	}

	server.RegisterShutdownHook(func() {
		if err := exporter.Shutdown(context.Background()); err != nil {
			klog.Errorf("关闭 exporter 失败: %v", err)
		}
	})

	batchProcessor := tracesdk.NewBatchSpanProcessor(exporter)
	tenantProcessor := NewTenantIDProcessor(batchProcessor)

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		klog.Warnf("创建 resource 失败，使用默认 resource: %v", err)
		res = resource.Default()
	}

	TracerProvider = tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(tenantProcessor),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(TracerProvider)
}

// AddTenantIDProcessorToGlobalTracerProvider 在 Kitex 的 provider 创建之后，用我们包含 TenantIDProcessor 的 TracerProvider 替换全局 TracerProvider
// 必须在 Kitex 的 provider.NewOpenTelemetryProvider 创建之后调用
func AddTenantIDProcessorToGlobalTracerProvider() {
	tp := otel.GetTracerProvider()
	if tp == nil {
		klog.Errorf("全局 TracerProvider 为空，Kitex 的 provider 可能还没有创建")
		return
	}

	if _, ok := tp.(*tracesdk.TracerProvider); !ok {
		klog.Errorf("全局 TracerProvider 不是 *tracesdk.TracerProvider 类型，实际类型: %T", tp)
		return
	}

	if TracerProvider != nil {
		otel.SetTracerProvider(TracerProvider)
		return
	}

	klog.Warnf("没有预先设置的 TracerProvider，无法替换")
}
