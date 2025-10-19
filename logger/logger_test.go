package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Level != klog.LevelInfo {
		t.Errorf("DefaultConfig().Level = %v, want %v", config.Level, klog.LevelInfo)
	}
	if config.OutputPath != "" {
		t.Errorf("DefaultConfig().OutputPath = %v, want empty string", config.OutputPath)
	}
	if config.Format != "text" {
		t.Errorf("DefaultConfig().Format = %v, want text", config.Format)
	}
}

func TestNewLogger_WithNilConfig(t *testing.T) {
	logger := NewLogger(nil)
	if logger == nil {
		t.Error("NewLogger(nil) returned nil")
	}
	defer logger.Close()
}

func TestNewLogger_WithConfig(t *testing.T) {
	config := &Config{
		Level:      klog.LevelDebug,
		OutputPath: "",
		Format:     "json",
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Error("NewLogger(config) returned nil")
	}
	defer logger.Close()
}

func TestStandardLogger_LogLevels(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelDebug,
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 测试各个日志级别
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, level := range expectedLevels {
		if !strings.Contains(contentStr, level) {
			t.Errorf("Log file does not contain %s level", level)
		}
	}
}

func TestStandardLogger_WithField(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 测试 WithField
	fieldLogger := logger.WithField("user_id", "12345").WithField("request_id", "req-001")
	fieldLogger.Info("User action")

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "user_id=12345") {
		t.Error("Log file does not contain user_id field")
	}
	if !strings.Contains(contentStr, "request_id=req-001") {
		t.Error("Log file does not contain request_id field")
	}
}

func TestStandardLogger_JSONFormat(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: tmpFile.Name(),
		Format:     "json",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 测试 JSON 格式
	logger.Info("Test message")

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, `"level":"INFO"`) {
		t.Error("JSON log does not contain level field")
	}
	if !strings.Contains(contentStr, `"message":"Test message"`) {
		t.Error("JSON log does not contain message field")
	}
	if !strings.Contains(contentStr, `"time"`) {
		t.Error("JSON log does not contain time field")
	}
}

func TestStandardLogger_JSONFormatWithFields(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: tmpFile.Name(),
		Format:     "json",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 测试 JSON 格式带字段
	fieldLogger := logger.WithField("key1", "value1").WithField("key2", 123)
	fieldLogger.Info("Test message with fields")

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, `"key1":"value1"`) {
		t.Error("JSON log does not contain key1 field")
	}
	if !strings.Contains(contentStr, `"key2":123`) {
		t.Error("JSON log does not contain key2 field")
	}
}

func TestStandardLogger_LogLevelFiltering(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelWarn, // 只记录 WARN 和 ERROR
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 记录不同级别的日志
	logger.Debug("Debug message") // 应该被过滤
	logger.Info("Info message")   // 应该被过滤
	logger.Warn("Warn message")   // 应该被记录
	logger.Error("Error message") // 应该被记录

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "Debug message") {
		t.Error("Debug message should be filtered out")
	}
	if strings.Contains(contentStr, "Info message") {
		t.Error("Info message should be filtered out")
	}
	if !strings.Contains(contentStr, "Warn message") {
		t.Error("Warn message should be recorded")
	}
	if !strings.Contains(contentStr, "Error message") {
		t.Error("Error message should be recorded")
	}
}

func TestStandardLogger_Close(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)

	// 测试 Close 方法
	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// 再次关闭应该不会出错（但可能会返回错误，这是正常的）
	logger.Close() // 忽略错误，因为文件已经关闭
}

func TestStandardLogger_CloseWithStdout(t *testing.T) {
	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: "", // 使用标准输出
		Format:     "text",
	}
	logger := NewLogger(config)

	// 使用标准输出时，Close 应该不返回错误
	if err := logger.Close(); err != nil {
		t.Errorf("Close() with stdout returned error: %v", err)
	}
}

func TestStandardLogger_ConcurrentAccess(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 并发写入测试
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				logger.Info("Concurrent message %d-%d", id, j)
			}
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证文件内容
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	// 应该有 100 行日志（10 goroutines * 10 messages each）
	// 加上最后一个空行，所以是 101 行
	if len(lines) < 100 {
		t.Errorf("Expected at least 100 log lines, got %d", len(lines)-1)
	}
}

func TestStandardLogger_FormatMethods(t *testing.T) {
	// 创建一个临时文件用于测试
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := &Config{
		Level:      klog.LevelDebug, // 设置为 Debug 级别以包含所有日志
		OutputPath: tmpFile.Name(),
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	// 测试格式化方法
	logger.Infof("Formatted message: %s", "test")
	logger.Debugf("Debug formatted: %d", 123)
	logger.Warnf("Warning formatted: %v", true)

	// 读取文件内容验证
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Formatted message: test") {
		t.Error("Formatted message not found")
	}
	if !strings.Contains(contentStr, "Debug formatted: 123") {
		t.Error("Debug formatted message not found")
	}
	if !strings.Contains(contentStr, "Warning formatted: true") {
		t.Error("Warning formatted message not found")
	}
}

func BenchmarkStandardLogger_Info(b *testing.B) {
	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: "",
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message %d", i)
	}
}

func BenchmarkStandardLogger_InfoWithFields(b *testing.B) {
	config := &Config{
		Level:      klog.LevelInfo,
		OutputPath: "",
		Format:     "text",
	}
	logger := NewLogger(config)
	defer logger.Close()

	fieldLogger := logger.WithField("key1", "value1").WithField("key2", 123)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fieldLogger.Info("Benchmark message %d", i)
	}
}
