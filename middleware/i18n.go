package middleware

import (
	"gin-gorm-demo/common/locales"
	"gin-gorm-demo/conf"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// I18nMiddleware 国际化中间件
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		languageMatcher := language.NewMatcher([]language.Tag{
			language.English, // 默认语言
			language.Chinese, // 支持的语言
			// 添加更多支持的语言...
		})
		tag, _ := language.MatchStrings(languageMatcher, acceptLang)

		// 将匹配到的语言标签设置到上下文中
		c.Set("locales", tag.String())

		// 懒加载语言包
		locales.LazyLoadLang(tag.String(), conf.Settings.Lang.CodeConfFilePath)

		// 继续处理请求
		c.Next()
	}
}
