package kvconfig

import (
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/grayscalecloud/kitexcommon/model"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

type CommonConfig struct {
	Env   string
	Kitex model.Kitex `yaml:"kitex"`
	MySQL model.MySQL `yaml:"mysql"`
	Redis model.Redis `yaml:"redis"`
	OTel  model.OTel  `yaml:"otel"`
}

func GetCommonConfig(registryAddr string) (*CommonConfig, error) {
	client, err := api.NewClient(&api.Config{Address: registryAddr})
	if err != nil {
		fmt.Println("Error creating Consul client:", err)
		return nil, err
	}
	//获取配置
	content, _, err := client.KV().Get("onebids/common", nil)
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	conf := new(CommonConfig)
	err = yaml.Unmarshal(content.Value, &conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	return conf, nil
}

func GetKvConfig[T any](registryAddr string, keyName string) (*T, error) {
	client, err := api.NewClient(&api.Config{Address: registryAddr})
	if err != nil {
		fmt.Println("Error creating Consul client:", err)
		return nil, err
	}
	//获取配置
	content, _, err := client.KV().Get(keyName, nil)
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	conf := new(T)
	err = yaml.Unmarshal(content.Value, &conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	return conf, nil
}

func GetPasetoPubConfig(registryAddr string) (*model.PasetoConfig, error) {
	client, err := api.NewClient(&api.Config{Address: registryAddr})
	if err != nil {
		fmt.Println("Error creating Consul client:", err)
		return nil, err
	}
	//获取配置
	content, _, err := client.KV().Get("onebids/pasetopub", nil)
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	conf := new(model.PasetoConfig)
	err = yaml.Unmarshal(content.Value, &conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}

	return conf, nil
}
func GetPasetoSecretConfig(registryAddr string) (*model.PasetoConfig, error) {
	client, err := api.NewClient(&api.Config{Address: registryAddr})
	if err != nil {
		fmt.Println("Error creating Consul client:", err)
		return nil, err
	}
	//获取配置
	content, _, err := client.KV().Get("onebids/pasetosecret", nil)
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	conf := new(model.PasetoConfig)
	err = yaml.Unmarshal(content.Value, &conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}

	return conf, nil
}
