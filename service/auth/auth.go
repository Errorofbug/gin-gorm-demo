package auth

import (
	"gin-gorm-demo/common/entity"
	"gin-gorm-demo/models"
	"github.com/gin-gonic/gin"
)

// Service 鉴权服务
type Service struct {
	C *gin.Context
}

// Login 鉴权
func (s *Service) Login(dto entity.LoginDTO) (int64, error) {
	return models.NewUserDao().CheckAuth(dto.WorkNo, dto.Password)
}
