# Nacos 配置监听使用指南

本指南将详细说明如何在你的 Kitex 项目中监听 Nacos 配置中心的配置更新。

## 快速开始

### 1. 基本配置监听

```go
package main

import (
    "log"
    "os"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 1. 设置环境变量
    os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
    os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
    os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

    // 2. 初始化配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeNacos)
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
        log.Printf("配置已更新: %s", content)
        
        // 重新获取解析后的配置
        commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
        if err != nil {
            log.Printf("重新获取配置失败: %v", err)
            return
        }
        
        // 使用新配置
        log.Printf("Kitex 服务: %s", commonConfig.Kitex.Service)
        log.Printf("Kitex 地址: %s", commonConfig.Kitex.Address)
    })
    
    if err != nil {
        log.Fatalf("监听配置失败: %v", err)
    }

    // 5. 保持程序运行
    select {}
}
```

### 2. 监听多个配置

```go
// 定义要监听的配置
configs := []kvconfig.ConfigRequest{
    {Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
    {Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
    {Key: "pasetosecret", DataId: "pasetosecret", Group: "DEFAULT_GROUP"},
}

// 批量监听配置
err = nacosClient.WatchConfigs(configs, func(key, content string) {
    log.Printf("配置 [%s] 已更新: %s", key, content)
    
    // 根据配置类型执行不同操作
    switch key {
    case "common":
        log.Println("通用配置已更新，可能需要重启服务")
    case "pasetopub":
        log.Println("Paseto 公钥已更新")
    case "pasetosecret":
        log.Println("Paseto 密钥已更新")
    }
})
```

## 在 Kitex 服务中集成配置监听

### 1. 服务启动时初始化配置监听

```go
package main

import (
    "log"
    "os"
    "github.com/cloudwego/kitex/server"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 1. 设置环境变量
    os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
    os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
    os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

    // 2. 初始化配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeNacos)
    if err != nil {
        log.Fatalf("初始化配置工厂失败: %v", err)
    }

    // 3. 获取初始配置
    commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Fatalf("获取初始配置失败: %v", err)
    }

    // 4. 启动配置监听
    startConfigWatch()

    // 5. 创建并启动 Kitex 服务
    svr := server.NewServer(
        server.WithServiceAddr(commonConfig.Kitex.Address),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: commonConfig.Kitex.Service,
        }),
    )

    err = svr.Run()
    if err != nil {
        log.Fatalf("服务器启动失败: %v", err)
    }
}

func startConfigWatch() {
    nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()
    if nacosClient == nil {
        log.Fatalf("获取 Nacos 客户端失败")
    }

    // 监听通用配置
    err := nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
        log.Println("配置已更新，正在重新加载...")
        
        // 重新获取配置
        newConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
        if err != nil {
            log.Printf("重新获取配置失败: %v", err)
            return
        }
        
        // 更新全局配置变量
        updateGlobalConfig(newConfig)
    })
    
    if err != nil {
        log.Fatalf("监听配置失败: %v", err)
    }
}

var globalConfig *model.CommonConfig

func updateGlobalConfig(config *model.CommonConfig) {
    globalConfig = config
    log.Printf("全局配置已更新: Kitex 服务=%s, 地址=%s", 
        config.Kitex.Service, config.Kitex.Address)
}
```

### 2. 配置变化时的处理逻辑

```go
// 数据库配置变化处理
func handleDatabaseConfigChange(config *model.CommonConfig) {
    log.Printf("数据库配置已更新，新 DSN: %s", config.MySQL.DSN)
    
    // 关闭旧连接
    if oldDB != nil {
        oldDB.Close()
    }
    
    // 使用新 DSN 创建新连接
    newDB, err := sql.Open("mysql", config.MySQL.DSN)
    if err != nil {
        log.Printf("创建新数据库连接失败: %v", err)
        return
    }
    
    oldDB = newDB
    log.Println("数据库连接已更新")
}

// Redis 配置变化处理
func handleRedisConfigChange(config *model.CommonConfig) {
    log.Printf("Redis 配置已更新，新地址: %s", config.Redis.Address)
    
    // 关闭旧连接
    if oldRedis != nil {
        oldRedis.Close()
    }
    
    // 使用新配置创建新连接
    newRedis := redis.NewClient(&redis.Options{
        Addr:     config.Redis.Address,
        Username: config.Redis.Username,
        Password: config.Redis.Password,
        DB:       config.Redis.DB,
    })
    
    oldRedis = newRedis
    log.Println("Redis 连接已更新")
}

// Paseto 配置变化处理
func handlePasetoConfigChange() {
    log.Println("Paseto 配置已更新")
    
    // 重新获取 Paseto 配置
    pubConfig, err := kvconfig.GetPasetoPubConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Printf("获取 Paseto 公钥配置失败: %v", err)
        return
    }
    
    secretConfig, err := kvconfig.GetPasetoSecretConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Printf("获取 Paseto 密钥配置失败: %v", err)
        return
    }
    
    // 更新全局 Paseto 配置
    updatePasetoConfig(pubConfig, secretConfig)
    log.Println("Paseto 配置已更新")
}
```

## 高级用法

### 1. 使用配置监听器

```go
// 创建配置监听器
watcher, err := NewConfigWatcher(
    []string{"127.0.0.1:8848"},
    "your-namespace-id",
    "DEFAULT_GROUP",
)
if err != nil {
    log.Fatalf("创建配置监听器失败: %v", err)
}
defer watcher.Stop()

// 添加配置变化回调
watcher.AddCallback(func(dataId string, config *model.CommonConfig) {
    log.Printf("配置 [%s] 已更新", dataId)
    
    // 根据不同配置类型执行不同操作
    switch dataId {
    case "common":
        handleCommonConfigChange(config)
    case "pasetopub":
        handlePasetoConfigChange()
    }
})

// 监听配置
err = watcher.WatchConfig("common", "DEFAULT_GROUP")
if err != nil {
    log.Fatalf("监听配置失败: %v", err)
}
```

### 2. 带重试的配置监听

```go
// 带重试的配置获取
content, err := nacosClient.GetConfigWithRetry("common", "DEFAULT_GROUP", 3, 2*time.Second)
if err != nil {
    log.Printf("获取配置失败: %v", err)
} else {
    log.Printf("获取到的配置: %s", content)
}
```

### 3. 带上下文的配置监听

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

content, err := nacosClient.GetConfigWithContext(ctx, "common", "DEFAULT_GROUP")
if err != nil {
    log.Printf("获取配置失败: %v", err)
} else {
    log.Printf("获取到的配置: %s", content)
}
```

## 环境变量配置

### Nacos 配置

```bash
# Nacos 服务器地址（多个地址用逗号分隔）
export NACOS_SERVER_ADDR="127.0.0.1:8848,127.0.0.1:8849"

# 命名空间 ID
export NACOS_NAMESPACE_ID="your-namespace-id"

# 配置分组
export NACOS_GROUP="DEFAULT_GROUP"
```

### 在 Docker 中使用

```dockerfile
ENV NACOS_SERVER_ADDR=127.0.0.1:8848
ENV NACOS_NAMESPACE_ID=your-namespace-id
ENV NACOS_GROUP=DEFAULT_GROUP
```

### 在 Kubernetes 中使用

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: your-service
spec:
  template:
    spec:
      containers:
      - name: your-service
        env:
        - name: NACOS_SERVER_ADDR
          value: "nacos-service:8848"
        - name: NACOS_NAMESPACE_ID
          value: "your-namespace-id"
        - name: NACOS_GROUP
          value: "DEFAULT_GROUP"
```

## 注意事项

1. **配置格式**: 确保 Nacos 中的配置是有效的 YAML 格式
2. **错误处理**: 配置监听失败时不会影响服务正常运行，但会记录错误日志
3. **性能考虑**: 避免在配置变化回调中执行耗时操作
4. **线程安全**: 配置变化回调可能在不同 goroutine 中执行，注意线程安全
5. **优雅关闭**: 程序退出时应该停止配置监听

## 故障排除

### 常见问题

1. **连接失败**: 检查 Nacos 服务器地址和端口是否正确
2. **配置解析失败**: 检查配置内容是否为有效的 YAML 格式
3. **监听不生效**: 检查命名空间 ID 和分组名称是否正确
4. **权限问题**: 检查是否有访问配置的权限

### 调试方法

```go
// 启用详细日志
os.Setenv("NACOS_LOG_LEVEL", "debug")

// 检查配置是否存在
content, err := nacosClient.GetConfig("common", "DEFAULT_GROUP")
if err != nil {
    log.Printf("配置不存在或获取失败: %v", err)
} else {
    log.Printf("配置内容: %s", content)
}
```

## 完整示例

参考 `examples/` 目录下的示例文件：

- `simple_watch_example.go`: 简单配置监听示例
- `config_watch_example.go`: 高级配置监听示例
- `kitex_integration_example.go`: Kitex 集成示例
