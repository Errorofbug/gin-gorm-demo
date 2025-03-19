package main

import (
	"fmt"
	"gin-gorm-demo/common/gredis"
	"gin-gorm-demo/common/locales"
	"gin-gorm-demo/common/logging"
	"gin-gorm-demo/conf"
	"gin-gorm-demo/models"
	"gin-gorm-demo/router"
	"github.com/gin-gonic/gin"
	"net/http"

	_ "gin-gorm-demo/docs"
)

// init 初始化服务
func init() {
	// 注意初始化顺序
	conf.InitSettings()
	logging.InitLog()
	models.InitDB()
	gredis.InitRedis()
	locales.InitI18n()
}

//	@title			Golang Gin API
//	@version		1.0
//	@description	Gin+Gorm项目脚手架

//	@host		localhost:9000
//	@BasePath	/api/v1

// @license.name	MIT
func main() {
	gin.SetMode(conf.Settings.Server.RunMode)

	r := router.InitRouter()
	readTimeout := conf.Settings.Server.ReadTimeout
	writeTimeout := conf.Settings.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", conf.Settings.Server.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Info("HTTP Server Listening on " + endPoint)
	server.ListenAndServe()

	// 实现优雅关机
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
