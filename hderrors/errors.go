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

	return &BusinessError{
		Code:    code,
		Message: message,
		Cause:   err,
		Stack:   captureStack(2),
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

	// 如果是自定义错误，保留错误码
	var e *BusinessError
	if errors.As(err, &e) {
		return &BusinessError{
			Code:    e.Code,
			Message: message,
			Cause:   e.Cause,
			Stack:   captureStack(2),
		}
	}

	return &BusinessError{
		Code:    NewDefaultEnumsType(-1), // 默认错误码
		Message: message,
		Cause:   err,
		Stack:   captureStack(2),
	}
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 获取原始错误
func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// GetCode 获取错误码
func (e *BusinessError) GetCode() int64 {
	return e.Code.ToInt()
}

func (e *BusinessError) BizStatusCode() int32 {
	return int32(e.Code.ToInt())
}
func (e *BusinessError) BizMessage() string {
	return e.Message
}
func (e *BusinessError) BizExtra() map[string]string {
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
