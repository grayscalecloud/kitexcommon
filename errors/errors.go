// Package errors 提供了一个通用的错误处理机制
// 支持错误包装、错误堆栈和错误类型检查
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// 标准错误类型
var (
	// 标准库错误
	New    = errors.New
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

// Error 自定义错误结构
// 包含错误码、错误消息和堆栈信息
type Error struct {
	Code    int         // 错误码
	Message string      // 错误消息
	Cause   error       // 原始错误
	Stack   []StackInfo // 堆栈信息
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
func NewError(code int, message string) *Error {
	return &Error{
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
func Wrap(err error, code int, message string) *Error {
	if err == nil {
		return nil
	}

	return &Error{
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
	if e, ok := err.(*Error); ok {
		return &Error{
			Code:    e.Code,
			Message: message,
			Cause:   e.Cause,
			Stack:   captureStack(2),
		}
	}

	return &Error{
		Code:    -1, // 默认错误码
		Message: message,
		Cause:   err,
		Stack:   captureStack(2),
	}
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 获取原始错误
func (e *Error) Unwrap() error {
	return e.Cause
}

// GetCode 获取错误码
func (e *Error) GetCode() int {
	return e.Code
}

// GetMessage 获取错误消息
func (e *Error) GetMessage() string {
	return e.Message
}

// GetStack 获取堆栈信息
func (e *Error) GetStack() []StackInfo {
	return e.Stack
}

// FormatStack 格式化堆栈信息
func (e *Error) FormatStack() string {
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

		if !more {
			break
		}
	}

	return stack
}
