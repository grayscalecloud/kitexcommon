package examples

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/grayscalecloud/kitexcommon/kvconfig"
	"gopkg.in/yaml.v2"
)

// ConfigWatcher 配置监听器
type ConfigWatcher struct {
	client        *kvconfig.NacosConfigClient
	configs       map[string]*kvconfig.CommonConfig
	mu            sync.RWMutex
	callbacks     []func(string, *kvconfig.CommonConfig)
	stopChan      chan struct{}
}

// NewConfigWatcher 创建配置监听器
func NewConfigWatcher(serverAddrs []string, namespaceId, group string) (*ConfigWatcher, error) {
	client, err := kvconfig.NewNacosConfigClient(serverAddrs, namespaceId, group, "", "")
	if err != nil {
		return nil, fmt.Errorf("创建 Nacos 客户端失败: %w", err)
	}

	return &ConfigWatcher{
		client:    client,
		configs:   make(map[string]*kvconfig.CommonConfig),
		callbacks: make([]func(string, *kvconfig.CommonConfig), 0),
		stopChan:  make(chan struct{}),
	}, nil
}

// AddCallback 添加配置变化回调
func (w *ConfigWatcher) AddCallback(callback func(string, *kvconfig.CommonConfig)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks = append(w.callbacks, callback)
}

// WatchConfig 监听单个配置
func (w *ConfigWatcher) WatchConfig(dataId, group string) error {
	return w.client.ListenConfig(dataId, group, func(content string) {
		klog.Infof("配置 [%s] 发生变化", dataId)
		
		// 解析配置
		config := new(kvconfig.CommonConfig)
		err := yaml.Unmarshal([]byte(content), config)
		if err != nil {
			klog.Errorf("解析配置失败 [%s]: %v", dataId, err)
			return
		}

		// 更新缓存
		w.mu.Lock()
		w.configs[dataId] = config
		w.mu.Unlock()

		// 触发回调
		w.mu.RLock()
		for _, callback := range w.callbacks {
			go callback(dataId, config)
		}
		w.mu.RUnlock()
	})
}

// WatchMultipleConfigs 监听多个配置
func (w *ConfigWatcher) WatchMultipleConfigs(configs []kvconfig.ConfigRequest) error {
	for _, req := range configs {
		err := w.WatchConfig(req.DataId, req.Group)
		if err != nil {
			return fmt.Errorf("监听配置失败 [%s]: %w", req.DataId, err)
		}
	}
	return nil
}

// GetConfig 获取缓存的配置
func (w *ConfigWatcher) GetConfig(dataId string) (*kvconfig.CommonConfig, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	config, exists := w.configs[dataId]
	return config, exists
}

// Stop 停止监听
func (w *ConfigWatcher) Stop() {
	close(w.stopChan)
	w.client.Close()
}

// 基本配置监听示例
func BasicConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 3. 添加配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		fmt.Printf("配置 [%s] 已更新:\n", dataId)
		fmt.Printf("  Kitex 服务: %s\n", config.Kitex.Service)
		fmt.Printf("  Kitex 地址: %s\n", config.Kitex.Address)
		fmt.Printf("  MySQL DSN: %s\n", config.MySQL.DSN)
		fmt.Printf("  Redis 地址: %s\n", config.Redis.Address)
		fmt.Println("---")
	})

	// 4. 监听配置
	err = watcher.WatchConfig("common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("开始监听配置变化，按 Ctrl+C 退出...")
	
	// 5. 保持程序运行
	select {}
}

// 多配置监听示例
func MultipleConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 3. 添加配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		fmt.Printf("配置 [%s] 已更新，时间: %s\n", dataId, time.Now().Format("2006-01-02 15:04:05"))
		
		// 根据不同配置类型执行不同操作
		switch dataId {
		case "common":
			fmt.Println("  通用配置已更新，可能需要重启服务")
		case "pasetopub":
			fmt.Println("  Paseto 公钥已更新，需要重新加载")
		case "pasetosecret":
			fmt.Println("  Paseto 密钥已更新，需要重新加载")
		default:
			fmt.Printf("  未知配置类型: %s\n", dataId)
		}
	})

	// 4. 监听多个配置
	configs := []kvconfig.ConfigRequest{
		{Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
		{Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
		{Key: "pasetosecret", DataId: "pasetosecret", Group: "DEFAULT_GROUP"},
	}

	err = watcher.WatchMultipleConfigs(configs)
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("开始监听多个配置变化，按 Ctrl+C 退出...")
	
	// 5. 保持程序运行
	select {}
}

// 配置监听与 Kitex 集成示例
func KitexConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 3. 添加配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		fmt.Printf("配置 [%s] 已更新，正在重新加载服务配置...\n", dataId)
		
		// 这里可以执行实际的配置重载逻辑
		// 例如：重新连接数据库、更新日志级别、重启服务等
		switch dataId {
		case "common":
			// 重新加载通用配置
			fmt.Println("  重新加载通用配置...")
			// 可以在这里调用你的服务重载逻辑
			
		case "pasetopub", "pasetosecret":
			// 重新加载 Paseto 配置
			fmt.Println("  重新加载 Paseto 配置...")
			// 可以在这里重新初始化 Paseto 相关组件
		}
	})

	// 4. 监听配置
	err = watcher.WatchConfig("common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	// 5. 模拟 Kitex 服务运行
	fmt.Println("Kitex 服务已启动，开始监听配置变化...")
	
	// 这里可以启动你的 Kitex 服务
	// 配置变化时会自动触发回调函数
	
	// 保持程序运行
	select {}
}

// 配置监听与数据库重连示例
func DatabaseConfigWatchExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 3. 添加数据库配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		if dataId == "common" {
			fmt.Printf("数据库配置已更新，新 DSN: %s\n", config.MySQL.DSN)
			
			// 这里可以执行数据库重连逻辑
			// 例如：关闭旧连接，使用新 DSN 创建新连接
			fmt.Println("  正在重新连接数据库...")
			// 你的数据库重连逻辑
			
			fmt.Printf("  Redis 配置已更新，新地址: %s\n", config.Redis.Address)
			fmt.Println("  正在重新连接 Redis...")
			// 你的 Redis 重连逻辑
		}
	})

	// 4. 监听配置
	err = watcher.WatchConfig("common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("数据库配置监听已启动，按 Ctrl+C 退出...")
	
	// 5. 保持程序运行
	select {}
}

// 配置监听与优雅关闭示例
func GracefulShutdownExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}

	// 3. 添加配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		fmt.Printf("配置 [%s] 已更新\n", dataId)
	})

	// 4. 监听配置
	err = watcher.WatchConfig("common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	// 5. 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 6. 在 goroutine 中运行监听
	go func() {
		defer watcher.Stop()
		select {
		case <-ctx.Done():
			fmt.Println("配置监听已停止")
		}
	}()

	fmt.Println("配置监听已启动，按 Ctrl+C 退出...")
	
	// 7. 等待中断信号
	select {
	case <-ctx.Done():
		fmt.Println("程序正在关闭...")
	}
}

// 配置监听与错误处理示例
func ErrorHandlingExample() {
	// 1. 设置环境变量
	os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
	os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
	os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

	// 2. 创建配置监听器
	watcher, err := NewConfigWatcher(
		[]string{"127.0.0.1:8848"},
		"your-namespace-id",
		"DEFAULT_GROUP",
	)
	if err != nil {
		log.Fatalf("创建配置监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 3. 添加带错误处理的配置变化回调
	watcher.AddCallback(func(dataId string, config *kvconfig.CommonConfig) {
		fmt.Printf("配置 [%s] 已更新\n", dataId)
		
		// 模拟配置处理可能出现的错误
		if config.Kitex.Service == "" {
			fmt.Printf("  警告: 配置 [%s] 中 Kitex 服务名为空\n", dataId)
		}
		
		if config.MySQL.DSN == "" {
			fmt.Printf("  警告: 配置 [%s] 中 MySQL DSN 为空\n", dataId)
		}
		
		// 验证配置有效性
		if config.Kitex.Address == "" {
			fmt.Printf("  错误: 配置 [%s] 中 Kitex 地址为空，使用默认值\n", dataId)
			config.Kitex.Address = "0.0.0.0:8080"
		}
	})

	// 4. 监听配置
	err = watcher.WatchConfig("common", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("监听配置失败: %v", err)
	}

	fmt.Println("配置监听已启动（带错误处理），按 Ctrl+C 退出...")
	
	// 5. 保持程序运行
	select {}
}
