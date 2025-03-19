package middleware

import (
	"gin-gorm-demo/common/app"
	"gin-gorm-demo/common/e"
	"gin-gorm-demo/common/gredis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		appG := app.Gin{C: c}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appG.Response(http.StatusUnauthorized, e.InvalidParams, data)
			c.Abort()
			return
		}

		// 解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			appG.Response(http.StatusUnauthorized, e.InvalidParams, data)
			c.Abort()
			return
		}

		// 获取缓存中的token，二者进行比对
		token := parts[1]
		key := gredis.GenerateKey(gredis.PrefixAuthToken, token)
		cachedToken, err := gredis.Get(key)
		if err != nil {
			appG.Response(http.StatusUnauthorized, e.ErrorAuthCheckTokenFail, data)
			c.Abort()
			return
		}
		if string(cachedToken) != token {
			appG.Response(http.StatusUnauthorized, e.ErrorAuthCheckTokenFail, data)
			c.Abort()
			return
		}

		// 刷新token
		ttl, err := gredis.TTL(key)
		if err != nil {
			appG.Response(http.StatusUnauthorized, e.ErrorAuthCheckTokenFail, data)
			c.Abort()
			return
		}

		if ttl < 24*60*60 {
			err := gredis.Expire(key, gredis.TimeoutAuthToken)
			if err != nil {
				appG.Response(http.StatusUnauthorized, e.ErrorAuthCheckTokenTimeout, data)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
