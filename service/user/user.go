package user

import (
	"gin-gorm-demo/common/entity"
	"gin-gorm-demo/common/util"
	"gin-gorm-demo/models"
	"github.com/gin-gonic/gin"
)

// Service 用户Service
type Service struct {
	C *gin.Context
}

// GetUserList 获取用户列表
func (s *Service) GetUserList(dto entity.GetUserListDTO) (entity.GetUserListVO, error) {
	var (
		users []*models.User
	)

	where := make(map[string]interface{})
	if dto.ID != 0 {
		where["id"] = dto.ID
	}
	if dto.WorkNo != "" {
		where["work_no"] = dto.WorkNo
	}
	if dto.RealName != "" {
		where["real_name"] = dto.RealName
	}
	if dto.Sex != -1 {
		where["sex"] = dto.Sex
	}
	if dto.Role != 0 {
		where["role"] = dto.Role
	}
	if dto.Department != "" {
		where["department"] = dto.Department
	}

	users, err := models.NewUserDao().GetList(where, util.GetPage(dto.Page), util.GetPageSize(dto.PageSize))
	if err != nil {
		return entity.GetUserListVO{}, err
	}

	count, err := models.NewUserDao().GetTotal()
	if err != nil {
		return entity.GetUserListVO{}, err
	}

	vo := entity.GetUserListVO{
		Total: count,
	}
	for _, v := range users {
		item := entity.GetUserListItemVO{
			WorkNo:     v.WorkNo,
			RealName:   v.RealName,
			Role:       v.Role,
			Sex:        v.Sex,
			Department: v.Department,
			Email:      v.Email,
			CreateTime: util.FormatTime(v.CreatedTime),
		}
		vo.List = append(vo.List, item)
	}

	return vo, nil
}
