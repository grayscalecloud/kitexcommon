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

package serversuite

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/monitor"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	registryconsul "github.com/kitex-contrib/registry-consul"
)

type CommonServerSuite struct {
	CurrentServiceName string
	RegistryAddr       string
	OtelEndpoint       string
}

func (s CommonServerSuite) Options() []server.Option {
	opts := []server.Option{
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
	}

	r, err := registryconsul.NewConsulRegister(s.RegistryAddr)
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// ... consul 配置代码 ...

	// 初始化 OpenTelemetry Provider
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(s.CurrentServiceName), // 添加服务名
		provider.WithExportEndpoint(s.OtelEndpoint),
		provider.WithEnableMetrics(false),
		provider.WithEnableTracing(true),
		provider.WithInsecure(),
	)

	// 添加 TenantIDProcessor 到全局 TracerProvider
	monitor.AddTenantIDProcessorToGlobalTracerProvider()

	// 注册关闭钩子
	server.RegisterShutdownHook(func() {
		if err := p.Shutdown(context.Background()); err != nil {
			klog.Errorf("Failed to shutdown OpenTelemetry provider: %v", err)
		}
	})

	klog.Infof("初始化 otel provider: 当前名字称：%s 注册地址：%s 上报地址：%s",
		s.CurrentServiceName, s.RegistryAddr, s.OtelEndpoint)

	opts = append(opts,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithTracer(prometheus.NewServerTracer(s.CurrentServiceName, "",
			prometheus.WithDisableServer(true),
			prometheus.WithRegistry(monitor.Reg))),
	)

	return opts
}
