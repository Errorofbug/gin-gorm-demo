package entity

// LoginDTO 登录请求
type LoginDTO struct {
	WorkNo   string `json:"work_no" validate:"required,max=10"`
	Password string `json:"password" validate:"required,max=255"`
}

// LoginVO 登录响应
type LoginVO struct {
	Token string `json:"token"`
}
