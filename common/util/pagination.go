package util

import (
	"gin-gorm-demo/conf"
)

// GetPage 获取分页参数
func GetPage(page int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * conf.Settings.App.PageSize
	}

	return result
}

// GetPageSize 获取分页大小参数
func GetPageSize(pageSize int) int {
	if pageSize > 0 {
		return pageSize
	}
	return conf.Settings.App.PageSize
}
