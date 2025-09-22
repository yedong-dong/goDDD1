package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User 用户模型
type User struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	UID        uint       `gorm:"not null;unique" json:"uid"` // 用户唯一标识，从10000开始自增
	Username   string     `gorm:"size:50;not null;unique" json:"username"`
	Email      string     `gorm:"size:100;not null;unique" json:"email"`
	Password   string     `gorm:"size:100;not null" json:"password"`
	Level      uint       `gorm:"default:1" json:"level"`                // 用户等级，默认为1级
	Experience uint       `gorm:"default:0" json:"experience"`           // 用户经验值
	TotalSpent uint       `gorm:"default:0" json:"total_spent"`          // 用户总消费金额
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	IsDeleted  string     `gorm:"default:0;size:1" json:"is_deleted"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	// 如果UID为空，则自动生成从10000开始的UID
	if u.UID == 0 {
		var maxUID uint
		// 查询当前最大的UID值
		scope.DB().Model(&User{}).Select("COALESCE(MAX(uid), 9999)").Row().Scan(&maxUID)
		// 设置新的UID值，确保从10000开始
		if maxUID < 10000 {
			u.UID = 10000
		} else {
			u.UID = maxUID + 1
		}
	}
	
	// 设置默认等级为1
	if u.Level == 0 {
		u.Level = 1
	}
	
	// 这里可以添加密码加密等逻辑
	return nil
}
