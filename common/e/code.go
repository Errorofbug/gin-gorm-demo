package e

// 请求响应码
const (
	Success       = 200
	InvalidParams = 400
	Error         = 500
)

// 鉴权响应码
const (
	ErrorAuthCheckTokenFail    = 20001 // Token鉴权失败
	ErrorAuthCheckTokenTimeout = 20002 // Token鉴权超时
	ErrorAuthToken             = 20003 // 用户名或密码错误
)

// 用户接口
const (
	ErrorGetUsersFail = 30001 // 获取用户列表失败
	ErrorGetUserFail  = 30002 // 获取用户信息失败
)
