package locales

import (
	"fmt"
	"gin-gorm-demo/common/logging"
	"gin-gorm-demo/common/util"
	"gin-gorm-demo/conf"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"sync"
)

var (
	Bundle          *i18n.Bundle
	loadedLanguages = make(map[string]bool) // 用于缓存已加载的语言包
	mu              sync.Mutex              // 用于保护并发访问
)

// InitI18n 初始化i18n
func InitI18n() {
	if conf.Settings.App.DefaultLanguage == "zh" {
		Bundle = i18n.NewBundle(language.Chinese)
	} else {
		Bundle = i18n.NewBundle(language.English)
	}
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	path, err := util.GetFileFullPath(conf.Settings.Lang.DefaultLanguage,
		conf.Settings.Lang.CodeConfFilePath, conf.Settings.Lang.ConfFileExt)

	if err != nil {
		logging.Fatal(err)
	}
	_, err = Bundle.LoadMessageFile(path)
	if err != nil {
		logging.Fatal(err)
		return
	}
}

// LazyLoadLang 懒加载语言包
func LazyLoadLang(lang, confPath string) {
	mu.Lock()
	defer mu.Unlock()

	if !loadedLanguages[lang] {
		path, err := util.GetFileFullPath(lang, confPath, conf.Settings.Lang.ConfFileExt)
		fmt.Printf(path)

		if err != nil {
			logging.Error(err.Error())
			return
		}
		_, err = Bundle.LoadMessageFile(path)
		if err != nil {
			logging.Error(err.Error())
			return
		}
		loadedLanguages[lang] = true
	}
}
