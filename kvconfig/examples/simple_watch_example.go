package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/grayscalecloud/kitexcommon/kvconfig"
)

// SimpleConfigWatchExample 简单的配置监听示例
func SimpleConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 初始化全局配置工厂
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err := kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}

	// 3. 获取 Nacos 客户端
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 4. 监听配置变化
	err = nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		fmt.Printf("配置已更新，新内容:\n%s\n", content)
		
		// 重新获取解析后的配置
		commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
		if err != nil {
			klog.Errorf("重新获取配置失败: %v", err)
			return
		}
		
		// 使用新配置
		fmt.Printf("Kitex 服务: %s\n", commonConfig.Kitex.Service)
		fmt.Printf("Kitex 地址: %s\n", commonConfig.Kitex.Address)
		fmt.Printf("MySQL DSN: %s\n", commonConfig.MySQL.DSN)
	})
	
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("配置监听已启动，按 Ctrl+C 退出...")
	
	// 5. 保持程序运行
	select {}
}

// MultipleConfigWatchExample 监听多个配置的简单示例
func SimpleMultipleConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 初始化全局配置工厂
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err := kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}

	// 3. 获取 Nacos 客户端
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 4. 定义要监听的配置
	configs := []kvconfig.ConfigRequest{
		{Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
		{Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
		{Key: "pasetosecret", DataId: "pasetosecret", Group: "DEFAULT_GROUP"},
	}

	// 5. 批量监听配置
	err = nacosClient.WatchConfigs(configs, func(key, content string) {
		fmt.Printf("配置 [%s] 已更新:\n%s\n", key, content)
		
		// 根据配置类型执行不同操作
		switch key {
		case "common":
			fmt.Println("  通用配置已更新，可能需要重启服务")
		case "pasetopub":
			fmt.Println("  Paseto 公钥已更新")
		case "pasetosecret":
			fmt.Println("  Paseto 密钥已更新")
		}
	})
	
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("多配置监听已启动，按 Ctrl+C 退出...")
	
	// 6. 保持程序运行
	select {}
}

// KitexServiceWatchExample 在 Kitex 服务中监听配置的示例
func KitexServiceWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 初始化全局配置工厂
	options := &kvconfig.ConfigFactoryOptions{
		ConfigType: kvconfig.ConfigTypeNacos,
	}
	err := kvconfig.InitGlobalConfigFactory(options)
	if err != nil {
		log.Fatalf("初始化配置工厂失败: %v", err)
	}

	// 3. 获取初始配置
	commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取初始配置失败: %v", err)
	}

	fmt.Printf("初始配置: Kitex 服务=%s, 地址=%s\n", 
		commonConfig.Kitex.Service, commonConfig.Kitex.Address)

	// 4. 获取 Nacos 客户端并开始监听
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 5. 监听配置变化
	err = nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		fmt.Println("配置已更新，正在重新加载...")
		
		// 重新获取配置
		newConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
		if err != nil {
			klog.Errorf("重新获取配置失败: %v", err)
			return
		}
		
		// 检查配置是否真的发生了变化
		if newConfig.Kitex.Service != commonConfig.Kitex.Service {
			fmt.Printf("Kitex 服务名已从 %s 更改为 %s\n", 
				commonConfig.Kitex.Service, newConfig.Kitex.Service)
			commonConfig.Kitex.Service = newConfig.Kitex.Service
		}
		
		if newConfig.Kitex.Address != commonConfig.Kitex.Address {
			fmt.Printf("Kitex 地址已从 %s 更改为 %s\n", 
				commonConfig.Kitex.Address, newConfig.Kitex.Address)
			commonConfig.Kitex.Address = newConfig.Kitex.Address
		}
		
		if newConfig.MySQL.DSN != commonConfig.MySQL.DSN {
			fmt.Printf("MySQL DSN 已更新\n")
			commonConfig.MySQL.DSN = newConfig.MySQL.DSN
		}
		
		if newConfig.Redis.Address != commonConfig.Redis.Address {
			fmt.Printf("Redis 地址已从 %s 更改为 %s\n", 
				commonConfig.Redis.Address, newConfig.Redis.Address)
			commonConfig.Redis.Address = newConfig.Redis.Address
		}
		
		fmt.Println("配置重载完成")
	})
	
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("Kitex 服务配置监听已启动，按 Ctrl+C 退出...")
	
	// 6. 这里可以启动你的 Kitex 服务
	// 配置变化时会自动触发上面的回调函数
	
	// 保持程序运行
	select {}
}
