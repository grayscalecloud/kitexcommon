package enum

// SubmissionStatus 进件状态枚举
// 用于表示商户或用户提交申请的处理状态，如支付机构商户入驻申请等场景
type SubmissionStatus int

const (
	// SubmissionStatusUnknown 未知状态 - 初始状态或异常状态
	SubmissionStatusUnknown SubmissionStatus = iota
	// SubmissionStatusSubmitted 已提交 - 用户已提交申请，等待系统处理
	SubmissionStatusSubmitted
	// SubmissionStatusUnderReview 审核中 - 申请正在人工或系统审核中
	SubmissionStatusUnderReview
	// SubmissionStatusApproved 已批准 - 申请已通过审核
	SubmissionStatusApproved
	// SubmissionStatusRejected 已拒绝 - 申请被拒绝
	SubmissionStatusRejected
	// SubmissionStatusCancelled 已取消 - 用户主动取消申请
	SubmissionStatusCancelled
	// SubmissionStatusExpired 已过期 - 申请长时间未处理导致过期
	SubmissionStatusExpired
	// SubmissionStatusPendingSupplement 待补充材料 - 申请材料不完整，需要用户补充
	SubmissionStatusPendingSupplement
	// SubmissionStatusSupplementing 补充材料中 - 用户正在补充材料
	SubmissionStatusSupplementing
	// SubmissionStatusSuccess 进件成功 - 申请处理完成且成功
	SubmissionStatusSuccess
)

// String 返回进件状态的字符串表示
func (ss SubmissionStatus) String() string {
	switch ss {
	case SubmissionStatusSubmitted:
		return "submitted"
	case SubmissionStatusUnderReview:
		return "under_review"
	case SubmissionStatusApproved:
		return "approved"
	case SubmissionStatusRejected:
		return "rejected"
	case SubmissionStatusCancelled:
		return "cancelled"
	case SubmissionStatusExpired:
		return "expired"
	case SubmissionStatusPendingSupplement:
		return "pending_supplement"
	case SubmissionStatusSupplementing:
		return "supplementing"
	case SubmissionStatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

// CNString 返回进件状态的中文字符串表示
func (ss SubmissionStatus) CNString() string {
	switch ss {
	case SubmissionStatusSubmitted:
		return "已提交"
	case SubmissionStatusUnderReview:
		return "审核中"
	case SubmissionStatusApproved:
		return "已批准"
	case SubmissionStatusRejected:
		return "已拒绝"
	case SubmissionStatusCancelled:
		return "已取消"
	case SubmissionStatusExpired:
		return "已过期"
	case SubmissionStatusPendingSupplement:
		return "待补充材料"
	case SubmissionStatusSupplementing:
		return "补充材料中"
	case SubmissionStatusSuccess:
		return "进件成功"
	default:
		return "未知"
	}
}

// IsValid 检查进件状态是否有效
func (ss SubmissionStatus) IsValid() bool {
	return ss >= SubmissionStatusUnknown && ss <= SubmissionStatusSuccess
}

// SubmissionStatusFromString 从字符串创建进件状态
func SubmissionStatusFromString(s string) SubmissionStatus {
	switch s {
	case "submitted":
		return SubmissionStatusSubmitted
	case "under_review":
		return SubmissionStatusUnderReview
	case "approved":
		return SubmissionStatusApproved
	case "rejected":
		return SubmissionStatusRejected
	case "cancelled":
		return SubmissionStatusCancelled
	case "expired":
		return SubmissionStatusExpired
	case "pending_supplement":
		return SubmissionStatusPendingSupplement
	case "supplementing":
		return SubmissionStatusSupplementing
	case "success":
		return SubmissionStatusSuccess
	default:
		return SubmissionStatusUnknown
	}
}

// SubmissionStatusFromCNString 从中文字符串创建进件状态
func SubmissionStatusFromCNString(s string) SubmissionStatus {
	switch s {
	case "已提交":
		return SubmissionStatusSubmitted
	case "审核中":
		return SubmissionStatusUnderReview
	case "已批准":
		return SubmissionStatusApproved
	case "已拒绝":
		return SubmissionStatusRejected
	case "已取消":
		return SubmissionStatusCancelled
	case "已过期":
		return SubmissionStatusExpired
	case "待补充材料":
		return SubmissionStatusPendingSupplement
	case "补充材料中":
		return SubmissionStatusSupplementing
	case "进件成功":
		return SubmissionStatusSuccess
	default:
		return SubmissionStatusUnknown
	}
}
