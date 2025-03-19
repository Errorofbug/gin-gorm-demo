package v1

import (
	"gin-gorm-demo/common/app"
	"gin-gorm-demo/common/e"
	"gin-gorm-demo/common/entity"
	"gin-gorm-demo/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserList 获取用户列表
//
//	@Summary	获取用户列表
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id			query		int64	false	"用户ID"
//	@Param		work_no		query		string	false	"用户工号"
//	@Param		real_name	query		string	false	"用户姓名"
//	@Param		role		query		int		false	"用户角色"
//	@Param		sex			query		int		false	"用户性别"
//	@Param		department	query		string	false	"用户部门"
//	@Param		page		query		int		false	"分页"
//	@Param		page_size	query		int		false	"分页大小"
//	@Success	200			{object}	app.Response{data=entity.GetUserListVO}
//	@Failure	400			{object}	app.Response
//	@Failure	500			{object}	app.Response
//	@Router		/api/v1/getuserlist [get]
func GetUserList(c *gin.Context) {
	appG := app.Gin{C: c}

	var req entity.GetUserListDTO
	if err := app.ValidateAndBind(c, &req); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	if c.Query("sex") == "" {
		req.Sex = -1
	}

	userService := user.Service{C: c}
	vo, err := userService.GetUserList(req)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetUsersFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, vo)
}
