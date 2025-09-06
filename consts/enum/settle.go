package enum

// SettleStatus 结算状态枚举
// 用于表示交易结算的处理状态
type SettleStatus int

const (
	// SettleStatusUnknown 未知状态
	SettleStatusUnknown SettleStatus = iota
	// SettleStatusPending 待结算 - 交易已完成，等待结算
	SettleStatusPending
	// SettleStatusProcessing 结算中 - 正在进行结算处理
	SettleStatusProcessing
	// SettleStatusSuccess 结算成功 - 结算已完成
	SettleStatusSuccess
	// SettleStatusFailed 结算失败 - 结算处理失败
	SettleStatusFailed
	// SettleStatusCancelled 已取消 - 结算被取消
	SettleStatusCancelled
	// SettleStatusVerifying 验证中 - 结算信息正在验证中
	SettleStatusVerifying
	// SettleStatusPartial 部分结算 - 只完成了部分结算
	SettleStatusPartial
)

// String 返回结算状态的字符串表示
func (ss SettleStatus) String() string {
	switch ss {
	case SettleStatusPending:
		return "pending"
	case SettleStatusProcessing:
		return "processing"
	case SettleStatusSuccess:
		return "success"
	case SettleStatusFailed:
		return "failed"
	case SettleStatusCancelled:
		return "cancelled"
	case SettleStatusVerifying:
		return "verifying"
	case SettleStatusPartial:
		return "partial"
	default:
		return "unknown"
	}
}

// CNString 返回结算状态的中文字符串表示
func (ss SettleStatus) CNString() string {
	switch ss {
	case SettleStatusPending:
		return "待结算"
	case SettleStatusProcessing:
		return "结算中"
	case SettleStatusSuccess:
		return "结算成功"
	case SettleStatusFailed:
		return "结算失败"
	case SettleStatusCancelled:
		return "已取消"
	case SettleStatusVerifying:
		return "验证中"
	case SettleStatusPartial:
		return "部分结算"
	default:
		return "未知"
	}
}

// IsValid 检查结算状态是否有效
func (ss SettleStatus) IsValid() bool {
	return ss == SettleStatusUnknown ||
		ss == SettleStatusPending ||
		ss == SettleStatusProcessing ||
		ss == SettleStatusSuccess ||
		ss == SettleStatusFailed ||
		ss == SettleStatusCancelled ||
		ss == SettleStatusVerifying ||
		ss == SettleStatusPartial
}

// SettleStatusFromString 从字符串创建结算状态
func SettleStatusFromString(s string) SettleStatus {
	switch s {
	case "pending":
		return SettleStatusPending
	case "processing":
		return SettleStatusProcessing
	case "success":
		return SettleStatusSuccess
	case "failed":
		return SettleStatusFailed
	case "cancelled":
		return SettleStatusCancelled
	case "verifying":
		return SettleStatusVerifying
	case "partial":
		return SettleStatusPartial
	default:
		return SettleStatusUnknown
	}
}

// SettleStatusFromCNString 从中文字符串创建结算状态
func SettleStatusFromCNString(s string) SettleStatus {
	switch s {
	case "待结算":
		return SettleStatusPending
	case "结算中":
		return SettleStatusProcessing
	case "结算成功":
		return SettleStatusSuccess
	case "结算失败":
		return SettleStatusFailed
	case "已取消":
		return SettleStatusCancelled
	case "验证中":
		return SettleStatusVerifying
	case "部分结算":
		return SettleStatusPartial
	default:
		return SettleStatusUnknown
	}
}
