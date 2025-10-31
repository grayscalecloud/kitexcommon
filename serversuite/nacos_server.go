package serversuite

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/hdmodel"
	"github.com/grayscalecloud/kitexcommon/monitor"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.opentelemetry.io/otel"
)

const (
	// DefaultNacosPort 默认 Nacos 端口
	DefaultNacosPort = 8848
	// DefaultNacosAddr 默认 Nacos 地址
	DefaultNacosAddr = "127.0.0.1"
	// DefaultTimeoutMs 默认超时时间（毫秒）
	DefaultTimeoutMs = 5000
	// DefaultLogDir 默认日志目录
	DefaultLogDir = "/tmp/nacos/log"
	// DefaultCacheDir 默认缓存目录
	DefaultCacheDir = "/tmp/nacos/cache"
	// DefaultLogLevel 默认日志级别
	DefaultLogLevel = "info"
)

// NacosServerSuite Nacos 服务端套件配置
type NacosServerSuite struct {
	// CurrentServiceName 当前服务名称
	CurrentServiceName string
	// RegistryAddr 注册中心地址，格式：host:port
	RegistryAddr string
	// NacosPort Nacos 端口，当 RegistryAddr 未指定端口时使用
	NacosPort uint64
	// NamespaceId Nacos 命名空间 ID
	NamespaceId string
	// Username Nacos 认证用户名
	Username string
	// Password Nacos 认证密码
	Password string
	// Monitor 监控配置
	Monitor *hdmodel.Monitor
}

// parseNacosAddr 解析 Nacos 地址和端口
func (s NacosServerSuite) parseNacosAddr() (string, uint64, error) {
	if s.RegistryAddr == "" {
		return DefaultNacosAddr, DefaultNacosPort, nil
	}

	addr := strings.Split(s.RegistryAddr, ":")
	if len(addr) < 1 || addr[0] == "" {
		return DefaultNacosAddr, DefaultNacosPort, nil
	}

	serverAddr := addr[0]
	var serverPort uint64 = DefaultNacosPort

	// 解析端口
	if len(addr) >= 2 && addr[1] != "" {
		port, err := strconv.ParseUint(addr[1], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("无效的端口号 '%s': %w", addr[1], err)
		}
		serverPort = port
	} else if s.NacosPort != 0 {
		serverPort = s.NacosPort
	}

	return serverAddr, serverPort, nil
}

// validateConfig 验证配置参数
func (s NacosServerSuite) validateConfig() error {
	if s.CurrentServiceName == "" {
		return fmt.Errorf("服务名称不能为空")
	}

	// 验证 Monitor 配置
	if s.Monitor != nil {
		if s.Monitor.OTel.Enable && s.Monitor.OTel.Endpoint == "" {
			return fmt.Errorf("启用 OpenTelemetry 时必须指定端点地址")
		}
	}

	return nil
}

// createNacosClient 创建 Nacos 客户端
func (s NacosServerSuite) createNacosClient() (naming_client.INamingClient, error) {
	serverAddr, serverPort, err := s.parseNacosAddr()
	if err != nil {
		return nil, fmt.Errorf("解析 Nacos 地址失败: %w", err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(serverAddr, serverPort),
	}

	cc := constant.ClientConfig{
		NamespaceId:         s.NamespaceId,
		TimeoutMs:           DefaultTimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              DefaultLogDir,
		CacheDir:            DefaultCacheDir,
		LogLevel:            DefaultLogLevel,
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
		return nil, fmt.Errorf("创建 Nacos 客户端失败: %w", err)
	}

	return cli, nil
}

// setupOpenTelemetry 设置 OpenTelemetry
func (s NacosServerSuite) setupOpenTelemetry() ([]server.Option, error) {
	if s.Monitor == nil || !s.Monitor.OTel.Enable || s.Monitor.OTel.Endpoint == "" {
		return nil, nil
	}

	// **关键：在 Kitex provider 创建之前，先设置包含 TenantIDProcessor 的 TracerProvider**
	// 这样 Kitex 的 provider 可能会使用已存在的全局 TracerProvider，或者我们的 processor 会被包含
	monitor.SetupTracerProviderWithTenantID(s.CurrentServiceName, s.Monitor.OTel.Endpoint)

	// 检查设置后的全局 TracerProvider
	beforeTP := otel.GetTracerProvider()
	klog.Infof("Kitex provider 创建前，全局 TracerProvider: %T, 是否是我们设置的: %v",
		beforeTP, beforeTP == monitor.TracerProvider)

	// 然后创建 Kitex 的 provider
	// 注意：如果 Kitex 的 provider 会创建新的 TracerProvider，我们需要在创建后再次尝试添加 processor
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(s.CurrentServiceName),
		provider.WithExportEndpoint(s.Monitor.OTel.Endpoint),
		provider.WithEnableMetrics(false),
		provider.WithEnableTracing(s.Monitor.OTel.Enable),
		provider.WithInsecure(),
	)

	// 检查 Kitex provider 创建后的全局 TracerProvider
	afterTP := otel.GetTracerProvider()
	klog.Infof("Kitex provider 创建后，全局 TracerProvider: %T, 是否是我们设置的: %v",
		afterTP, afterTP == monitor.TracerProvider)

	if afterTP != monitor.TracerProvider {
		klog.Warnf("Kitex provider 覆盖了我们的 TracerProvider，尝试添加 processor")
		// 如果 Kitex 的 provider 覆盖了我们的 TracerProvider，尝试添加 processor
		monitor.AddTenantIDProcessorToGlobalTracerProvider()
	} else {
		klog.Infof("Kitex provider 使用我们设置的 TracerProvider，TenantIDProcessor 应该已生效")
	}

	// 注册关闭钩子
	server.RegisterShutdownHook(func() {
		if err := p.Shutdown(context.Background()); err != nil {
			klog.Errorf("关闭 OpenTelemetry provider 失败: %v", err)
		}
	})

	klog.Infof("初始化 otel provider: 当前服务名称：%s 注册地址：%s 上报地址：%s",
		s.CurrentServiceName, s.RegistryAddr, s.Monitor.OTel.Endpoint)

	return nil, nil
}

// setupTracing 设置链路追踪
func (s NacosServerSuite) setupTracing() []server.Option {
	if s.Monitor == nil || !s.Monitor.OTel.Enable {
		return nil
	}

	return []server.Option{
		server.WithSuite(tracing.NewServerSuite()),
		server.WithTracer(prometheus.NewServerTracer(s.CurrentServiceName, "",
			prometheus.WithDisableServer(true),
			prometheus.WithRegistry(monitor.Reg))),
	}
}

// Options 返回服务器选项配置
func (s NacosServerSuite) Options() []server.Option {
	// 验证配置
	if err := s.validateConfig(); err != nil {
		klog.Fatalf("配置验证失败: %v", err)
	}

	var opts []server.Option

	// 创建 Nacos 客户端
	cli, err := s.createNacosClient()
	if err != nil {
		klog.Fatalf("创建 Nacos 客户端失败: %v", err)
	}

	// 设置注册中心
	r := registry.NewNacosRegistry(cli)
	opts = append(opts, server.WithRegistry(r))

	// 设置 OpenTelemetry
	if _, err := s.setupOpenTelemetry(); err != nil {
		klog.Fatalf("设置 OpenTelemetry 失败: %v", err)
	}

	// 设置服务基本信息
	opts = append(opts,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
	)

	// 设置链路追踪
	if tracingOpts := s.setupTracing(); tracingOpts != nil {
		opts = append(opts, tracingOpts...)
	}

	return opts
}
