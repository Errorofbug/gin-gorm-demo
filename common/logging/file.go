package logging

import (
	"fmt"
	"gin-gorm-demo/conf"
	"time"
)

// getLogFilePath 获取日志文件路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", conf.Settings.App.RuntimeRootPath, conf.Settings.App.LogSavePath)
}

// getLogFileName 获取日志文件名
func getLogFileName() string {
	return fmt.Sprintf("%s_%s.%s",
		conf.Settings.App.LogSaveName,
		time.Now().Format(conf.Settings.App.TimeFormat),
		conf.Settings.App.LogFileExt,
	)
}
