package logging

import (
	"fmt"
	"gin-gorm-demo/common/util"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

// 日志级别
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// InitLog 设置日志记录器
func InitLog() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = util.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// Debug 输出调试级别的日志
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

// Info 输出信息级别的日志
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

// Warn 输出警告级别的日志
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

// Error 输出错误级别的日志
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

// Fatal 输出致命错误级别的日志
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

// setPrefix 设置日志输出的前缀
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
