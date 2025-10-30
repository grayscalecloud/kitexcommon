package ctxx

// Define key types for context values to avoid conflicts
type (
	tenantKey            struct{}
	userIDKey            struct{}
	requestIDKey         struct{}
	merchantIDKey        struct{}
	tenantIsolationKey   struct{}
	merchantIsolationKey struct{}
)

const (
	// 上下文键名定义 - 统一用于内部上下文管理和RPC调用间传递
	// RPC上下文传递键名定义
	TenantKey            = "tenant_id"
	UserKey              = "user_id"
	RequestKey           = "request_id"
	MerchantKey          = "merchant_id"
	TenantIsolationKey   = "tenant_isolation"
	MerchantIsolationKey = "merchant_isolation"
)
