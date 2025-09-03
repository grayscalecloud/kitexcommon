// Package logger 提供了一个通用的日志接口和实现
// 支持不同级别的日志记录和格式化
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	// DebugLevel 调试级别
	DebugLevel LogLevel = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 定义日志接口
type Logger interface {
	// Debug 记录调试级别的日志
	Debug(ctx context.Context, format string, args ...interface{})
	// Info 记录信息级别的日志
	Info(ctx context.Context, format string, args ...interface{})
	// Warn 记录警告级别的日志
	Warn(ctx context.Context, format string, args ...interface{})
	// Error 记录错误级别的日志
	Error(ctx context.Context, format string, args ...interface{})
	// WithField 添加字段到日志上下文
	WithField(key string, value interface{}) Logger
}

// Config 日志配置
type Config struct {
	// Level 日志级别
	Level LogLevel
	// OutputPath 输出路径，为空则输出到标准输出
	OutputPath string
	// Format 日志格式，支持 "text" 和 "json"
	Format string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		OutputPath: "",
		Format:     "text",
	}
}

// standardLogger 标准日志实现
type standardLogger struct {
	level  LogLevel
	writer io.Writer
	format string
	fields map[string]interface{}
}

// NewLogger 创建新的日志器
func NewLogger(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}

	var writer io.Writer
	if config.OutputPath == "" {
		writer = os.Stdout
	} else {
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file: %v, using stdout instead\n", err)
			writer = os.Stdout
		} else {
			writer = file
		}
	}

	return &standardLogger{
		level:  config.Level,
		writer: writer,
		format: config.Format,
		fields: make(map[string]interface{}),
	}
}

// log 记录日志
func (l *standardLogger) log(ctx context.Context, level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)

	var logLine string
	if l.format == "json" {
		// 简单的JSON格式
		fields := ""
		for k, v := range l.fields {
			fields += fmt.Sprintf(`"%s":%#v,`, k, v)
		}
		logLine = fmt.Sprintf(`{"time":"%s","level":"%s","message":"%s",%s}`,
			timestamp, level.String(), message, fields[:len(fields)-1])
	} else {
		// 文本格式
		fields := ""
		for k, v := range l.fields {
			fields += fmt.Sprintf("%s=%v ", k, v)
		}
		logLine = fmt.Sprintf("%s [%s] %s %s\n", timestamp, level.String(), message, fields)
	}

	_, _ = l.writer.Write([]byte(logLine))
}

// Debug 实现Debug日志
func (l *standardLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, args...)
}

// Info 实现Info日志
func (l *standardLogger) Info(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, InfoLevel, format, args...)
}

// Warn 实现Warn日志
func (l *standardLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WarnLevel, format, args...)
}

// Error 实现Error日志
func (l *standardLogger) Error(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ErrorLevel, format, args...)
}

// WithField 添加字段到日志上下文
func (l *standardLogger) WithField(key string, value interface{}) Logger {
	newLogger := &standardLogger{
		level:  l.level,
		writer: l.writer,
		format: l.format,
		fields: make(map[string]interface{}, len(l.fields)+1),
	}

	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// 添加新字段
	newLogger.fields[key] = value
	return newLogger
}
