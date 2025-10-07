package examples

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grayscalecloud/kitexcommon/kvconfig"
)

// ConsulWatchExample Consul 配置监听示例
func ConsulWatchExample() {
	// 1. 设置环境变量
	os.Setenv("CONSUL_SERVER_ADDR", "127.0.0.1:8500")
	os.Setenv("CONSUL_NAMESPACE_ID", "test-namespace")
	os.Setenv("CONSUL_GROUP", "DEFAULT_GROUP")

	// 2. 创建 Consul 客户端
	client, err := kvconfig.NewConsulConfigClient(
		"127.0.0.1:8500",
		"test-namespace",
		"DEFAULT_GROUP",
		"",
		"",
	)
	if err != nil {
		log.Fatalf("创建 Consul 客户端失败: %v", err)
	}
	defer client.Close()

	// 3. 使用 channel 方式监听配置
	fmt.Println("=== 使用 Channel 方式监听配置 ===")
	configChan, err := client.WatchConfig("test-config", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("启动配置监听失败: %v", err)
	}

	// 启动一个 goroutine 处理配置变化
	go func() {
		for config := range configChan {
			if config == "" {
				fmt.Println("配置被删除")
			} else {
				fmt.Printf("配置发生变化: %s\n", config)
			}
		}
		fmt.Println("配置监听已停止")
	}()

	// 4. 使用回调函数方式监听配置
	fmt.Println("\n=== 使用回调函数方式监听配置 ===")
	err = client.ListenConfig("another-config", "DEFAULT_GROUP", func(content string) {
		if content == "" {
			fmt.Println("另一个配置被删除")
		} else {
			fmt.Printf("另一个配置发生变化: %s\n", content)
		}
	})
	if err != nil {
		log.Printf("启动另一个配置监听失败: %v", err)
	}

	// 5. 模拟配置变化（这里只是演示，实际配置变化由 Consul 服务器触发）
	fmt.Println("\n=== 模拟配置变化 ===")

	// 发布一个测试配置
	testConfig := `database:
  host: "localhost"
  port: 3306
  username: "test"
  password: "test123"

redis:
  address: "localhost:6379"
  password: ""
  db: 0`

	err = client.PublishConfig("test-config", "DEFAULT_GROUP", testConfig)
	if err != nil {
		log.Printf("发布测试配置失败: %v", err)
	} else {
		fmt.Println("✅ 测试配置发布成功")
	}

	// 等待一段时间让监听生效
	time.Sleep(2 * time.Second)

	// 更新配置
	updatedConfig := `database:
  host: "localhost"
  port: 3306
  username: "test"
  password: "newpassword123"

redis:
  address: "localhost:6379"
  password: "redis123"
  db: 1`

	err = client.PublishConfig("test-config", "DEFAULT_GROUP", updatedConfig)
	if err != nil {
		log.Printf("更新测试配置失败: %v", err)
	} else {
		fmt.Println("✅ 测试配置更新成功")
	}

	// 等待一段时间让监听生效
	time.Sleep(2 * time.Second)

	// 6. 停止监听
	fmt.Println("\n=== 停止监听 ===")
	err = client.StopListenConfig("test-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("停止配置监听失败: %v", err)
	} else {
		fmt.Println("✅ 配置监听已停止")
	}

	err = client.StopListenConfig("another-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("停止另一个配置监听失败: %v", err)
	} else {
		fmt.Println("✅ 另一个配置监听已停止")
	}

	// 等待一段时间确保停止生效
	time.Sleep(1 * time.Second)
	fmt.Println("=== 示例完成 ===")
}

// ConsulWatchWithContext 使用 Context 控制监听生命周期
func ConsulWatchWithContext() {
	// 1. 创建 Consul 客户端
	client, err := kvconfig.NewConsulConfigClient(
		"127.0.0.1:8500",
		"test-namespace",
		"DEFAULT_GROUP",
		"",
		"",
	)
	if err != nil {
		log.Fatalf("创建 Consul 客户端失败: %v", err)
	}
	defer client.Close()

	// 2. 创建 Context 用于控制监听生命周期
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. 启动配置监听
	configChan, err := client.WatchConfig("context-config", "DEFAULT_GROUP")
	if err != nil {
		log.Fatalf("启动配置监听失败: %v", err)
	}

	// 4. 使用 select 处理配置变化和 Context 取消
	go func() {
		for {
			select {
			case config, ok := <-configChan:
				if !ok {
					fmt.Println("配置通道已关闭")
					return
				}
				if config == "" {
					fmt.Println("配置被删除")
				} else {
					fmt.Printf("配置发生变化: %s\n", config)
				}
			case <-ctx.Done():
				fmt.Println("Context 已取消，停止监听")
				client.StopListenConfig("context-config", "DEFAULT_GROUP")
				return
			}
		}
	}()

	// 5. 等待 Context 超时或取消
	<-ctx.Done()
	fmt.Println("=== Context 控制示例完成 ===")
}

// ConsulMultipleWatch 监听多个配置
func ConsulMultipleWatch() {
	// 1. 创建 Consul 客户端
	client, err := kvconfig.NewConsulConfigClient(
		"127.0.0.1:8500",
		"test-namespace",
		"DEFAULT_GROUP",
		"",
		"",
	)
	if err != nil {
		log.Fatalf("创建 Consul 客户端失败: %v", err)
	}
	defer client.Close()

	// 2. 监听多个配置
	configs := []string{"config1", "config2", "config3"}
	configChans := make(map[string]<-chan string)

	for _, configName := range configs {
		configChan, err := client.WatchConfig(configName, "DEFAULT_GROUP")
		if err != nil {
			log.Printf("启动配置 %s 监听失败: %v", configName, err)
			continue
		}
		configChans[configName] = configChan
		fmt.Printf("✅ 开始监听配置: %s\n", configName)
	}

	// 3. 使用 goroutine 处理所有配置变化
	for configName, configChan := range configChans {
		go func(name string, ch <-chan string) {
			for config := range ch {
				if config == "" {
					fmt.Printf("配置 %s 被删除\n", name)
				} else {
					fmt.Printf("配置 %s 发生变化: %s\n", name, config)
				}
			}
			fmt.Printf("配置 %s 监听已停止\n", name)
		}(configName, configChan)
	}

	// 4. 等待一段时间
	time.Sleep(10 * time.Second)

	// 5. 停止所有监听
	fmt.Println("\n=== 停止所有监听 ===")
	err = client.StopAllListenConfigs()
	if err != nil {
		log.Printf("停止所有监听失败: %v", err)
	} else {
		fmt.Println("✅ 所有配置监听已停止")
	}

	fmt.Println("=== 多配置监听示例完成 ===")
}
