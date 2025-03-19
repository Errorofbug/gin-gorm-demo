package models

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type (
	Where   = map[string]interface{} // 查询条件
	Fields  = []string               // 字段
	Appends = map[string]interface{} // 附加选项
	Updates = map[string]interface{} // 更新字段
)

// BaseDao 基础DAO
type BaseDao struct {
	DB        *gorm.DB // MySQL 连接池对象
	TableName string   // 表名
}

// NewBaseDao 创建BaseDao
func NewBaseDao(db *gorm.DB, tableName string) *BaseDao {
	return &BaseDao{
		DB:        db,
		TableName: tableName,
	}
}

// BaseSelect 基础查询，不进行对象绑定
// where: 查询条件，例如 Where{"id >": 1} 或 Where{"id": 1}
// fields: 查询字段，例如 Fields{"id", "name"}
// appends: 附加选项，例如 Appends{"ORDER BY": "id"}
func (dao *BaseDao) BaseSelect(where Where, fields Fields, appends Appends) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 将[]string转换为逗号分隔的字符串
	fieldsStr := "*"
	if len(fields) > 0 {
		fieldsStr = strings.Join(fields, ", ")
	}

	query := dao.DB.Table(dao.TableName).Select(fieldsStr)
	for key, value := range where {
		query = buildWhereClause(query, key, value)
	}

	// 附加选项
	query = applyAppends(query, appends)

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// BaseSelectConvert 基础查询，进行对象绑定
// where: 查询条件，例如 Where{"id >": 1} 或 Where{"id": 1}
// fields: 查询字段，例如 Fields{"id", "name"}
// appends: 附加选项，例如 Appends{"ORDER BY": "id"}
// out: 查询结果的绑定对象（指针）
func (dao *BaseDao) BaseSelectConvert(where Where, fields Fields, appends Appends, out interface{}) error {
	fieldsStr := "*"
	if len(fields) > 0 {
		fieldsStr = strings.Join(fields, ", ")
	}

	query := dao.DB.Table(dao.TableName).Select(fieldsStr)
	for key, value := range where {
		query = buildWhereClause(query, key, value)
	}

	// 附加选项
	query = applyAppends(query, appends)

	return query.Find(out).Error
}

// BaseSQLQuery 原生SQL查询并绑定结果
// sql: 需要执行的SQL语句，可包含占位符（如"WHERE id > ?"）
// cond: SQL参数，支持以下形式：
//   - 单个值：BaseSQLQuery("...WHERE id > ?", 1, &result)
//   - 切片：BaseSQLQuery("...WHERE id IN (?)", []int{1,2}, &result)
//   - map（命名参数）：BaseSQLQuery("...WHERE name=@name", map[string]interface{}{"name":"张三"}, &result)
//
// out: 查询结果绑定对象（指针）
func (dao *BaseDao) BaseSQLQuery(sql string, cond interface{}, out interface{}) error {
	// 创建原生SQL查询
	query := dao.DB.Raw(sql, cond)

	// 使用Scan方法绑定结果
	if err := query.Scan(out).Error; err != nil {
		// 处理特殊错误类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 清空输出对象
			reflect.ValueOf(out).Elem().Set(reflect.Zero(reflect.TypeOf(out).Elem()))
			return nil
		}
		return err
	}
	return nil
}

// BaseUpdate 基础更新
// where: 更新条件，例如 Where{"id = ?": 1}
// updates: 更新字段，例如 Updates{"name": "张三"}
func (dao *BaseDao) BaseUpdate(where Where, updates Updates) (int64, error) {
	query := dao.DB.Table(dao.TableName)
	for key, value := range where {
		query = buildWhereClause(query, key, value)
	}
	result := query.Updates(updates)
	return result.RowsAffected, result.Error
}

// BaseInsert 基础插入
// value: 插入的数据（Model对象）
func (dao *BaseDao) BaseInsert(value interface{}) (int64, error) {
	result := dao.DB.Table(dao.TableName).Create(value)
	if result.Error != nil {
		return 0, result.Error
	}

	// 获取插入最后一条记录的主键ID
	id := reflect.ValueOf(value).Elem().FieldByName("ID").Int()
	return id, nil
}

// BaseDelete 基础删除
// where: 删除条件，例如 Where{"id = ?": 1}
func (dao *BaseDao) BaseDelete(where Where) (int64, error) {
	query := dao.DB.Table(dao.TableName)
	for key, value := range where {
		query = buildWhereClause(query, key, value)
	}
	result := query.Delete(nil)
	return result.RowsAffected, result.Error
}

// Count 基础计数
// where: 计数条件，例如 Where{"id > ?": 1}
func (dao *BaseDao) Count(where Where) (int64, error) {
	var count int64
	query := dao.DB.Table(dao.TableName)
	for key, value := range where {
		query = buildWhereClause(query, key, value)
	}
	err := query.Count(&count).Error
	return count, err
}

// ModelToMap 将传入的Model转换为map[string]interface{}
// model: 需要转换的Model
func (dao *BaseDao) ModelToMap(model interface{}) (Updates, error) {
	// 获取模型的反射值和类型
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 解引用指针
	}

	// 遍历结构体字段，提取字段名和值
	modelMap := make(map[string]interface{})
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取字段的gorm标签
		gormTag := field.Tag.Get("gorm")
		if gormTag == "-" {
			continue // 忽略被标记为忽略的字段
		}

		// 获取字段名（使用gorm标签中的column名称）
		columnName := gormTag
		if columnName == "" {
			columnName = field.Name // 默认使用字段名
		}

		// 将字段值添加到 modelMap 中
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			modelMap[columnName] = fieldValue.Interface()
		}
	}

	return modelMap, nil
}

// applyAppends 处理附加查询选项
func applyAppends(query *gorm.DB, appends Appends) *gorm.DB {
	for key, value := range appends {
		switch key {
		case "order":
			if strValue, ok := value.(string); ok {
				query = query.Order(strValue)
			}
		case "limit":
			if intValue, ok := value.(int); ok {
				query = query.Limit(intValue)
			}
		case "offset":
			if intValue, ok := value.(int); ok {
				query = query.Offset(intValue)
			}
		case "group":
			if strValue, ok := value.(string); ok {
				query = query.Group(strValue)
			}
		case "having":
			if strValue, ok := value.(string); ok {
				query = query.Having(strValue)
			}
		case "distinct":
			if boolValue, ok := value.(bool); ok && boolValue {
				query = query.Distinct()
			}
		default:
		}
	}
	return query
}

// buildWhereClause 根据key和value构建WHERE子句
func buildWhereClause(query *gorm.DB, key string, value interface{}) *gorm.DB {
	// 检查 key 中是否包含符号
	if hasOperator(key) {
		// 如果包含符号，直接使用 key 和 value
		return query.Where(key+" ?", value)
	}
	// 如果不包含符号，默认使用等号
	return query.Where(key+" = ?", value)
}

// hasOperator 检查 key 中是否包含 SQL 操作符
func hasOperator(key string) bool {
	// 常见的 SQL 操作符
	operators := []string{"=", "!=", "<>", ">", "<", ">=", "<=", " LIKE ", " IN "}
	for _, op := range operators {
		if contains(key, op) {
			return true
		}
	}
	return false
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
