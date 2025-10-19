package logger

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"runtime"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/grayscalecloud/kitexcommon/tools"

	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var projectRoot string

func init() {
	// 使用编译时的文件路径来确定项目根目录
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		projectRoot = ""
		return
	}
	// 获取base包所在目录的父目录作为项目根目录
	projectRoot = filepath.Dir(filepath.Dir(file))
}

// GetProjectRoot 返回项目根目录
func GetProjectRoot() string {
	return projectRoot
}

type TraceLogger struct {
	*kitexlogrus.Logger
	prefix string
	fields map[string]interface{} // 支持字段存储
}

func NewTraceLogger(prefix string) *TraceLogger {
	return &TraceLogger{
		Logger: kitexlogrus.NewLogger(),
		prefix: prefix,
		fields: make(map[string]interface{}),
	}
}

func getCallerInfo(skip int, projectRoot string) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	if projectRoot != "" {
		// 将绝对路径转换为相对于项目根目录的路径
		if rel, err := filepath.Rel(projectRoot, file); err == nil {
			file = rel
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func (l *TraceLogger) logWithTrace(ctx context.Context, level, msg string) {
	caller := getCallerInfo(4, l.prefix)
	span := trace.SpanFromContext(ctx)
	tenantId := tools.GetTenant(ctx)

	if span.IsRecording() {
		// 添加错误处理，确保属性添加不会失败
		defer func() {
			if r := recover(); r != nil {
				// 如果添加属性失败，至少记录到日志
				l.Logger.CtxErrorf(ctx, "Failed to add trace attributes: %v", r)
			}
		}()

		// 准备属性列表
		attrs := []attribute.KeyValue{
			attribute.String("level", level),
			attribute.String("message", msg),
			attribute.String("caller", caller),
			attribute.String("tenantId", tenantId),
		}

		// 添加自定义字段
		for k, v := range l.fields {
			attrs = append(attrs, attribute.String("field."+k, fmt.Sprintf("%v", v)))
		}

		span.AddEvent("log", trace.WithAttributes(attrs...))
	}
}

// 实现 klog.Logger 接口
func (l *TraceLogger) Trace(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "TRACE", fmt.Sprint(v...))
	l.Logger.Trace(v...)
}

func (l *TraceLogger) Debug(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "DEBUG", fmt.Sprint(v...))
	l.Logger.Debug(v...)
}

func (l *TraceLogger) Info(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "INFO", fmt.Sprint(v...))
	l.Logger.Info(v...)
}

func (l *TraceLogger) Notice(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "NOTICE", fmt.Sprint(v...))
	l.Logger.Notice(v...)
}

func (l *TraceLogger) Warn(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "WARN", fmt.Sprint(v...))
	l.Logger.Warn(v...)
}

func (l *TraceLogger) Error(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "ERROR", fmt.Sprint(v...))
	l.Logger.Error(v...)
}

func (l *TraceLogger) Fatal(v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "FATAL", fmt.Sprint(v...))
	l.Logger.Fatal(v...)
}

// 实现 klog.FormatLogger 接口
func (l *TraceLogger) Tracef(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "TRACE", fmt.Sprintf(format, v...))
	l.Logger.Tracef(format, v...)
}

func (l *TraceLogger) Debugf(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "DEBUG", fmt.Sprintf(format, v...))
	l.Logger.Debugf(format, v...)
}

func (l *TraceLogger) Infof(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "INFO", fmt.Sprintf(format, v...))
	l.Logger.Infof(format, v...)
}

func (l *TraceLogger) Noticef(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "NOTICE", fmt.Sprintf(format, v...))
	l.Logger.Noticef(format, v...)
}

func (l *TraceLogger) Warnf(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "WARN", fmt.Sprintf(format, v...))
	l.Logger.Warnf(format, v...)
}

func (l *TraceLogger) Errorf(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "ERROR", fmt.Sprintf(format, v...))
	l.Logger.Errorf(format, v...)
}

func (l *TraceLogger) Fatalf(format string, v ...interface{}) {
	ctx := context.Background()
	l.logWithTrace(ctx, "FATAL", fmt.Sprintf(format, v...))
	l.Logger.Fatalf(format, v...)
}

// 实现 klog.CtxLogger 接口
func (l *TraceLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "TRACE", fmt.Sprintf(format, v...))
	l.Logger.CtxTracef(ctx, format, v...)
}

func (l *TraceLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "DEBUG", fmt.Sprintf(format, v...))
	l.Logger.CtxDebugf(ctx, format, v...)
}

func (l *TraceLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "INFO", fmt.Sprintf(format, v...))
	l.Logger.CtxInfof(ctx, format, v...)
}

func (l *TraceLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "NOTICE", fmt.Sprintf(format, v...))
	l.Logger.CtxNoticef(ctx, format, v...)
}

func (l *TraceLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "WARN", fmt.Sprintf(format, v...))
	l.Logger.CtxWarnf(ctx, format, v...)
}

func (l *TraceLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "ERROR", fmt.Sprintf(format, v...))
	l.Logger.CtxErrorf(ctx, format, v...)
}

func (l *TraceLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "FATAL", fmt.Sprintf(format, v...))
	l.Logger.CtxFatalf(ctx, format, v...)
}

// 实现 klog.Control 接口
func (l *TraceLogger) SetLevel(level klog.Level) {
	l.Logger.SetLevel(level)
}

func (l *TraceLogger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

// WithField 添加字段到日志上下文
func (l *TraceLogger) WithField(key string, value interface{}) Logger {
	// 创建一个新的 TraceLogger 实例，保持相同的配置
	newLogger := &TraceLogger{
		Logger: l.Logger,
		prefix: l.prefix,
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

// Close 关闭日志器，释放资源
func (l *TraceLogger) Close() error {
	// TraceLogger 使用 kitexlogrus.Logger，通常不需要手动关闭
	// 但为了接口一致性，提供这个方法
	return nil
}
