// Package validator 提供了一个通用的数据验证机制
// 支持结构体字段验证、自定义验证规则和错误消息
package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator 验证器接口
type Validator interface {
	// Validate 验证数据
	Validate(data interface{}) (bool, []ValidationError)
	// AddRule 添加验证规则
	AddRule(field string, rule Rule)
}

// Rule 验证规则接口
type Rule interface {
	// Validate 验证字段值
	Validate(value interface{}) (bool, string)
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误消息
}

// defaultValidator 默认验证器实现
type defaultValidator struct {
	rules map[string][]Rule // 字段名 -> 规则列表
}

// NewValidator 创建一个新的验证器
func NewValidator() Validator {
	return &defaultValidator{
		rules: make(map[string][]Rule),
	}
}

// AddRule 添加验证规则
func (v *defaultValidator) AddRule(field string, rule Rule) {
	if _, ok := v.rules[field]; !ok {
		v.rules[field] = make([]Rule, 0)
	}
	v.rules[field] = append(v.rules[field], rule)
}

// Validate 验证数据
func (v *defaultValidator) Validate(data interface{}) (bool, []ValidationError) {
	val := reflect.ValueOf(data)

	// 如果是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 只支持结构体验证
	if val.Kind() != reflect.Struct {
		return false, []ValidationError{{
			Field:   "",
			Message: "只支持结构体验证",
		}}
	}

	typ := val.Type()
	errors := make([]ValidationError, 0)

	// 遍历所有字段
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := val.Field(i).Interface()

		// 检查是否有该字段的规则
		if rules, ok := v.rules[fieldName]; ok {
			for _, rule := range rules {
				// 验证字段
				if valid, message := rule.Validate(fieldValue); !valid {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: message,
					})
				}
			}
		}
	}

	return len(errors) == 0, errors
}

// Required 必填规则
type Required struct {
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (r *Required) Validate(value interface{}) (bool, string) {
	if value == nil {
		return false, getMessage(r.Message, "字段不能为空")
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return val.String() != "", getMessage(r.Message, "字段不能为空")
	case reflect.Slice, reflect.Map, reflect.Array:
		return val.Len() > 0, getMessage(r.Message, "字段不能为空")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0, getMessage(r.Message, "字段不能为零")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0, getMessage(r.Message, "字段不能为零")
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0, getMessage(r.Message, "字段不能为零")
	case reflect.Bool:
		return val.Bool(), getMessage(r.Message, "字段必须为true")
	default:
		return true, ""
	}
}

// MinLength 最小长度规则
type MinLength struct {
	Length  int    // 最小长度
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (m *MinLength) Validate(value interface{}) (bool, string) {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		str := val.String()
		if utf8.RuneCountInString(str) < m.Length {
			return false, getMessage(m.Message, fmt.Sprintf("字段长度不能小于%d", m.Length))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() < m.Length {
			return false, getMessage(m.Message, fmt.Sprintf("字段长度不能小于%d", m.Length))
		}
	default:
		return false, getMessage(m.Message, "不支持的类型")
	}

	return true, ""
}

// MaxLength 最大长度规则
type MaxLength struct {
	Length  int    // 最大长度
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (m *MaxLength) Validate(value interface{}) (bool, string) {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		str := val.String()
		if utf8.RuneCountInString(str) > m.Length {
			return false, getMessage(m.Message, fmt.Sprintf("字段长度不能大于%d", m.Length))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() > m.Length {
			return false, getMessage(m.Message, fmt.Sprintf("字段长度不能大于%d", m.Length))
		}
	default:
		return false, getMessage(m.Message, "不支持的类型")
	}

	return true, ""
}

// Pattern 正则表达式规则
type Pattern struct {
	Regex   string // 正则表达式
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (p *Pattern) Validate(value interface{}) (bool, string) {
	val := reflect.ValueOf(value)

	if val.Kind() != reflect.String {
		return false, getMessage(p.Message, "不支持的类型")
	}

	str := val.String()
	matched, err := regexp.MatchString(p.Regex, str)
	if err != nil {
		return false, "正则表达式错误: " + err.Error()
	}

	if !matched {
		return false, getMessage(p.Message, "字段格式不正确")
	}

	return true, ""
}

// Email 邮箱规则
type Email struct {
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (e *Email) Validate(value interface{}) (bool, string) {
	val := reflect.ValueOf(value)

	if val.Kind() != reflect.String {
		return false, getMessage(e.Message, "不支持的类型")
	}

	str := val.String()
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, str)

	if !matched {
		return false, getMessage(e.Message, "邮箱格式不正确")
	}

	return true, ""
}

// Phone 手机号规则
type Phone struct {
	Message string // 自定义错误消息
}

// Validate 验证字段值
func (p *Phone) Validate(value interface{}) (bool, string) {
	val := reflect.ValueOf(value)

	if val.Kind() != reflect.String {
		return false, getMessage(p.Message, "不支持的类型")
	}

	str := val.String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, str)

	if !matched {
		return false, getMessage(p.Message, "手机号格式不正确")
	}

	return true, ""
}

// CustomRule 自定义规则
type CustomRule struct {
	Func    func(interface{}) (bool, string) // 自定义验证函数
	Message string                           // 自定义错误消息
}

// Validate 验证字段值
func (c *CustomRule) Validate(value interface{}) (bool, string) {
	valid, msg := c.Func(value)
	if !valid {
		return false, getMessage(c.Message, msg)
	}
	return true, ""
}

// getMessage 获取错误消息
func getMessage(custom, defaultMsg string) string {
	if custom != "" {
		return custom
	}
	return defaultMsg
}

// FormatErrors 格式化验证错误
func FormatErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return sb.String()
}
