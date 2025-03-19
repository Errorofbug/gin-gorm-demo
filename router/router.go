package router

import (
	v1 "gin-gorm-demo/api/v1"
	"gin-gorm-demo/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.I18nMiddleware())

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/login", v1.Login)
		apiv1.Use(middleware.AuthMiddleware())
		{
			// 获取用户列表
			apiv1.GET("/getuserlist", v1.GetUserList)
		}
	}

	return r
}
