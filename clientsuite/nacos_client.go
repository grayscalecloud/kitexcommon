package clientsuite

import (
	"strconv"
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/grayscalecloud/kitexcommon/utils"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/v2/resolver"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosClientSuite struct {
	CurrentServiceName string
	NacosAddr          string
	NacosPort          uint64
	NamespaceId        string
	Username           string
	Password           string
}

func (s NacosClientSuite) Options() []client.Option {
	// 如果以 ： 开头，则默认为本机地址这里强制指定一下，不然服务发现可能出现不可用的IP
	if strings.HasPrefix(s.NacosAddr, ":") {
		s.NacosAddr = utils.MustGetLocalIPv4() + s.NacosAddr
	}

	var serverAddr string
	var serverPort uint64

	if s.NacosAddr != "" {
		addr := strings.Split(s.NacosAddr, ":")
		if len(addr) >= 1 {
			serverAddr = addr[0]
		}
		if len(addr) >= 2 {
			// 修复类型转换问题，使用strconv.ParseUint进行转换
			port, err := strconv.ParseUint(addr[1], 10, 64)
			if err != nil {
				serverPort = 8848 // 默认端口
			} else {
				serverPort = port
			}
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

	r := resolver.NewNacosResolver(cli)
	opts := []client.Option{
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		client.WithSuite(tracing.NewClientSuite()),
	}

	return opts
}
