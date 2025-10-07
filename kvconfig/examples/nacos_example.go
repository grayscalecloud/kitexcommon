package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grayscalecloud/kitexcommon/kvconfig"
	"gopkg.in/yaml.v2"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

// NacosConfigExample Nacos 配置使用示例
func NacosConfigExample() {
	// 1. 创建 Nacos 配置客户端
	serverAddrs := []string{"127.0.0.1:8848"}
	namespaceId := "your-namespace-id"
	group := "DEFAULT_GROUP"

	client, err := kvconfig.NewNacosConfigClient(serverAddrs, namespaceId, group, "", "")
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	defer client.Close()

	// 2. 获取通用配置
	commonConfig, err := client.GetCommonConfig(group)
	if err != nil {
		log.Printf("获取通用配置失败: %v", err)
	} else {
		fmt.Printf("通用配置: %+v\n", commonConfig)
	}

	// 3. 获取 Paseto 配置
	pasetoPubConfig, err := client.GetPasetoPubConfig(group)
	if err != nil {
		log.Printf("获取 Paseto 公钥配置失败: %v", err)
	} else {
		fmt.Printf("Paseto 公钥配置: %+v\n", pasetoPubConfig)
	}

	// 4. 获取自定义配置
	type CustomConfig struct {
		Database string `yaml:"database"`
		Cache    string `yaml:"cache"`
	}

	customConfigStr, err := client.GetConfig("custom-config", group)
	if err != nil {
		log.Printf("获取自定义配置失败: %v", err)
		return
	}
	
	// 解析自定义配置
	var customConfig CustomConfig
	err = yaml.Unmarshal([]byte(customConfigStr), &customConfig)
	if err != nil {
		log.Printf("获取自定义配置失败: %v", err)
	} else {
		fmt.Printf("自定义配置: %+v\n", customConfig)
	}

	// 5. 监听配置变化
	err = client.ListenConfig("common", group, func(content string) {
		fmt.Printf("配置发生变化，新内容: %s\n", content)
	})
	if err != nil {
		log.Printf("监听配置失败: %v", err)
	}

	// 6. 带重试的配置获取
	content, err := client.GetConfigWithRetry("common", group, 3, 2*time.Second)
	if err != nil {
		log.Printf("带重试的配置获取失败: %v", err)
	} else {
		fmt.Printf("获取到的配置内容: %s\n", content)
	}

	// 7. 批量获取配置
	configs := []kvconfig.ConfigRequest{
		{Key: "common", DataId: "common", Group: group},
		{Key: "pasetopub", DataId: "pasetopub", Group: group},
		{Key: "custom", DataId: "custom-config", Group: group},
	}

	batchResults, err := client.BatchGetConfigs(configs)
	if err != nil {
		log.Printf("批量获取配置失败: %v", err)
	} else {
		fmt.Printf("批量获取结果: %+v\n", batchResults)
	}

	// 8. 带上下文的配置获取
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	content, err = client.GetConfigWithContext(ctx, "common", group)
	if err != nil {
		log.Printf("带上下文的配置获取失败: %v", err)
	} else {
		fmt.Printf("带上下文获取的配置: %s\n", content)
	}
}

// NacosConfigWithCustomConfigExample 使用自定义配置的示例
func NacosConfigWithCustomConfigExample() {
	// 创建自定义配置
	config := &kvconfig.NacosConfig{
		ServerConfigs: []constant.ServerConfig{
			{
				IpAddr: "127.0.0.1",
				Port:   8848,
			},
		},
		ClientConfig: constant.ClientConfig{
			NamespaceId:         "your-namespace-id",
			TimeoutMs:           10000,
			NotLoadCacheAtStart: false,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "debug",
		},
	}

	client, err := kvconfig.NewNacosConfigClientWithConfig(config)
	if err != nil {
		log.Fatalf("创建自定义 Nacos 客户端失败: %v", err)
	}
	defer client.Close()

	// 使用客户端
	commonConfig, err := client.GetCommonConfig("DEFAULT_GROUP")
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}

	fmt.Printf("使用自定义配置获取的结果: %+v\n", commonConfig)
}

// NacosConfigManagementExample 配置管理示例
func NacosConfigManagementExample() {
	client, err := kvconfig.NewNacosConfigClient(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
		"",
		"",
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	defer client.Close()

	// 发布配置
	configContent := `
kitex:
  service: "example-service"
  address: "0.0.0.0:8080"
  metrics_port: "9090"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "info"
  log_file_name: "example.log"
  log_max_size: 100
  log_max_backups: 3
  log_max_age: 7

mysql:
  dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  address: "localhost:6379"
  username: ""
  password: ""
  db: 0

otel:
  endpoint: "http://localhost:4317"
  insecure: true
`

	err = client.PublishConfig("example-config", "DEFAULT_GROUP", configContent)
	if err != nil {
		log.Printf("发布配置失败: %v", err)
	} else {
		fmt.Println("配置发布成功")
	}

	// 获取刚发布的配置
	content, err := client.GetConfig("example-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("获取配置失败: %v", err)
	} else {
		fmt.Printf("获取到的配置: %s\n", content)
	}

	// 删除配置
	err = client.DeleteConfig("example-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("删除配置失败: %v", err)
	} else {
		fmt.Println("配置删除成功")
	}
}

// NacosConfigWatchExample 配置监听示例
func NacosConfigWatchExample() {
	client, err := kvconfig.NewNacosConfigClient(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
		"",
		"",
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	defer client.Close()

	// 监听多个配置
	configs := []kvconfig.ConfigRequest{
		{Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
		{Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
	}

	err = client.WatchConfigs(configs, func(key, content string) {
		fmt.Printf("配置 [%s] 发生变化:\n%s\n", key, content)
	})
	if err != nil {
		log.Printf("监听配置失败: %v", err)
		return
	}

	fmt.Println("开始监听配置变化，按 Ctrl+C 退出...")

	// 保持程序运行以监听配置变化
	select {}
}
