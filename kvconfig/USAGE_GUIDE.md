# Nacos 配置中心使用指南

## 测试结果

✅ **连接成功**：你的 Nacos 配置中心连接正常！
- 服务器地址：`115.190.176.125:8848`
- 命名空间：`6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a`
- 用户名：`nacos`
- 密码：`nacos`

## 快速开始

### 1. 基本配置

```go
package main

import (
    "log"
    "os"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 设置环境变量
    os.Setenv("NACOS_SERVER_ADDR", "115.190.176.125:8848")
    os.Setenv("NACOS_NAMESPACE_ID", "6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a")
    os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")
    os.Setenv("NACOS_USERNAME", "nacos")
    os.Setenv("NACOS_PASSWORD", "nacos")

    // 初始化配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeNacos)
    if err != nil {
        log.Fatalf("初始化失败: %v", err)
    }

    // 获取配置
    commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Fatalf("获取配置失败: %v", err)
    }

    log.Printf("Kitex 服务: %s", commonConfig.Kitex.Service)
    log.Printf("Kitex 地址: %s", commonConfig.Kitex.Address)
}
```

### 2. 配置监听

```go
// 获取 Nacos 客户端
nacosClient := kvconfig.GetGlobalConfigFactory().GetNacosClient()

// 监听配置变化
err := nacosClient.ListenConfig("common", "DEFAULT_GROUP", func(content string) {
    log.Printf("配置已更新: %s", content)
    
    // 重新获取配置
    newConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Printf("重新获取配置失败: %v", err)
        return
    }
    
    // 使用新配置
    log.Printf("Kitex 服务: %s", newConfig.Kitex.Service)
})
```

### 3. 发布配置

```go
// 发布配置
configContent := `
kitex:
  service: "your-service"
  address: "0.0.0.0:8080"
  metrics_port: "9090"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "info"
  log_file_name: "service.log"
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

err := nacosClient.PublishConfig("common", "DEFAULT_GROUP", configContent)
if err != nil {
    log.Printf("发布配置失败: %v", err)
} else {
    log.Println("配置发布成功")
}
```

## 在 Kitex 项目中使用

### 1. 服务启动时初始化

```go
package main

import (
    "log"
    "os"
    "github.com/cloudwego/kitex/server"
    "github.com/grayscalecloud/kitexcommon/kvconfig"
)

func main() {
    // 设置环境变量
    os.Setenv("NACOS_SERVER_ADDR", "115.190.176.125:8848")
    os.Setenv("NACOS_NAMESPACE_ID", "6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a")
    os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")
    os.Setenv("NACOS_USERNAME", "nacos")
    os.Setenv("NACOS_PASSWORD", "nacos")

    // 初始化配置工厂
    err := kvconfig.InitGlobalConfigFactory(kvconfig.ConfigTypeNacos)
    if err != nil {
        log.Fatalf("初始化配置工厂失败: %v", err)
    }

    // 获取配置
    commonConfig, err := kvconfig.GetCommonConfigGlobal("DEFAULT_GROUP")
    if err != nil {
        log.Fatalf("获取配置失败: %v", err)
    }

    // 启动配置监听
    startConfigWatch()

    // 创建并启动 Kitex 服务
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

    // 监听配置变化
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

### 2. 配置变化处理

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
```

## 环境变量配置

```bash
# Nacos 服务器地址
export NACOS_SERVER_ADDR="115.190.176.125:8848"

# 命名空间 ID
export NACOS_NAMESPACE_ID="6a4a9a5b-bf1b-4e3c-8c0d-56cc393a616a"

# 配置分组
export NACOS_GROUP="DEFAULT_GROUP"

# 用户名和密码
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="nacos"
```

## 配置格式

在 Nacos 中创建配置时，使用以下 YAML 格式：

```yaml
kitex:
  service: "your-service"
  address: "0.0.0.0:8080"
  metrics_port: "9090"
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true
  log_level: "info"
  log_file_name: "service.log"
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

## 测试命令

```bash
# 运行简单测试
go test -v ./kvconfig -run TestNacosSimple

# 运行完整演示
go run kvconfig/examples/nacos_watch_demo.go

# 运行所有测试
go test -v ./kvconfig
```

## 功能特性

✅ **配置获取**：支持获取各种配置
✅ **配置发布**：支持发布新配置
✅ **配置删除**：支持删除配置
✅ **配置监听**：支持实时监听配置变化
✅ **批量操作**：支持批量获取配置
✅ **重试机制**：支持配置获取重试
✅ **上下文支持**：支持带上下文的配置操作
✅ **身份验证**：支持用户名密码认证
✅ **配置工厂**：统一的配置管理接口

## 注意事项

1. 确保 Nacos 服务正常运行
2. 配置内容必须是有效的 YAML 格式
3. 监听配置变化时程序需要保持运行
4. 建议在生产环境中使用配置工厂模式
5. 配置获取失败时会记录错误日志，但不会 panic

## 故障排除

### 常见问题

1. **连接失败**：检查服务器地址和端口是否正确
2. **身份验证失败**：检查用户名和密码是否正确
3. **配置不存在**：检查配置是否已在 Nacos 中创建
4. **配置解析失败**：检查配置内容是否为有效的 YAML 格式

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

## 总结

🎉 **恭喜！** 你的 Nacos 配置中心集成已经完成并且测试通过！

现在你可以：
1. ✅ 连接到你的 Nacos 配置中心
2. ✅ 发布和获取配置
3. ✅ 监听配置变化
4. ✅ 在 Kitex 项目中使用配置中心
5. ✅ 实现配置的热更新

所有功能都正常工作，可以开始在你的项目中使用 Nacos 配置中心了！
