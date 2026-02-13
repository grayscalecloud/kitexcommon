// Package hderrors  提供了一个通用的错误处理机制
// 支持错误包装、错误堆栈和错误类型检查
package hderrors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// 标准错误类型
var (
	// New 标准库错误
	New    = errors.New
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

type EnumsType interface {
	ToInt() int64
}

type DefaultEnumsType struct {
	Code int64
}

func NewDefaultEnumsType(code int64) *DefaultEnumsType {
	return &DefaultEnumsType{
		Code: code,
	}
}

func (e *DefaultEnumsType) ToInt() int64 {
	return e.Code
}

// BusinessError 自定义错误结构
// 包含错误码、错误消息和堆栈信息
type BusinessError struct {
	Code    EnumsType   // 错误码
	Message string      // 错误消息
	Cause   error       // 原始错误
	Stack   []StackInfo // 堆栈信息
	extra   map[string]string
}

// StackInfo 堆栈信息
type StackInfo struct {
	File     string // 文件名
	Line     int    // 行号
	Function string // 函数名
}

// NewError 创建一个新的错误
//
// 参数:
//   - code: 错误码
//   - message: 错误消息
//
// 返回:
//   - 自定义错误
func NewError(code EnumsType, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Stack:   captureStack(2),
		extra:   make(map[string]string),
	}
}

// Wrap 包装一个已有的错误
//
// 参数:
//   - err: 原始错误
//   - code: 错误码
//   - message: 错误消息
//
// 返回:
//   - 包装后的错误
func Wrap(err error, code EnumsType, message string) *BusinessError {
	if err == nil {
		return nil
	}

	extra := make(map[string]string)
	// 如果原始错误是 BusinessError，保留其 extra 信息，并用其 message 替换传入的 message
	var be *BusinessError
	if errors.As(err, &be) {
		if be.extra != nil {
			for k, v := range be.extra {
				extra[k] = v
			}
		}
	}

	msg := message
	if be != nil {
		msg = be.Message
	}

	return &BusinessError{
		Code:    code,
		Message: msg,
		Cause:   err,
		Stack:   captureStack(2),
		extra:   extra,
	}
}

// WrapWithMessage 用新消息包装错误
//
// 参数:
//   - err: 原始错误
//   - message: 错误消息
//
// 返回:
//   - 包装后的错误
func WrapWithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	// 如果是自定义错误，保留错误码和额外信息
	var e *BusinessError
	if errors.As(err, &e) {
		extra := make(map[string]string)
		if e.extra != nil {
			for k, v := range e.extra {
				extra[k] = v
			}
		}
		return &BusinessError{
			Code:    e.Code,
			Message: message,
			Cause:   e.Cause,
			Stack:   e.Stack, // 保留原始堆栈
			extra:   extra,   // 复制额外的错误信息
		}
	}

	return &BusinessError{
		Code:    NewDefaultEnumsType(-1), // 默认错误码
		Message: message,
		Cause:   err,
		Stack:   captureStack(2),
		extra:   make(map[string]string),
	}
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e == nil {
		return "<nil>"
	}

	codeStr := "0"
	if e.Code != nil {
		codeStr = fmt.Sprintf("%d", e.Code.ToInt())
	}

	if e.Cause != nil {
		return fmt.Sprintf("[BusinessError:%s] %s: %v", codeStr, e.Message, e.Cause)
	}
	return fmt.Sprintf("[BusinessError:%s] %s", codeStr, e.Message)
}

// Unwrap 获取原始错误
func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// GetCode 获取错误码
func (e *BusinessError) GetCode() int64 {
	return e.Code.ToInt()
}

// BizStatusCode 获取业务错误状态码（用于RPC返回）
func (e *BusinessError) BizStatusCode() int32 {
	return int32(e.Code.ToInt())
}

// BizMessage 获取业务错误消息（用于RPC返回）
func (e *BusinessError) BizMessage() string {
	return e.Message
}

// BizExtra 获取业务错误的额外信息（用于RPC返回）
func (e *BusinessError) BizExtra() map[string]string {
	if e.extra == nil {
		return make(map[string]string)
	}
	return e.extra
}

// GetMessage 获取错误消息
func (e *BusinessError) GetMessage() string {
	return e.Message
}

// GetStack 获取堆栈信息
func (e *BusinessError) GetStack() []StackInfo {
	return e.Stack
}

// FormatStack 格式化堆栈信息
func (e *BusinessError) FormatStack() string {
	if len(e.Stack) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Stack Trace:\n")
	for _, frame := range e.Stack {
		sb.WriteString(fmt.Sprintf("  %s:%d %s\n", frame.File, frame.Line, frame.Function))
	}
	return sb.String()
}

// captureStack 捕获当前堆栈信息
func captureStack(skip int) []StackInfo {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])

	stack := make([]StackInfo, 0, n)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		stack = append(stack, StackInfo{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		})
	}

	return stack
}

// IsBusinessError 检查错误是否为业务错误
func IsBusinessError(err error) bool {
	var businessError *BusinessError
	ok := errors.As(err, &businessError)
	return ok
}

// SetExtra 设置额外的错误信息
//
// 参数:
//   - key: 键
//   - value: 值
func (e *BusinessError) SetExtra(key, value string) {
	if e == nil {
		return
	}
	if e.extra == nil {
		e.extra = make(map[string]string)
	}
	e.extra[key] = value
}

// SetExtras 批量设置额外的错误信息
//
// 参数:
//   - extras: 额外的错误信息映射
func (e *BusinessError) SetExtras(extras map[string]string) {
	if e == nil {
		return
	}
	if e.extra == nil {
		e.extra = make(map[string]string)
	}
	for k, v := range extras {
		e.extra[k] = v
	}
}

// GetExtra 获取额外的错误信息
//
// 参数:
//   - key: 键
//
// 返回:
//   - 值，如果不存在返回空字符串
func (e *BusinessError) GetExtra(key string) string {
	if e == nil || e.extra == nil {
		return ""
	}
	return e.extra[key]
}

// FormatError 格式化完整的错误信息，包括堆栈
//
// 返回:
//   - 包含堆栈信息的完整错误描述
func (e *BusinessError) FormatError() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder

	// 错误基本信息
	sb.WriteString(e.Error())

	// 额外信息
	if len(e.extra) > 0 {
		sb.WriteString("\nExtra Info:")
		for k, v := range e.extra {
			sb.WriteString(fmt.Sprintf("\n  %s: %s", k, v))
		}
	}

	// 堆栈信息
	stackStr := e.FormatStack()
	if stackStr != "" {
		sb.WriteString("\n")
		sb.WriteString(stackStr)
	}

	return sb.String()
}

// WithExtra 链式设置额外的错误信息
//
// 参数:
//   - key: 键
//   - value: 值
//
// 返回:
//   - 错误本身，支持链式调用
func (e *BusinessError) WithExtra(key, value string) *BusinessError {
	e.SetExtra(key, value)
	return e
}
