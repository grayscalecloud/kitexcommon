package enum

// ProductStatus 商品状态枚举
// 用于表示商品的各种状态
type ProductStatus int

const (
	// ProductStatusUnknown 未知状态
	ProductStatusUnknown ProductStatus = 0
	// ProductStatusOnline 上架 - 商品正在销售中
	ProductStatusOnline ProductStatus = 1
	// ProductStatusOffline 下架 - 商品已下架，不再销售
	ProductStatusOffline ProductStatus = 2
	// ProductStatusPending 待审核 - 商品提交审核中
	ProductStatusPending ProductStatus = 3
	// ProductStatusRejected 审核拒绝 - 商品审核被拒绝
	ProductStatusRejected ProductStatus = 4
)

// String 返回商品状态的字符串表示
func (ps ProductStatus) String() string {
	switch ps {
	case ProductStatusOnline:
		return "online"
	case ProductStatusOffline:
		return "offline"
	case ProductStatusPending:
		return "pending"
	case ProductStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

// CNString 返回商品状态的中文字符串表示
func (ps ProductStatus) CNString() string {
	switch ps {
	case ProductStatusOnline:
		return "上架"
	case ProductStatusOffline:
		return "下架"
	case ProductStatusPending:
		return "待审核"
	case ProductStatusRejected:
		return "审核拒绝"
	default:
		return "未知"
	}
}

// IsValid 检查商品状态是否有效
func (ps ProductStatus) IsValid() bool {
	return ps == ProductStatusUnknown ||
		ps == ProductStatusOnline ||
		ps == ProductStatusOffline ||
		ps == ProductStatusPending ||
		ps == ProductStatusRejected
}

// ProductStatusFromString 从字符串创建商品状态
func ProductStatusFromString(s string) ProductStatus {
	switch s {
	case "online":
		return ProductStatusOnline
	case "offline":
		return ProductStatusOffline
	case "pending":
		return ProductStatusPending
	case "rejected":
		return ProductStatusRejected
	default:
		return ProductStatusUnknown
	}
}

// ProductStatusFromCNString 从中文字符串创建商品状态
func ProductStatusFromCNString(s string) ProductStatus {
	switch s {
	case "上架":
		return ProductStatusOnline
	case "下架":
		return ProductStatusOffline
	case "待审核":
		return ProductStatusPending
	case "审核拒绝":
		return ProductStatusRejected
	default:
		return ProductStatusUnknown
	}
}
