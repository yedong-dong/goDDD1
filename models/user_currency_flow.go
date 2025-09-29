package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserCurrencyFlow struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	StoreID     uint      `gorm:"not null" json:"store_id"`
	CostType    string    `gorm:"size:20;not null" json:"cost_type"`
	Description string    `gorm:"size:255;not null" json:"description"`
	Price       int64     `gorm:"not null" json:"price"`
	Ctime       time.Time `gorm:"not null" json:"ctime"`
}

func (UserCurrencyFlow) TableName() string {
	return "user_currency_flow"
}

func (u *UserCurrencyFlow) BeforeCreate(scope *gorm.Scope) error {
	if u.Ctime.IsZero() {
		u.Ctime = time.Now()
	}
	return nil
}

type UserCurrencyFlowDTO struct {
	UserID   uint   `gorm:"not null" json:"user_id"`
	StoreID  uint   `gorm:"not null" json:"store_id"`
	CostType string `gorm:"size:20;not null" json:"cost_type"`
	Price    int64  `gorm:"not null" json:"price"`
}
