package e

import (
	"fmt"
	"gin-gorm-demo/common/locales"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetErrMsg(c *gin.Context, errno int) string {
	lang := c.MustGet("locales").(string)
	loc := i18n.NewLocalizer(locales.Bundle, lang)
	errMsg := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: fmt.Sprintf("%d", errno),
	})

	return errMsg
}
