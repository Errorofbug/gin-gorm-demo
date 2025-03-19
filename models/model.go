package models

import (
	"fmt"
	"gorm.io/gorm/schema"
	"log"
	"time"

	"gin-gorm-demo/common/logging"
	"gin-gorm-demo/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// BaseModel 基础模型
type BaseModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedTime time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

// SoftDeleteBaseModel 带软删（使用新版软删除特性）
type SoftDeleteBaseModel struct {
	BaseModel
	IsDeleted bool `gorm:"column:is_deleted;default:0" json:"is_deleted"`
}

// InitDB 连接数据库（新版配置）
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Settings.Database.User,
		conf.Settings.Database.Password,
		conf.Settings.Database.Server,
		conf.Settings.Database.Name)

	gormConfig := &gorm.Config{
		// 单数表名配置
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// 注册全局回调
	registerCallbacks(db)

	// 连接池配置
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

// 注册回调函数
func registerCallbacks(db *gorm.DB) {
	// 错误日志回调
	registerErrorCallback(db, "create")
	registerErrorCallback(db, "update")
	registerErrorCallback(db, "delete")
	registerErrorCallback(db, "query")
	registerErrorCallback(db, "row")
	registerErrorCallback(db, "raw")

	// 注册钩子来处理创建和更新的时间戳
	db.Callback().Create().Before("gorm:create").Register("set_created_time", setCreatedTime)
	db.Callback().Update().Before("gorm:update").Register("set_updated_time", setUpdatedTime)
}

// 错误日志回调注册
func registerErrorCallback(db *gorm.DB, operation string) {
	callbackName := "log_after_" + operation
	afterHookName := "gorm:after_" + operation

	switch operation {
	case "create":
		db.Callback().Create().After(afterHookName).Register(callbackName, logAfterSQL)
	case "update":
		db.Callback().Update().After(afterHookName).Register(callbackName, logAfterSQL)
	case "delete":
		db.Callback().Delete().After(afterHookName).Register(callbackName, logAfterSQL)
	case "query":
		db.Callback().Query().After(afterHookName).Register(callbackName, logAfterSQL)
	case "row":
		db.Callback().Row().After("gorm:row").Register(callbackName, logAfterSQL)
	case "raw":
		db.Callback().Raw().After("gorm:raw").Register(callbackName, logAfterSQL)
	}
}

// 日志回调函数
func logAfterSQL(tx *gorm.DB) {
	if tx.Error != nil {
		// 排除记录不存在的警告
		if tx.Error != gorm.ErrRecordNotFound {
			logging.Error(tx.Error)
		}
	} else {
		logStr := fmt.Sprintf("[Table: %s]SQL: \"%s\", Vars: %v",
			tx.Statement.Table, tx.Statement.SQL.String(), tx.Statement.Vars)
		logging.Debug(logStr)
	}
}

// 设置创建时间戳
func setCreatedTime(tx *gorm.DB) {
	if tx.Statement.Schema != nil {
		now := time.Now()

		// 设置 CreatedTime
		if createTimeField := tx.Statement.Schema.LookUpField("CreatedTime"); createTimeField != nil {
			if tx.Statement.ReflectValue.CanAddr() {
				if createTimeField.GORMDataType == "time" {
					_ = createTimeField.Set(tx.Statement.Context, tx.Statement.ReflectValue, now)
				}
			}
		}

		// 设置 UpdatedTime
		if updateTimeField := tx.Statement.Schema.LookUpField("UpdatedTime"); updateTimeField != nil {
			if tx.Statement.ReflectValue.CanAddr() {
				if updateTimeField.GORMDataType == "time" {
					_ = updateTimeField.Set(tx.Statement.Context, tx.Statement.ReflectValue, now)
				}
			}
		}
	}
}

// 设置更新时间戳
func setUpdatedTime(tx *gorm.DB) {
	if tx.Statement.Schema != nil {
		now := time.Now()

		// 设置 UpdatedTime
		if updateTimeField := tx.Statement.Schema.LookUpField("UpdatedTime"); updateTimeField != nil {
			if tx.Statement.ReflectValue.CanAddr() {
				if updateTimeField.GORMDataType == "time" {
					_ = updateTimeField.Set(tx.Statement.Context, tx.Statement.ReflectValue, now)
				}
			}
		}
	}
}
