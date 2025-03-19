package v1

import (
	"gin-gorm-demo/common/app"
	"gin-gorm-demo/common/e"
	"gin-gorm-demo/common/entity"
	"gin-gorm-demo/common/gredis"
	"gin-gorm-demo/common/util"
	"gin-gorm-demo/service/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 登录
//
//	@Summary	登录
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		work_no		body		string	true	"用户工号"
//	@Param		password	body		string	true	"密码"
//	@Success	200			{object}	app.Response{data=entity.LoginVO}
//	@Success	400			{object}	app.Response
//	@Failure	500			{object}	app.Response
//	@Router		/api/v1/login [post]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}

	var req entity.LoginDTO

	if err := app.ValidateAndBind(c, &req); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	authService := auth.Service{C: c}
	id, err := authService.Login(req)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthCheckTokenFail, nil)
		return
	}

	if id == 0 {
		appG.Response(http.StatusUnauthorized, e.ErrorAuthToken, nil)
		return
	}

	token, err := util.GenerateToken(id)
	// 这里也可以不用jwt，直接使用md5生成token
	//token := util.EncodeMD5(util.ToStr(id))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthToken, nil)
		return
	}

	key := gredis.GenerateKey(gredis.PrefixAuthToken, token)
	err = gredis.Set(key, token, gredis.TimeoutAuthToken)
	if err != nil {
		return
	}

	vo := entity.LoginVO{
		Token: token,
	}

	appG.Response(http.StatusOK, e.Success, vo)
}
