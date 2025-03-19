package entity

// GetUserListDTO 获取用户列表请求
type GetUserListDTO struct {
	ID         int64  `form:"id" validate:"omitempty,min=1"`
	WorkNo     string `form:"work_no" validate:"omitempty,max=10"`
	RealName   string `form:"real_name" validate:"omitempty,max=16"`
	Role       int    `form:"role" validate:"omitempty,oneof=1 2"`
	Sex        int    `form:"sex" validate:"omitempty,oneof=0 1"`
	Department string `form:"department" validate:"omitempty,max=32"`

	Page     int `form:"page" validate:"omitempty,min=1"`
	PageSize int `form:"page_size" validate:"omitempty,min=1"`
}

// GetUserListItemVO 用户列表元素
type GetUserListItemVO struct {
	WorkNo     string `json:"work_no"`
	RealName   string `json:"real_name"`
	Role       int    `json:"role"`
	Sex        int    `json:"sex"`
	Department string `json:"department"`
	Email      string `json:"email"`
	CreateTime string `json:"create_time"`
}

// GetUserListVO 获取用户列表响应
type GetUserListVO struct {
	List  []GetUserListItemVO `json:"list"`
	Total int64               `json:"total"`
}
