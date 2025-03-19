package app

import (
	"gin-gorm-demo/common/logging"
)

// MarkErrors 日志记录错误
func MarkErrors(errors []ValidationError) {
	for _, err := range errors {
		logging.Info(err.Field, err.Message)
	}
}
