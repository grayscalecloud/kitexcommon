package enum

// WithdrawStatus 提现状态枚举
// 用于表示用户提现申请的处理状态
type WithdrawStatus int

const (
	// WithdrawStatusUnknown 未知状态 - 初始状态或异常状态
	WithdrawStatusUnknown WithdrawStatus = iota
	// WithdrawStatusPending 待处理 - 提现申请已提交，等待处理
	WithdrawStatusPending
	// WithdrawStatusProcessing 处理中 - 提现申请正在处理中
	WithdrawStatusProcessing
	// WithdrawStatusSuccess 提现成功 - 提现申请处理完成且成功
	WithdrawStatusSuccess
	// WithdrawStatusFailed 提现失败 - 提现申请处理失败
	WithdrawStatusFailed
	// WithdrawStatusCancelled 已取消 - 用户主动取消提现申请
	WithdrawStatusCancelled
	// WithdrawStatusFrozen 冻结中 - 提现申请因风控等原因被冻结
	WithdrawStatusFrozen
	// WithdrawStatusPendingConfirm 待确认 - 提现申请需要进一步确认
	WithdrawStatusPendingConfirm
	// WithdrawStatusRejected 已驳回 - 提现申请被驳回
	WithdrawStatusRejected
)

// String 返回提现状态的字符串表示
func (ws WithdrawStatus) String() string {
	switch ws {
	case WithdrawStatusPending:
		return "pending"
	case WithdrawStatusProcessing:
		return "processing"
	case WithdrawStatusSuccess:
		return "success"
	case WithdrawStatusFailed:
		return "failed"
	case WithdrawStatusCancelled:
		return "cancelled"
	case WithdrawStatusFrozen:
		return "frozen"
	case WithdrawStatusPendingConfirm:
		return "pending_confirm"
	case WithdrawStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

// CNString 返回提现状态的中文字符串表示
func (ws WithdrawStatus) CNString() string {
	switch ws {
	case WithdrawStatusPending:
		return "待处理"
	case WithdrawStatusProcessing:
		return "处理中"
	case WithdrawStatusSuccess:
		return "提现成功"
	case WithdrawStatusFailed:
		return "提现失败"
	case WithdrawStatusCancelled:
		return "已取消"
	case WithdrawStatusFrozen:
		return "冻结中"
	case WithdrawStatusPendingConfirm:
		return "待确认"
	case WithdrawStatusRejected:
		return "已驳回"
	default:
		return "未知"
	}
}

// IsValid 检查提现状态是否有效
func (ws WithdrawStatus) IsValid() bool {
	return ws >= WithdrawStatusUnknown && ws <= WithdrawStatusRejected
}

// WithdrawStatusFromString 从字符串创建提现状态
func WithdrawStatusFromString(s string) WithdrawStatus {
	switch s {
	case "pending":
		return WithdrawStatusPending
	case "processing":
		return WithdrawStatusProcessing
	case "success":
		return WithdrawStatusSuccess
	case "failed":
		return WithdrawStatusFailed
	case "cancelled":
		return WithdrawStatusCancelled
	case "frozen":
		return WithdrawStatusFrozen
	case "pending_confirm":
		return WithdrawStatusPendingConfirm
	case "rejected":
		return WithdrawStatusRejected
	default:
		return WithdrawStatusUnknown
	}
}

// WithdrawStatusFromCNString 从中文字符串创建提现状态
func WithdrawStatusFromCNString(s string) WithdrawStatus {
	switch s {
	case "待处理":
		return WithdrawStatusPending
	case "处理中":
		return WithdrawStatusProcessing
	case "提现成功":
		return WithdrawStatusSuccess
	case "提现失败":
		return WithdrawStatusFailed
	case "已取消":
		return WithdrawStatusCancelled
	case "冻结中":
		return WithdrawStatusFrozen
	case "待确认":
		return WithdrawStatusPendingConfirm
	case "已驳回":
		return WithdrawStatusRejected
	default:
		return WithdrawStatusUnknown
	}
}
