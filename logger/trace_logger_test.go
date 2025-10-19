package logger

import (
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
)

func TestNewTraceLogger(t *testing.T) {
	prefix := "test-prefix"
	logger := NewTraceLogger(prefix)

	if logger == nil {
		t.Fatal("NewTraceLogger returned nil")
	}

	if logger.prefix != prefix {
		t.Errorf("Expected prefix %s, got %s", prefix, logger.prefix)
	}

	if logger.fields == nil {
		t.Error("Fields map not initialized")
	}

	if len(logger.fields) != 0 {
		t.Errorf("Expected empty fields map, got %d fields", len(logger.fields))
	}
}

func TestTraceLogger_InterfaceImplementation(t *testing.T) {
	logger := NewTraceLogger("test")

	// 验证实现了 Logger 接口
	var _ Logger = logger
}

func TestTraceLogger_WithField(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试添加单个字段
	fieldLogger := logger.WithField("key1", "value1")

	if fieldLogger == nil {
		t.Error("WithField returned nil")
	}

	// 验证字段被正确添加
	traceLogger, ok := fieldLogger.(*TraceLogger)
	if !ok {
		t.Error("WithField did not return TraceLogger")
	}

	if len(traceLogger.fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(traceLogger.fields))
	}

	if traceLogger.fields["key1"] != "value1" {
		t.Errorf("Expected field value 'value1', got %v", traceLogger.fields["key1"])
	}
}

func TestTraceLogger_WithField_Chaining(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试链式添加字段
	fieldLogger := logger.WithField("key1", "value1").WithField("key2", 123).WithField("key3", true)

	traceLogger, ok := fieldLogger.(*TraceLogger)
	if !ok {
		t.Error("WithField chaining did not return TraceLogger")
	}

	if len(traceLogger.fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(traceLogger.fields))
	}

	expectedFields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	for k, v := range expectedFields {
		if traceLogger.fields[k] != v {
			t.Errorf("Expected field %s=%v, got %v", k, v, traceLogger.fields[k])
		}
	}
}

func TestTraceLogger_WithField_Isolation(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试字段隔离
	fieldLogger1 := logger.WithField("key1", "value1")
	fieldLogger2 := logger.WithField("key2", "value2")

	traceLogger1, ok1 := fieldLogger1.(*TraceLogger)
	traceLogger2, ok2 := fieldLogger2.(*TraceLogger)

	if !ok1 || !ok2 {
		t.Error("WithField did not return TraceLogger")
	}

	// 验证字段隔离
	if len(traceLogger1.fields) != 1 {
		t.Errorf("Expected 1 field in logger1, got %d", len(traceLogger1.fields))
	}

	if len(traceLogger2.fields) != 1 {
		t.Errorf("Expected 1 field in logger2, got %d", len(traceLogger2.fields))
	}

	if traceLogger1.fields["key1"] != "value1" {
		t.Error("Field isolation failed for logger1")
	}

	if traceLogger2.fields["key2"] != "value2" {
		t.Error("Field isolation failed for logger2")
	}
}

func TestTraceLogger_Close(t *testing.T) {
	logger := NewTraceLogger("test")

	// Close 方法应该不返回错误
	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestTraceLogger_StandardMethods(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试标准方法不会 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Standard method panicked: %v", r)
		}
	}()

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
}

func TestTraceLogger_FormatMethods(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试格式化方法不会 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Format method panicked: %v", r)
		}
	}()

	logger.Debugf("Debug formatted: %s", "test")
	logger.Infof("Info formatted: %d", 123)
	logger.Warnf("Warn formatted: %v", true)
	logger.Errorf("Error formatted: %s", "error")
}

func TestTraceLogger_WithFieldAndStandardMethods(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试带字段的标准方法
	fieldLogger := logger.WithField("user_id", "12345").WithField("request_id", "req-001")

	// 这些调用不应该 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Field logger method panicked: %v", r)
		}
	}()

	fieldLogger.Debug("Debug message with fields")
	fieldLogger.Info("Info message with fields")
	fieldLogger.Warn("Warn message with fields")
	fieldLogger.Error("Error message with fields")
}

func TestGetCallerInfo(t *testing.T) {
	// 测试 getCallerInfo 函数
	caller := getCallerInfo(1, "")
	if caller == "" {
		t.Error("getCallerInfo returned empty string")
	}

	// 测试带项目根目录的情况
	callerWithRoot := getCallerInfo(1, "/some/project/root")
	if callerWithRoot == "" {
		t.Error("getCallerInfo with root returned empty string")
	}
}

func TestProjectRootInitialization(t *testing.T) {
	// 验证项目根目录被正确初始化
	if projectRoot == "" {
		t.Error("Project root not initialized")
	}
}

func TestTraceLogger_SetLevel(t *testing.T) {
	logger := NewTraceLogger("test")

	// 测试 SetLevel 方法
	logger.SetLevel(klog.LevelDebug)
	logger.SetLevel(klog.LevelInfo)
	logger.SetLevel(klog.LevelWarn)
	logger.SetLevel(klog.LevelError)
}

func BenchmarkTraceLogger_Info(b *testing.B) {
	logger := NewTraceLogger("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("Benchmark message %d", i)
	}
}

func BenchmarkTraceLogger_InfoWithFields(b *testing.B) {
	logger := NewTraceLogger("benchmark")
	fieldLogger := logger.WithField("key1", "value1").WithField("key2", 123)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fieldLogger.Infof("Benchmark message %d", i)
	}
}

func BenchmarkTraceLogger_WithField(b *testing.B) {
	logger := NewTraceLogger("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithField("key", i)
	}
}
