package enum

// InvoiceStatus 开票状态枚举
// 用于表示发票申请的处理状态
type InvoiceStatus int

const (
	// InvoiceStatusUnknown 未知状态 - 初始状态或异常状态
	InvoiceStatusUnknown InvoiceStatus = iota
	// InvoiceStatusApplied 已申请 - 用户已提交开票申请
	InvoiceStatusApplied
	// InvoiceStatusProcessing 处理中 - 开票申请正在处理中
	InvoiceStatusProcessing
	// InvoiceStatusIssued 已开具 - 发票已成功开具
	InvoiceStatusIssued
	// InvoiceStatusFailed 开具失败 - 发票开具失败
	InvoiceStatusFailed
	// InvoiceStatusCancelled 已取消 - 用户主动取消开票申请
	InvoiceStatusCancelled
	// InvoiceStatusMailed 已邮寄 - 发票已邮寄给用户
	InvoiceStatusMailed
	// InvoiceStatusSigned 已签收 - 用户已签收发票
	InvoiceStatusSigned
	// InvoiceStatusRejected 已驳回 - 开票申请被驳回
	InvoiceStatusRejected
	// InvoiceStatusPendingConfirm 待确认 - 开票申请需要进一步确认信息
	InvoiceStatusPendingConfirm
)

// String 返回开票状态的字符串表示
func (is InvoiceStatus) String() string {
	switch is {
	case InvoiceStatusApplied:
		return "applied"
	case InvoiceStatusProcessing:
		return "processing"
	case InvoiceStatusIssued:
		return "issued"
	case InvoiceStatusFailed:
		return "failed"
	case InvoiceStatusCancelled:
		return "cancelled"
	case InvoiceStatusMailed:
		return "mailed"
	case InvoiceStatusSigned:
		return "signed"
	case InvoiceStatusRejected:
		return "rejected"
	case InvoiceStatusPendingConfirm:
		return "pending_confirm"
	default:
		return "unknown"
	}
}

// CNString 返回开票状态的中文字符串表示
func (is InvoiceStatus) CNString() string {
	switch is {
	case InvoiceStatusApplied:
		return "已申请"
	case InvoiceStatusProcessing:
		return "处理中"
	case InvoiceStatusIssued:
		return "已开具"
	case InvoiceStatusFailed:
		return "开具失败"
	case InvoiceStatusCancelled:
		return "已取消"
	case InvoiceStatusMailed:
		return "已邮寄"
	case InvoiceStatusSigned:
		return "已签收"
	case InvoiceStatusRejected:
		return "已驳回"
	case InvoiceStatusPendingConfirm:
		return "待确认"
	default:
		return "未知"
	}
}

// IsValid 检查开票状态是否有效
func (is InvoiceStatus) IsValid() bool {
	return is >= InvoiceStatusUnknown && is <= InvoiceStatusPendingConfirm
}

// InvoiceStatusFromString 从字符串创建开票状态
func InvoiceStatusFromString(s string) InvoiceStatus {
	switch s {
	case "applied":
		return InvoiceStatusApplied
	case "processing":
		return InvoiceStatusProcessing
	case "issued":
		return InvoiceStatusIssued
	case "failed":
		return InvoiceStatusFailed
	case "cancelled":
		return InvoiceStatusCancelled
	case "mailed":
		return InvoiceStatusMailed
	case "signed":
		return InvoiceStatusSigned
	case "rejected":
		return InvoiceStatusRejected
	case "pending_confirm":
		return InvoiceStatusPendingConfirm
	default:
		return InvoiceStatusUnknown
	}
}

// InvoiceStatusFromCNString 从中文字符串创建开票状态
func InvoiceStatusFromCNString(s string) InvoiceStatus {
	switch s {
	case "已申请":
		return InvoiceStatusApplied
	case "处理中":
		return InvoiceStatusProcessing
	case "已开具":
		return InvoiceStatusIssued
	case "开具失败":
		return InvoiceStatusFailed
	case "已取消":
		return InvoiceStatusCancelled
	case "已邮寄":
		return InvoiceStatusMailed
	case "已签收":
		return InvoiceStatusSigned
	case "已驳回":
		return InvoiceStatusRejected
	case "待确认":
		return InvoiceStatusPendingConfirm
	default:
		return InvoiceStatusUnknown
	}
}
