package app

import (
	"errors"
	"fmt"
	"gin-gorm-demo/common/e"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

// ValidateAndBind 绑定请求参数并验证
func ValidateAndBind(c *gin.Context, req interface{}) error {
	// 绑定查询参数到结构体中
	if err := c.ShouldBind(req); err != nil {
		return err
	}

	// 创建一个验证器
	validate := validator.New()

	// 验证结构体
	if err := validate.Struct(req); err != nil {
		// 将验证错误转换为自定义格式
		validationErrors := TranslateValidationErrors(err)
		// 日志记录错误
		MarkErrors(validationErrors)
		return errors.New(e.GetErrMsg(c, e.InvalidParams))
	}

	return nil
}

// ValidationError 自定义验证错误结构
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误信息
}

// TranslateValidationErrors 将 validator.ValidationErrors 转换为自定义格式
func TranslateValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			// 生成详细的错误信息
			message := GenerateErrorMessage(e)
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(), // 字段名
				Message: message,   // 自定义错误信息
			})
		}
	}

	return validationErrors
}

// GenerateErrorMessage 生成详细的错误信息
func GenerateErrorMessage(e validator.FieldError) string {
	tag := e.Tag()
	param := e.Param()
	value := e.Value()

	switch tag {
	case "min":
		return "min is " + param + " (current value: " + formatReqValue(value) + ")"
	case "max":
		return "max is " + param + " (current value: " + formatReqValue(value) + ")"
	case "oneof":
		return "must be one of " + param + " (current value: " + formatReqValue(value) + ")"
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address (current value: " + formatReqValue(value) + ")"
	case "len":
		return "length must be " + param + " (current length: " + formatReqValue(value) + ")"
	default:
		return "failed validation rule: " + tag
	}
}

// formatReqValue 将字段值转换为字符串
func formatReqValue(value interface{}) string {
	return strings.TrimSpace(strings.Trim(fmt.Sprintf("%v", value), "[]"))
}
