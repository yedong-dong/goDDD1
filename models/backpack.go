package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Backpack struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`      // 用户ID
	StoreID   uint       `gorm:"not null;index" json:"store_id"`     // 商品ID
	Quantity  int64      `gorm:"not null;default:0" json:"quantity"` // 拥有数量
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	// 关联关系
	User  User  `gorm:"foreignkey:UserID" json:"user,omitempty"`   // 关联用户
	Store Store `gorm:"foreignkey:StoreID" json:"store,omitempty"` // 关联商品
}

func (Backpack) TableName() string {
	return "backpacks"
}

func (b *Backpack) BeforeCreate(scope *gorm.Scope) error {
	if b.Quantity < 0 {
		b.Quantity = 0
	}
	return nil
}

type BackpackItem struct {
	ID       uint   `json:"id"`
	StoreID  uint   `json:"store_id"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Price    int64  `json:"price"`
	CostType string `json:"cost_type"`
}
