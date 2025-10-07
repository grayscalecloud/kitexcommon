package examples

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grayscalecloud/kitexcommon/kvconfig"
)

func RunTestNacos() {
	// 1. 设置你提供的 Nacos 配置
	os.Setenv("NACOS_SERVER_ADDR", "115.190.176.125:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	fmt.Println("=== Nacos 配置中心测试 ===")
	fmt.Printf("服务器地址: %s\n", os.Getenv("NACOS_SERVER_ADDR"))
	fmt.Printf("命名空间: %s\n", os.Getenv("NACOS_NAMESPACE_ID"))
	fmt.Printf("分组: %s\n", os.Getenv("NACOS_GROUP"))
	fmt.Println()

	// 2. 测试创建 Nacos 客户端
	fmt.Println("1. 测试创建 Nacos 客户端...")
	client, err := kvconfig.NewNacosConfigClient(
		[]string{"115.190.176.125:8848"},
		"6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a",
		"DEFAULT_GROUP",
		"nacos",
		"nacos",
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ Nacos 客户端创建成功")

	// 3. 测试获取配置
	fmt.Println("\n2. 测试获取配置...")

	// 测试获取 common 配置
	fmt.Println("尝试获取 'common' 配置...")
	commonContent, err := client.GetConfig("common", "DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 获取 'common' 配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取 'common' 配置成功:\n%s\n", commonContent)
	}

	// 测试获取 pasetopub 配置
	fmt.Println("尝试获取 'pasetopub' 配置...")
	pasetoPubContent, err := client.GetConfig("pasetopub", "DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 获取 'pasetopub' 配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取 'pasetopub' 配置成功:\n%s\n", pasetoPubContent)
	}

	// 测试获取 pasetosecret 配置
	fmt.Println("尝试获取 'pasetosecret' 配置...")
	pasetoSecretContent, err := client.GetConfig("pasetosecret", "DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 获取 'pasetosecret' 配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取 'pasetosecret' 配置成功:\n%s\n", pasetoSecretContent)
	}

	// 4. 测试配置工厂
	fmt.Println("\n3. 测试配置工厂...")
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err = kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}
	fmt.Println("✅ 配置工厂初始化成功")

	// 测试通过工厂获取配置
	fmt.Println("通过工厂获取通用配置...")
	commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 通过工厂获取配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 通过工厂获取配置成功:\n")
		fmt.Printf("  Kitex 服务: %s\n", commonConfig.Kitex.Service)
		fmt.Printf("  Kitex 地址: %s\n", commonConfig.Kitex.Address)
		fmt.Printf("  MySQL DSN: %s\n", commonConfig.MySQL.DSN)
		fmt.Printf("  Redis 地址: %s\n", commonConfig.Redis.Address)
	}

	// 5. 测试配置监听
	fmt.Println("\n4. 测试配置监听...")
	fmt.Println("开始监听 'common' 配置变化（10秒后自动停止）...")

	// 启动配置监听
	err = client.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		fmt.Printf("🔔 配置发生变化！新内容:\n%s\n", content)
	})
	if err != nil {
		fmt.Printf("❌ 启动配置监听失败: %v\n", err)
	} else {
		fmt.Println("✅ 配置监听已启动")

		// 等待 10 秒
		fmt.Println("等待配置变化...（10秒后自动停止）")
		time.Sleep(10 * time.Second)
		fmt.Println("配置监听测试完成")
	}

	// 6. 测试发布配置
	fmt.Println("\n5. 测试发布配置...")
	testConfigContent := `
kitex:
  service: "test-service"
  address: "0.0.0.0:8080"
  metrics_port: "9090"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "info"
  log_file_name: "test.log"
  log_max_size: 100
  log_max_backups: 3
  log_max_age: 7

mysql:
  dsn: "test:test@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  address: "localhost:6379"
  username: ""
  password: ""
  db: 0

otel:
  endpoint: "http://localhost:4317"
  insecure: true
`

	fmt.Println("发布测试配置 'test-config'...")
	err = client.PublishConfig("test-config", "DEFAULT_GROUP", testConfigContent)
	if err != nil {
		fmt.Printf("❌ 发布配置失败: %v\n", err)
	} else {
		fmt.Println("✅ 配置发布成功")

		// 验证配置是否发布成功
		fmt.Println("验证配置是否发布成功...")
		retrievedContent, err := client.GetConfig("test-config", "DEFAULT_GROUP")
		if err != nil {
			fmt.Printf("❌ 获取刚发布的配置失败: %v\n", err)
		} else {
			fmt.Printf("✅ 配置发布验证成功:\n%s\n", retrievedContent)
		}

		// 清理测试配置
		fmt.Println("清理测试配置...")
		err = client.DeleteConfig("test-config", "DEFAULT_GROUP")
		if err != nil {
			fmt.Printf("❌ 删除测试配置失败: %v\n", err)
		} else {
			fmt.Println("✅ 测试配置已清理")
		}
	}

	// 7. 测试批量获取配置
	fmt.Println("\n6. 测试批量获取配置...")
	configs := []kvconfig.ConfigRequest{
		{Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
		{Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
		{Key: "pasetosecret", DataId: "pasetosecret", Group: "DEFAULT_GROUP"},
	}

	batchResults, err := client.BatchGetConfigs(configs)
	if err != nil {
		fmt.Printf("❌ 批量获取配置失败: %v\n", err)
	} else {
		fmt.Println("✅ 批量获取配置成功:")
		for key, content := range batchResults {
			fmt.Printf("  %s: %s\n", key, content[:min(50, len(content))]+"...")
		}
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("如果所有测试都显示 ✅，说明 Nacos 配置中心连接正常！")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
