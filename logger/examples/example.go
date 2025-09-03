package main

import (
	"context"

	"github.com/grayscalecloud/kitexcommon/logger"
)

func main() {
	// 创建默认配置的日志器
	log := logger.NewLogger(nil)

	// 创建自定义配置的日志器
	customLog := logger.NewLogger(&logger.Config{
		Level:      logger.DebugLevel,
		OutputPath: "app.log", // 输出到文件
		Format:     "text",    // 使用文本格式
	})

	// 使用上下文
	ctx := context.Background()

	// 记录不同级别的日志
	log.Info(ctx, "这是一条信息日志")
	log.Warn(ctx, "这是一条警告日志，参数: %s", "示例参数")
	log.Error(ctx, "这是一条错误日志")

	// 使用字段
	logWithFields := customLog.WithField("module", "user-service").
		WithField("requestID", "12345")

	logWithFields.Debug(ctx, "这是一条带有字段的调试日志")
	logWithFields.Info(ctx, "用户登录成功，用户ID: %d", 10001)
}
