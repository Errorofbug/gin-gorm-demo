package gredis

// Redis前缀key
const (
	PrefixAuthToken = "auth:token"
)

// Redis key过期时间
const (
	TimeoutAuthToken = 7 * 24 * 60 * 60
)
