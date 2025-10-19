// Package logger 提供了一个通用的日志接口和实现
// 基于 Kitex 的 klog 接口，支持不同级别的日志记录和格式化
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

// levelString 返回 klog.Level 的字符串表示
func levelString(level klog.Level) string {
	switch level {
	case klog.LevelTrace:
		return "TRACE"
	case klog.LevelDebug:
		return "DEBUG"
	case klog.LevelInfo:
		return "INFO"
	case klog.LevelNotice:
		return "NOTICE"
	case klog.LevelWarn:
		return "WARN"
	case klog.LevelError:
		return "ERROR"
	case klog.LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 定义日志接口，基于 Kitex 的 klog.FullLogger
type Logger interface {
	klog.FullLogger
	// WithField 添加字段到日志上下文
	WithField(key string, value interface{}) Logger
	// Close 关闭日志器，释放资源
	Close() error
}

// Config 日志配置
type Config struct {
	// Level 日志级别
	Level klog.Level
	// OutputPath 输出路径，为空则输出到标准输出
	OutputPath string
	// Format 日志格式，支持 "text" 和 "json"
	Format string
	// MaxSize 单个日志文件的最大大小（字节），0 表示不限制
	MaxSize int64
	// MaxBackups 保留的旧日志文件数量，0 表示不删除
	MaxBackups int
	// MaxAge 日志文件保留天数，0 表示不删除
	MaxAge int
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      klog.LevelInfo,
		OutputPath: "",
		Format:     "text",
		MaxSize:    0, // 不限制大小
		MaxBackups: 0, // 不删除旧文件
		MaxAge:     0, // 不按时间删除
	}
}

// standardLogger 标准日志实现
type standardLogger struct {
	level      klog.Level
	writer     io.Writer
	format     string
	fields     map[string]interface{}
	file       *os.File   // 保存文件句柄，用于关闭
	config     *Config    // 保存配置，用于轮转
	mu         sync.Mutex // 保护轮转操作
	rotateTime time.Time  // 下次轮转时间
}

// NewLogger 创建新的日志器
func NewLogger(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}

	var writer io.Writer
	var file *os.File
	if config.OutputPath == "" {
		writer = os.Stdout
	} else {
		var err error
		file, err = os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file: %v, using stdout instead\n", err)
			writer = os.Stdout
			file = nil
		} else {
			writer = file
		}
	}

	return &standardLogger{
		level:      config.Level,
		writer:     writer,
		format:     config.Format,
		fields:     make(map[string]interface{}),
		file:       file,
		config:     config,
		rotateTime: time.Now().Add(24 * time.Hour), // 默认每天轮转一次
	}
}

// log 记录日志
func (l *standardLogger) log(level klog.Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// 检查是否需要轮转
	l.mu.Lock()
	if l.file != nil && l.config.MaxSize > 0 {
		if l.shouldRotate() {
			l.rotate()
		}
	}
	l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)

	var logLine string
	if l.format == "json" {
		// 简单的JSON格式
		fields := ""
		if len(l.fields) > 0 {
			for k, v := range l.fields {
				fields += fmt.Sprintf(`"%s":%#v,`, k, v)
			}
			// 移除最后一个逗号
			fields = fields[:len(fields)-1]
		}

		if fields != "" {
			logLine = fmt.Sprintf(`{"time":"%s","level":"%s","message":"%s",%s}`,
				timestamp, levelString(level), message, fields)
		} else {
			logLine = fmt.Sprintf(`{"time":"%s","level":"%s","message":"%s"}`,
				timestamp, levelString(level), message)
		}
	} else {
		// 文本格式
		fields := ""
		for k, v := range l.fields {
			fields += fmt.Sprintf("%s=%v ", k, v)
		}
		logLine = fmt.Sprintf("%s [%s] %s %s\n", timestamp, levelString(level), message, fields)
	}

	if _, err := l.writer.Write([]byte(logLine)); err != nil {
		// 如果写入失败，尝试输出到标准错误
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
	}
}

// 实现 klog.Logger 接口
func (l *standardLogger) Trace(v ...interface{}) {
	l.log(klog.LevelTrace, "%v", v...)
}

func (l *standardLogger) Debug(v ...interface{}) {
	l.log(klog.LevelDebug, "%v", v...)
}

func (l *standardLogger) Info(v ...interface{}) {
	l.log(klog.LevelInfo, "%v", v...)
}

func (l *standardLogger) Notice(v ...interface{}) {
	l.log(klog.LevelNotice, "%v", v...)
}

func (l *standardLogger) Warn(v ...interface{}) {
	l.log(klog.LevelWarn, "%v", v...)
}

func (l *standardLogger) Error(v ...interface{}) {
	l.log(klog.LevelError, "%v", v...)
}

func (l *standardLogger) Fatal(v ...interface{}) {
	l.log(klog.LevelFatal, "%v", v...)
}

// 实现 klog.FormatLogger 接口
func (l *standardLogger) Tracef(format string, v ...interface{}) {
	l.log(klog.LevelTrace, format, v...)
}

func (l *standardLogger) Debugf(format string, v ...interface{}) {
	l.log(klog.LevelDebug, format, v...)
}

func (l *standardLogger) Infof(format string, v ...interface{}) {
	l.log(klog.LevelInfo, format, v...)
}

func (l *standardLogger) Noticef(format string, v ...interface{}) {
	l.log(klog.LevelNotice, format, v...)
}

func (l *standardLogger) Warnf(format string, v ...interface{}) {
	l.log(klog.LevelWarn, format, v...)
}

func (l *standardLogger) Errorf(format string, v ...interface{}) {
	l.log(klog.LevelError, format, v...)
}

func (l *standardLogger) Fatalf(format string, v ...interface{}) {
	l.log(klog.LevelFatal, format, v...)
}

// 实现 klog.CtxLogger 接口
func (l *standardLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelTrace, format, v...)
}

func (l *standardLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelDebug, format, v...)
}

func (l *standardLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelInfo, format, v...)
}

func (l *standardLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelNotice, format, v...)
}

func (l *standardLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelWarn, format, v...)
}

func (l *standardLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelError, format, v...)
}

func (l *standardLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.log(klog.LevelFatal, format, v...)
}

// 实现 klog.Control 接口
func (l *standardLogger) SetLevel(level klog.Level) {
	l.level = level
}

func (l *standardLogger) SetOutput(w io.Writer) {
	l.writer = w
}

// WithField 添加字段到日志上下文
func (l *standardLogger) WithField(key string, value interface{}) Logger {
	newLogger := &standardLogger{
		level:      l.level,
		writer:     l.writer,
		format:     l.format,
		fields:     make(map[string]interface{}, len(l.fields)+1),
		file:       l.file, // 共享文件句柄
		config:     l.config,
		rotateTime: l.rotateTime,
	}

	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// 添加新字段
	newLogger.fields[key] = value
	return newLogger
}

// Close 关闭日志器，释放资源
func (l *standardLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// shouldRotate 检查是否需要轮转
func (l *standardLogger) shouldRotate() bool {
	if l.file == nil || l.config.MaxSize <= 0 {
		return false
	}

	// 检查文件大小
	info, err := l.file.Stat()
	if err != nil {
		return false
	}

	return info.Size() >= l.config.MaxSize
}

// rotate 执行日志轮转
func (l *standardLogger) rotate() error {
	if l.file == nil {
		return nil
	}

	// 关闭当前文件
	if err := l.file.Close(); err != nil {
		return err
	}

	// 重命名当前日志文件
	oldPath := l.config.OutputPath
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	newPath := oldPath + "." + timestamp

	if err := os.Rename(oldPath, newPath); err != nil {
		// 如果重命名失败，重新打开原文件
		file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		l.file = file
		l.writer = file
		return err
	}

	// 创建新的日志文件
	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.file = file
	l.writer = file
	l.rotateTime = time.Now().Add(24 * time.Hour)

	// 清理旧文件
	go l.cleanupOldFiles()

	return nil
}

// cleanupOldFiles 清理旧的日志文件
func (l *standardLogger) cleanupOldFiles() {
	if l.config.MaxBackups <= 0 && l.config.MaxAge <= 0 {
		return
	}

	dir := filepath.Dir(l.config.OutputPath)
	baseName := filepath.Base(l.config.OutputPath)

	files, err := filepath.Glob(filepath.Join(dir, baseName+".*"))
	if err != nil {
		return
	}

	// 按修改时间排序
	sort.Slice(files, func(i, j int) bool {
		info1, err1 := os.Stat(files[i])
		info2, err2 := os.Stat(files[j])
		if err1 != nil || err2 != nil {
			return false
		}
		return info1.ModTime().After(info2.ModTime())
	})

	// 按数量删除
	if l.config.MaxBackups > 0 && len(files) > l.config.MaxBackups {
		for i := l.config.MaxBackups; i < len(files); i++ {
			os.Remove(files[i])
		}
		files = files[:l.config.MaxBackups]
	}

	// 按时间删除
	if l.config.MaxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -l.config.MaxAge)
		for _, file := range files {
			if info, err := os.Stat(file); err == nil {
				if info.ModTime().Before(cutoff) {
					os.Remove(file)
				}
			}
		}
	}
}
