# KitexCommon 配置中心集成

本包提供了多种配置中心的集成支持，包括 Consul 和 Nacos v2。

## 功能特性

- 支持 Consul 和 Nacos v2 配置中心
- 统一的配置获取接口
- 配置变化监听
- 配置发布和管理
- 批量配置操作
- 重试机制
- 上下文支持

## 快速开始

### 使用 Nacos 配置中心

#### 1. 基本使用

```go
package main

import (
    "log"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 创建 Nacos 配置客户端
    client, err := kvconfig.NewNacosConfigClient(
        []string{"127.0.0.1:8848"}, // Nacos 服务器地址
        "your-namespace-id",        // 命名空间 ID
        "DEFAULT_GROUP",            // 配置分组
    )
    if err != nil {
        log.Fatalf("创建 Nacos 客户端失败: %v", err)
    }
    defer client.Close()

    // 获取通用配置
    commonConfig, err := client.GetCommonConfig("DEFAULT_GROUP")
    if err != nil {
        log.Printf("获取配置失败: %v", err)
        return
    }
    
    log.Printf("配置: %+v", commonConfig)
}
```

#### 2. 使用配置工厂（推荐）

```go
package main

import (
    "log"
    "os"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 设置环境变量
    os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
    os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
    os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

    // 初始化全局配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeNacos)
    if err != nil {
        log.Fatalf("初始化配置工厂失败: %v", err)
    }

    // 使用全局函数获取配置
    commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Printf("获取配置失败: %v", err)
        return
    }
    
    log.Printf("配置: %+v", commonConfig)
}
```

### 使用 Consul 配置中心

```go
package main

import (
    "log"
    "os"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 设置环境变量
    os.Setenv("REGISTRY_ADDRESS", "127.0.0.1:8500")

    // 初始化全局配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeConsul)
    if err != nil {
        log.Fatalf("初始化配置工厂失败: %v", err)
    }

    // 使用全局函数获取配置
    commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Printf("获取配置失败: %v", err)
        return
    }
    
    log.Printf("配置: %+v", commonConfig)
}
```

## 配置结构

### 通用配置结构

```yaml
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
```

### Paseto 配置结构

```yaml
# pasetopub 配置
pub_key: "your-public-key"
implicit: "your-implicit"

# pasetosecret 配置
secret_key: "your-secret-key"
implicit: "your-implicit"
```

## 高级功能

### 配置监听

```go
// 监听单个配置
err := client.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
    log.Printf("配置发生变化: %s", content)
})

// 批量监听配置
configs := []kvconfig.ConfigRequest{
    {Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
    {Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
}

err = client.WatchConfigs(configs, func(key, content string) {
    log.Printf("配置 [%s] 发生变化: %s", key, content)
})
```

### 配置管理

```go
// 发布配置
configContent := `kitex:
  service: "example-service"
  address: "0.0.0.0:8080"`

err := client.PublishConfig("example-config", "DEFAULT_GROUP", configContent)

// 删除配置
err = client.DeleteConfig("example-config", "DEFAULT_GROUP")
```

### 重试机制

```go
// 带重试的配置获取
content, err := client.GetConfigWithRetry("common", "DEFAULT_GROUP", 3, 2*time.Second)
```

### 上下文支持

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

content, err := client.GetConfigWithContext(ctx, "common", "DEFAULT_GROUP")
```

## 环境变量

### Nacos 配置

- `NACOS_SERVER_ADDR`: Nacos 服务器地址，多个地址用逗号分隔
- `NACOS_NAMESPACE_ID`: 命名空间 ID
- `NACOS_GROUP`: 配置分组

### Consul 配置

- `REGISTRY_ADDRESS`: Consul 服务器地址
- `REGISTRY_ADDRESS_USERNAME`: Consul 用户名（可选）
- `REGISTRY_ADDRESS_PASSWORD`: Consul 密码（可选）

## 示例

完整的使用示例请参考 `examples/nacos_example.go` 文件。

## 注意事项

1. 确保 Nacos 或 Consul 服务正常运行
2. 配置内容必须是有效的 YAML 格式
3. 监听配置变化时，程序需要保持运行状态
4. 建议在生产环境中使用配置工厂模式
5. 配置获取失败时会记录错误日志，但不会 panic
