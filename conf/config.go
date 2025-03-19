package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
	"time"
)

// App 应用配置
type App struct {
	PageSize        int    `toml:"pageSize"`
	PrefixUrl       string `toml:"prefixUrl"`
	RuntimeRootPath string `toml:"runtimeRootPath"`
	LogSavePath     string `toml:"logSavePath"`
	LogSaveName     string `toml:"logSaveName"`
	LogFileExt      string `toml:"logFileExt"`
	TimeFormat      string `toml:"timeFormat"`
	DefaultLanguage string `toml:"defaultLanguage"`
}

// Server 服务器配置
type Server struct {
	RunMode      string        `toml:"runMode"`
	HttpPort     int           `toml:"httpPort"`
	ReadTimeout  time.Duration `toml:"readTimeout"`
	WriteTimeout time.Duration `toml:"writeTimeout"`
}

// Database mysql配置
type Database struct {
	Type     string `toml:"type"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Server   string `toml:"server"`
}

// Redis redis配置
type Redis struct {
	Host        string        `toml:"host"`
	Password    string        `toml:"password"`
	MaxIdle     int           `toml:"maxIdle"`
	MaxActive   int           `toml:"maxActive"`
	IdleTimeout time.Duration `toml:"idleTimeout"`
}

// Lang lang配置
type Lang struct {
	CodeConfFilePath string `toml:"codeConfFilePath"`
	ConfFileExt      string `toml:"confFileExt"`
	DefaultLanguage  string `toml:"defaultLanguage"`
	DefaultTimeZone  string `toml:"defaultTimeZone"`
}

// Config 全局配置
type Config struct {
	App      `toml:"app"`
	Server   `toml:"server"`
	Database `toml:"database"`
	Redis    `toml:"redis"`
	Lang     `toml:"lang"`
}

var (
	Settings       = &Config{}
	configFilePath = "./conf/conf.toml"
)

// InitSettings 初始化全局配置
func InitSettings() {
	path := flag.String("config_path", configFilePath, "config file")
	flag.Parse()

	if path == nil || len(*path) == 0 {
		os.Exit(1)
	}

	configPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if _, err := toml.DecodeFile(configPath, Settings); err != nil {
		log.Fatalf(err.Error())
	}

	Settings.Server.ReadTimeout = Settings.Server.ReadTimeout * time.Second
	Settings.Server.WriteTimeout = Settings.Server.WriteTimeout * time.Second
}
