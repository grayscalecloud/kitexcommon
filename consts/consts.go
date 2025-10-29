package consts

import "time"

const (
	Admin           = "Admin"
	User            = "User"
	ThirtyDays      = time.Hour * 24 * 30
	AccountID       = "accountID"
	From            = "x-from"
	ID              = "id"
	Language        = "language"
	DefaultLanguage = "zh"
	UserID          = "user_id"     // 用户ID常量
	TenantID        = "tenant_id"   // 租户ID常量
	TraceID         = "traceID"     // 追踪ID常量
	MerchantID      = "merchant_id" // 商户ID常量

	HlogFilePath = "./tmp/hlog/logs/"
	KlogFilePath = "./tmp/klog/logs/"

	MySqlDSN    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	MongoURI    = "mongodb://%s:%s@%s:%d"
	RabbitMqURI = "amqp://%s:%s@%s:%d/"

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	TCP = "tcp"

	FreePortAddress = "localhost:0"
	CorsAddress     = "http://localhost:3000"

	ConsulCheckInterval                       = "7s"
	ConsulCheckTimeout                        = "5s"
	ConsulCheckDeregisterCriticalServiceAfter = "15s"

	RedisProfileClientDB = 0
	RedisBlobClientDB    = 0
	RedisCarClientDB     = 0

	UserSnowflakeNode  = 2
	BlobSnowflakeNode  = 3
	AdminSnowflakeNode = 4

	LimitOfSomeCars     = 20
	LimitOfSomeTrips    = 20
	LimitOfSomeProfiles = 20
	LimitOfSomeUsers    = 20
)
