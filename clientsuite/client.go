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

package clientsuite

import (
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/grayscalecloud/kitexcommon/utils"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonClientSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonClientSuite) Options() []client.Option {
	// 如果以 ： 开头，则默认为本机地址这里强制指定一下，不然服务发现可能出现不可用的IP
	if strings.HasPrefix(s.RegistryAddr, ":") {
		s.RegistryAddr = utils.MustGetLocalIPv4() + s.RegistryAddr
	}
	r, err := consul.NewConsulResolver(s.RegistryAddr)
	if err != nil {
		panic(err)
	}
	opts := []client.Option{
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		//client.WithMuxConnection(1),                                // multiplexing 这个在windown 下不支持
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		client.WithSuite(tracing.NewClientSuite()),
	}

	return opts
}
