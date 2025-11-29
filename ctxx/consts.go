package ctxx

// Define key types for context values to avoid conflicts
// 注意：字符串类型的值使用 metainfo 存储，只有 bool 类型需要使用 context.WithValue
type (
	tenantIsolationKey   struct{}
	merchantIsolationKey struct{}
)

const (
	// 上下下文键名定义 - 统一用于内部上下文管理和RPC调用间传递
	// RPC上下文传递键名定义
	TenantKey            = "tenant_id"
	UserKey              = "user_id"
	RequestKey           = "request_id"
	MerchantKey          = "merchant_id"
	TenantIsolationKey   = "tenant_isolation"
	MerchantIsolationKey = "merchant_isolation"
	MemberKey            = "member_id"
	DonorKey             = "donor_id"
	AppTypeKey           = "app_type"
	ExpandedKey          = "expanded_info"
	TenantNameKey        = "tenant_name"
	MemberNameKey        = "member_name"
	MerchantNameKey      = "merchant_name"
	UserNameKey          = "user_name"
	DonorNameKey         = "donor_name"
	AppNameKey           = "app_name"
	AppIdKey             = "app_id"
	IpKey                = "ip"
	UserAgentKey         = "user_agent"

	// app type
	AppMerchant = "merchant"
	AppMember   = "member"
	AppDonor    = "donor"
	AppAdmin    = "admin"
	AppPlatform = "platform"
	AppCallback = "callback"
)
