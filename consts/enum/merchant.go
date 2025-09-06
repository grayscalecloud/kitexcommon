package enum

// MerchantStatus 商户状态枚举
// 用于表示商户的各种状态
type MerchantStatus int

const (
	// MerchantStatusUnknown 未知状态
	MerchantStatusUnknown MerchantStatus = iota
	// MerchantStatusActive 激活状态 - 商户正常营业中
	MerchantStatusActive
	// MerchantStatusInactive 非激活状态 - 商户暂停营业
	MerchantStatusInactive
	// MerchantStatusPending 待审核 - 商户信息提交审核中
	MerchantStatusPending
	// MerchantStatusRejected 审核拒绝 - 商户审核被拒绝
	MerchantStatusRejected
	// MerchantStatusFrozen 冻结状态 - 商户因违规等原因被冻结
	MerchantStatusFrozen
	// MerchantStatusClosed 注销状态 - 商户已注销
	MerchantStatusClosed
)

// String 返回商户状态的字符串表示
func (ms MerchantStatus) String() string {
	switch ms {
	case MerchantStatusActive:
		return "active"
	case MerchantStatusInactive:
		return "inactive"
	case MerchantStatusPending:
		return "pending"
	case MerchantStatusRejected:
		return "rejected"
	case MerchantStatusFrozen:
		return "frozen"
	case MerchantStatusClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// CNString 返回商户状态的中文字符串表示
func (ms MerchantStatus) CNString() string {
	switch ms {
	case MerchantStatusActive:
		return "激活"
	case MerchantStatusInactive:
		return "非激活"
	case MerchantStatusPending:
		return "待审核"
	case MerchantStatusRejected:
		return "审核拒绝"
	case MerchantStatusFrozen:
		return "冻结"
	case MerchantStatusClosed:
		return "注销"
	default:
		return "未知"
	}
}

// IsValid 检查商户状态是否有效
func (ms MerchantStatus) IsValid() bool {
	return ms == MerchantStatusUnknown ||
		ms == MerchantStatusActive ||
		ms == MerchantStatusInactive ||
		ms == MerchantStatusPending ||
		ms == MerchantStatusRejected ||
		ms == MerchantStatusFrozen ||
		ms == MerchantStatusClosed
}

// MerchantStatusFromString 从字符串创建商户状态
func MerchantStatusFromString(s string) MerchantStatus {
	switch s {
	case "active":
		return MerchantStatusActive
	case "inactive":
		return MerchantStatusInactive
	case "pending":
		return MerchantStatusPending
	case "rejected":
		return MerchantStatusRejected
	case "frozen":
		return MerchantStatusFrozen
	case "closed":
		return MerchantStatusClosed
	default:
		return MerchantStatusUnknown
	}
}

// MerchantStatusFromCNString 从中文字符串创建商户状态
func MerchantStatusFromCNString(s string) MerchantStatus {
	switch s {
	case "激活":
		return MerchantStatusActive
	case "非激活":
		return MerchantStatusInactive
	case "待审核":
		return MerchantStatusPending
	case "审核拒绝":
		return MerchantStatusRejected
	case "冻结":
		return MerchantStatusFrozen
	case "注销":
		return MerchantStatusClosed
	default:
		return MerchantStatusUnknown
	}
}
