package logger

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/grayscalecloud/kitexcommon/tools"

	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	// 使用编译时的文件路径来确定项目根目录
	//_, file, _, ok := runtime.Caller(0)
	//if !ok {
	//	projectRoot = ""
	//	return
	//}
	//// 获取base包所在目录的父目录作为项目根目录
	//projectRoot = filepath.Dir(filepath.Dir(file))
	//fmt.Println("根目录：", projectRoot)
}

type TraceLogger struct {
	*kitexlogrus.Logger
	prefix string
}

func NewTraceLogger(prefix string) *TraceLogger {
	return &TraceLogger{
		Logger: kitexlogrus.NewLogger(),
		prefix: prefix,
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
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("level", level),
			attribute.String("message", msg),
			attribute.String("caller", caller),
			attribute.String("tenantId", tenantId),
		))
	}
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
