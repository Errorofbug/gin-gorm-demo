package app

import (
	"gin-gorm-demo/common/e"
	"github.com/gin-gonic/gin"
)

// Gin gin应用
type Gin struct {
	C *gin.Context
}

// Response HTTP响应
type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

// Response 格式化返回HTTP响应
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		ErrNo:  errCode,
		ErrMsg: e.GetErrMsg(g.C, errCode),
		Data:   data,
	})
	return
}
