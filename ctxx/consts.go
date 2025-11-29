package ctxx

// Define key types for context values to avoid conflicts
type (
	tenantKey            struct{}
	userIDKey            struct{}
	requestIDKey         struct{}
	merchantIDKey        struct{}
	tenantIsolationKey   struct{}
	merchantIsolationKey struct{}
	memberIDKey          struct{}
	donorIDKey           struct{}
	appTypeKey           struct{}
	expandedKey          struct{}
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

	// app type
	AppMerchant = "merchant"
	AppMember   = "member"
	AppDonor    = "donor"
	AppAdmin    = "admin"
	AppPlatform = "platform"
	AppCallback = "callback"
)
