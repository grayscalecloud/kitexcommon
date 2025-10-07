package examples

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/kvconfig"
	"gopkg.in/yaml.v2"
)

// KitexServerExample Kitex 服务器集成示例
func KitexServerExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 初始化配置工厂
	options := &kvconfig.ConfigFactoryOptions{
		ServerAddr:  "127.0.0.1:8848",
		NamespaceId: "your-namespace-id",
		Group:       "DEFAULT_GROUP",
		Username:    "",
		Password:    "",
		ConfigType:  kvconfig.ConfigTypeNacos,
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

	// 4. 使用配置创建 Kitex 服务器
	addr, err := net.ResolveTCPAddr("tcp", commonConfig.Kitex.Address)
	if err != nil {
		log.Fatalf("解析地址失败: %v", err)
	}
	svr := server.NewServer(
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: commonConfig.Kitex.Service,
		}),
	)

	// 5. 启动服务器
	err = svr.Run()
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// KitexClientExample Kitex 客户端集成示例
func KitexClientExample() {
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

	// 4. 使用配置创建 Kitex 客户端
	// 这里需要根据你的实际服务定义来创建客户端
	klog.Infof("客户端配置: %+v", commonConfig.Kitex)
}

// ConfigWatchExample 配置监听示例
func ConfigWatchExample() {
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

	// 3. 获取 Nacos 客户端进行高级操作
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 4. 监听配置变化
	err = nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
		klog.Infof("配置发生变化: %s", content)
		
		// 重新解析配置
		// 这里可以重新初始化相关的服务配置
		// 比如重新连接数据库、更新日志级别等
	})
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	// 5. 保持程序运行
	select {}
}

// DatabaseConfigExample 数据库配置示例
func DatabaseConfigExample() {
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

	// 3. 获取数据库配置
	commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	// 4. 使用数据库配置
	klog.Infof("数据库 DSN: %s", commonConfig.MySQL.DSN)
	klog.Infof("Redis 地址: %s", commonConfig.Redis.Address)
}

// PasetoConfigExample Paseto 配置示例
func PasetoConfigExample() {
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

	// 3. 获取 Paseto 公钥配置
	pasetoPubConfig, err := kvconfig.GetPasetoPubConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取 Paseto 公钥配置失败: %v", err)
	}

	// 4. 获取 Paseto 密钥配置
	pasetoSecretConfig, err := kvconfig.GetPasetoSecretConfigGlobal("DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取 Paseto 密钥配置失败: %v", err)
	}

	// 5. 使用 Paseto 配置
	klog.Infof("Paseto 公钥: %s", pasetoPubConfig.PubKey)
	klog.Infof("Paseto 密钥: %s", pasetoSecretConfig.PubKey)
}

// CustomConfigExample 自定义配置示例
func CustomConfigExample() {
	// 定义自定义配置结构
	type CustomConfig struct {
		APIKey    string `yaml:"api_key"`
		APISecret string `yaml:"api_secret"`
		BaseURL   string `yaml:"base_url"`
	}

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

	// 3. 获取自定义配置
	customConfigStr, err := kvconfig.GetKvConfigGlobal("custom-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("获取自定义配置失败: %v", err)
		return
	}
	
	// 解析自定义配置
	var customConfig CustomConfig
	err = yaml.Unmarshal([]byte(customConfigStr), &customConfig)
	if err != nil {
		log.Fatalf("获取自定义配置失败: %v", err)
	}

	// 4. 使用自定义配置
	klog.Infof("API Key: %s", customConfig.APIKey)
	klog.Infof("API Secret: %s", customConfig.APISecret)
	klog.Infof("Base URL: %s", customConfig.BaseURL)
}

// ContextExample 上下文使用示例
func ContextExample() {
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

	// 3. 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4. 获取 Nacos 客户端
	nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
	if nacosClient == nil {
		log.Fatalf("获取 Nacos 客户端失败")
	}

	// 5. 使用上下文获取配置
	content, err := nacosClient.GetConfigWithContext(ctx, "common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	klog.Infof("配置内容: %s", content)
}
