package enum

// Status 状态枚举
// 用于表示各种业务状态，只包含启用和禁用两种状态
type Status int

const (
	// StatusUnknown 未知状态
	StatusUnknown Status = 99
	// StatusEnabled 启用状态 - 表示该资源或功能已启用，可以正常使用
	StatusEnabled Status = 1
	// StatusDisabled 禁用状态 - 表示该资源或功能已被禁用，无法使用
	StatusDisabled Status = 0
)

// String 返回状态的字符串表示
func (s Status) String() string {
	switch s {
	case StatusEnabled:
		return "enabled"
	case StatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// CNString 返回状态的中文字符串表示
func (s Status) CNString() string {
	switch s {
	case StatusEnabled:
		return "启用"
	case StatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

// IsValid 检查状态是否有效
func (s Status) IsValid() bool {
	return s == StatusUnknown || s == StatusEnabled || s == StatusDisabled
}

// StatusFromString 从字符串创建状态
func StatusFromString(s string) Status {
	switch s {
	case "enabled":
		return StatusEnabled
	case "disabled":
		return StatusDisabled
	default:
		return StatusUnknown
	}
}
