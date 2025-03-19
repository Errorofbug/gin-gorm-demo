package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// UserDao 用户 DAO
type UserDao struct {
	*BaseDao
}

// NewUserDao 创建一个 UserDao 实例
func NewUserDao() *UserDao {
	return &UserDao{
		BaseDao: NewBaseDao(db, "user"),
	}
}

// User 用户表
type User struct {
	SoftDeleteBaseModel
	WorkNo     string `gorm:"column:work_no" json:"work_no"`       // 工号
	Password   string `gorm:"column:password" json:"password"`     // 密码
	RealName   string `gorm:"column:real_name" json:"real_name"`   // 真实姓名
	Role       int    `gorm:"column:role" json:"role"`             // 身份，0-管理员，1-教师，2-学生
	Sex        int    `gorm:"column:sex" json:"sex"`               // 性别，0-女，1-男
	Department string `gorm:"column:department" json:"department"` // 所在部门
	Email      string `gorm:"column:email" json:"email"`           // 电子邮箱
}

// GetByID 根据ID获取
func (dao *UserDao) GetByID(id int64) (*User, error) {
	var user User
	where := Where{
		"id":         id,
		"is_deleted": 0,
	}
	fields := Fields{"*"}
	appends := Appends{}
	err := dao.BaseSelectConvert(where, fields, appends, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Insert 新增记录
func (dao *UserDao) Insert(user *User) (int64, error) {
	lastInsertID, err := dao.BaseInsert(user)
	if err != nil && lastInsertID != 0 {
		return lastInsertID, err
	}
	return 0, nil
}

// Update 根据ID更新记录
func (dao *UserDao) Update(user User) (bool, error) {
	toMap, err := dao.ModelToMap(user)
	if err != nil {
		return false, err
	}
	affectedRows, err := dao.BaseUpdate(Where{"id": user.ID}, toMap)
	if err != nil || affectedRows == 0 {
		return false, err
	}
	return true, nil
}

// Delete 根据ID删除记录（软删）
func (dao *UserDao) Delete(id int64) (bool, error) {
	affectedRows, err := dao.BaseUpdate(Where{"id": id}, Updates{"is_deleted": 1})
	if err != nil || affectedRows == 0 {
		return false, err
	}
	return true, nil
}

// GetList 根据条件获取用户列表
func (dao *UserDao) GetList(where Where, offset, limit int) ([]*User, error) {
	var users []*User
	where["is_deleted"] = 0
	fields := Fields{"*"}
	appends := Appends{
		"limit":  limit,
		"offset": offset,
	}
	err := dao.BaseSelectConvert(where, fields, appends, &users)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}

// GetTotal 获取用户总数
func (dao *UserDao) GetTotal() (int64, error) {
	count, err := dao.Count(Where{"is_deleted": 0})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CheckAuth 鉴权
func (dao *UserDao) CheckAuth(workNo, password string) (int64, error) {
	var user User
	where := Where{"work_no": workNo, "password": password, "is_deleted": 0}
	fields := Fields{"id"}
	appends := Appends{}
	err := dao.BaseSelectConvert(where, fields, appends, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	if user.ID > 0 {
		return user.ID, nil
	}

	return 0, nil
}
