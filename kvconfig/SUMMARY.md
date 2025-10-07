# Nacos v2 配置中心集成总结

## 已完成的功能

✅ **Nacos v2 配置中心客户端**
- 支持多服务器地址配置
- 支持命名空间和分组管理
- 支持配置获取、发布、删除
- 支持配置变化监听
- 支持批量配置操作
- 支持重试机制和上下文控制

✅ **配置工厂模式**
- 统一的配置获取接口
- 支持 Consul 和 Nacos 两种配置中心
- 全局配置工厂实例
- 环境变量自动配置

✅ **配置监听功能**
- 单个配置监听
- 多个配置批量监听
- 配置变化回调处理
- 线程安全的配置缓存

✅ **与现有代码兼容**
- 保持与现有 Consul 配置接口兼容
- 支持相同的配置结构体
- 支持 Paseto 配置获取

✅ **错误处理和日志**
- 完整的错误处理机制
- 详细的日志记录
- 优雅的错误恢复

✅ **测试和示例**
- 完整的单元测试
- 多种使用示例
- 详细的使用文档

## 文件结构

```
kvconfig/
├── nacos_config.go              # Nacos 配置客户端核心实现
├── config_factory.go            # 配置工厂，支持多种配置中心
├── conf.go                      # 原有 Consul 配置实现
├── consul_config.go             # 原有 Consul 配置实现
├── nacos_config_test.go         # 单元测试
├── README.md                    # 基本使用说明
├── CONFIG_WATCH_GUIDE.md        # 配置监听详细指南
├── SUMMARY.md                   # 本总结文档
└── examples/
    ├── nacos_example.go         # Nacos 基本使用示例
    ├── config_watch_example.go  # 配置监听高级示例
    ├── simple_watch_example.go  # 简单配置监听示例
    ├── kitex_integration_example.go # Kitex 集成示例
    └── quick_start.go           # 快速开始示例
```

## 快速使用

### 1. 基本配置获取

```go
// 设置环境变量
os.Setenv("NACOS_SERVER_ADDR", "127.0.0.1:8848")
os.Setenv("NACOS_NAMESPACE_ID", "your-namespace-id")
os.Setenv("NACOS_GROUP", "DEFAULT_GROUP")

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

### 3. 多配置监听

```go
// 定义要监听的配置
configs := []kvconfig.ConfigRequest{
    {Key: "common", DataId: "common", Group: "DEFAULT_GROUP"},
    {Key: "pasetopub", DataId: "pasetopub", Group: "DEFAULT_GROUP"},
}

// 批量监听
err = nacosClient.WatchConfigs(configs, func(key, content string) {
    log.Printf("配置 [%s] 已更新", key)
})
```

## 环境变量配置

```bash
# Nacos 服务器地址
export NACOS_SERVER_ADDR="127.0.0.1:8848"

# 命名空间 ID
export NACOS_NAMESPACE_ID="your-namespace-id"

# 配置分组
export NACOS_GROUP="DEFAULT_GROUP"
```

## 配置格式

在 Nacos 中创建配置时，使用以下 YAML 格式：

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

## 在 Kitex 项目中使用

1. **服务启动时初始化配置监听**
2. **配置变化时重新加载相关组件**
3. **支持数据库连接重连**
4. **支持 Paseto 配置热更新**
5. **支持日志级别动态调整**

## 测试

运行测试：
```bash
go test ./kvconfig -v
```

运行示例：
```bash
go run kvconfig/examples/quick_start.go
```

## 注意事项

1. 确保 Nacos 服务正常运行
2. 配置内容必须是有效的 YAML 格式
3. 监听配置变化时程序需要保持运行
4. 建议在生产环境中使用配置工厂模式
5. 配置获取失败时会记录错误日志，但不会 panic

## 下一步

- 可以添加配置加密支持
- 可以添加配置版本管理
- 可以添加配置回滚功能
- 可以添加配置变更审计日志
