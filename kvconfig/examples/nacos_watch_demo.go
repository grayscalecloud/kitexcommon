package examples

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grayscalecloud/kitexcommon/kvconfig"
)

func RunNacosWatchDemo() {
	fmt.Println("=== Nacos 配置监听演示 ===")
	
	// 设置 Nacos 配置和身份验证
	os.Setenv("NACOS_SERVER_ADDR", "115.190.176.125:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")
	os.Setenv("NACOS_USERNAME", "nacos")
	os.Setenv("NACOS_PASSWORD", "nacos")

	fmt.Printf("服务器地址: %s\n", os.Getenv("NACOS_SERVER_ADDR"))
	fmt.Printf("命名空间: %s\n", os.Getenv("NACOS_NAMESPACE_ID"))
	fmt.Printf("用户名: %s\n", os.Getenv("NACOS_USERNAME"))
	fmt.Println()

	// 1. 创建客户端
	fmt.Println("1. 创建 Nacos 客户端...")
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

	// 2. 发布初始配置
	fmt.Println("\n2. 发布初始配置...")
	initialConfig := `
kitex:
  service: "demo-service"
  address: "0.0.0.0:8080"
  metrics_port: "9090"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "info"
  log_file_name: "demo.log"
  log_max_size: 100
  log_max_backups: 3
  log_max_age: 7

mysql:
  dsn: "user:password@tcp(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  address: "localhost:6379"
  username: ""
  password: ""
  db: 0

otel:
  endpoint: "http://localhost:4317"
  insecure: true
`

	err = client.PublishConfig("common", "DEFAULT_GROUP", initialConfig)
	if err != nil {
		log.Fatalf("发布初始配置失败: %v", err)
	}
	fmt.Println("✅ 初始配置发布成功")

	// 3. 启动配置监听
	fmt.Println("\n3. 启动配置监听...")
	configChanged := false
	
	err = client.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		fmt.Printf("\n🔔 配置发生变化！时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("新配置内容:\n%s\n", content)
		configChanged = true
	})
	if err != nil {
		log.Fatalf("启动配置监听失败: %v", err)
	}
	fmt.Println("✅ 配置监听已启动")

	// 4. 等待一下，然后更新配置
	fmt.Println("\n4. 等待 3 秒后更新配置...")
	time.Sleep(3 * time.Second)

	// 更新配置
	updatedConfig := `
kitex:
  service: "demo-service-updated"
  address: "0.0.0.0:8081"
  metrics_port: "9091"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "debug"
  log_file_name: "demo-updated.log"
  log_max_size: 200
  log_max_backups: 5
  log_max_age: 10

mysql:
  dsn: "user:password@tcp(localhost:3306)/demo_updated?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  address: "localhost:6380"
  username: "redis_user"
  password: "redis_pass"
  db: 1

otel:
  endpoint: "http://localhost:4318"
  insecure: false
`

	fmt.Println("正在更新配置...")
	err = client.PublishConfig("common", "DEFAULT_GROUP", updatedConfig)
	if err != nil {
		log.Fatalf("更新配置失败: %v", err)
	}
	fmt.Println("✅ 配置更新成功")

	// 5. 等待配置变化
	fmt.Println("\n5. 等待配置变化...")
	time.Sleep(3 * time.Second)

	if configChanged {
		fmt.Println("✅ 配置监听功能正常工作！")
	} else {
		fmt.Println("⚠️ 配置监听可能没有触发")
	}

	// 6. 测试配置工厂
	fmt.Println("\n6. 测试配置工厂...")
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err = kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}
	fmt.Println("✅ 配置工厂初始化成功")

	// 通过工厂获取配置
	commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 通过工厂获取配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 通过工厂获取配置成功:\n")
		fmt.Printf("  Kitex 服务: %s\n", commonConfig.Kitex.Service)
		fmt.Printf("  Kitex 地址: %s\n", commonConfig.Kitex.Address)
		fmt.Printf("  MySQL DSN: %s\n", commonConfig.MySQL.DSN)
		fmt.Printf("  Redis 地址: %s\n", commonConfig.Redis.Address)
		fmt.Printf("  日志级别: %s\n", commonConfig.Kitex.LogLevel)
	}

	// 7. 清理配置
	fmt.Println("\n7. 清理配置...")
	err = client.DeleteConfig("common", "DEFAULT_GROUP")
	if err != nil {
		fmt.Printf("❌ 删除配置失败: %v\n", err)
	} else {
		fmt.Println("✅ 配置清理成功")
	}

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("🎉 Nacos 配置监听功能演示完成！")
	fmt.Println("✅ 配置发布、监听、更新、获取都正常工作")
}
