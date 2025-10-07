package examples

import (
	"log"
	"os"

	"github.com/grayscalecloud/kitexcommon/kvconfig"
)

func RunQuickStart() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 初始化配置工厂
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err := kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}

	// 3. 获取配置
	commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	log.Printf("配置获取成功:")
	log.Printf("  Kitex 服务: %s", commonConfig.Kitex.Service)
	log.Printf("  Kitex 地址: %s", commonConfig.Kitex.Address)
	log.Printf("  MySQL DSN: %s", commonConfig.MySQL.DSN)
	log.Printf("  Redis 地址: %s", commonConfig.Redis.Address)

	// 4. 开始监听配置变化
	startConfigWatch()

	// 5. 保持程序运行
	select {}
}

func startConfigWatch() {
	// 获取 Nacos 客户端
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 监听配置变化
	err := nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		log.Printf("配置已更新，新内容:\n%s", content)
		
		// 重新获取解析后的配置
		commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
		if err != nil {
			log.Printf("重新获取配置失败: %v", err)
			return
		}
		
		log.Printf("配置重载完成:")
		log.Printf("  Kitex 服务: %s", commonConfig.Kitex.Service)
		log.Printf("  Kitex 地址: %s", commonConfig.Kitex.Address)
	})
	
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	log.Println("配置监听已启动，按 Ctrl+C 退出...")
}
