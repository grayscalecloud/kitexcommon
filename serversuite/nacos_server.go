package serversuite

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/monitor"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosServerSuite struct {
	CurrentServiceName string
	RegistryAddr       string
	NacosPort          uint64
	NamespaceId        string
	OtelEndpoint       string
	EnableMetrics      bool
	EnableTracing      bool
	Username           string
	Password           string
}

func (s NacosServerSuite) Options() []server.Option {
	var opts []server.Option

	var serverAddr string
	var serverPort uint64

	if s.RegistryAddr != "" {
		addr := strings.Split(s.RegistryAddr, ":")
		if len(addr) >= 1 {
			serverAddr = addr[0]
		}
		if len(addr) >= 2 && addr[1] != "" {
			// 修复类型转换问题，使用strconv.ParseUint进行转换
			port, err := strconv.ParseUint(addr[1], 10, 64)
			if err != nil {
				serverPort = 8848 // 默认端口
			} else {
				serverPort = port
			}
		} else if s.NacosPort != 0 {
			serverPort = s.NacosPort
		} else {
			serverPort = 8848 // 默认端口
		}
	} else {
		serverAddr = "127.0.0.1"
		serverPort = 8848
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(serverAddr, serverPort),
	}

	cc := constant.ClientConfig{
		NamespaceId:         s.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		Username:            s.Username,
		Password:            s.Password,
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	r := registry.NewNacosRegistry(cli)
	opts = append(opts, server.WithRegistry(r))

	if s.OtelEndpoint != "" {
		// 初始化 OpenTelemetry Provider
		p := provider.NewOpenTelemetryProvider(
			provider.WithServiceName(s.CurrentServiceName), // 添加服务名
			provider.WithExportEndpoint(s.OtelEndpoint),
			provider.WithEnableMetrics(false),
			provider.WithEnableTracing(s.EnableTracing),
			provider.WithInsecure(),
		)

		// 注册关闭钩子
		server.RegisterShutdownHook(func() {
			if err := p.Shutdown(context.Background()); err != nil {
				klog.Errorf("Failed to shutdown OpenTelemetry provider: %v", err)
			}
		})

		klog.Infof("初始化 otel provider: 当前名字称：%s 注册地址：%s 上报地址：%s",
			s.CurrentServiceName, s.RegistryAddr, s.OtelEndpoint)
	}

	opts = append(opts,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
	)

	if s.EnableTracing {
		opts = append(opts,
			server.WithSuite(tracing.NewServerSuite()),
			server.WithTracer(prometheus.NewServerTracer(s.CurrentServiceName, "",
				prometheus.WithDisableServer(true),
				prometheus.WithRegistry(monitor.Reg))))
	}

	return opts
}
