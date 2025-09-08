package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// WalletType 钱包货币类型
type WalletType string

const (
	Coin    WalletType = "coin"
	Diamond WalletType = "diamond"
)

// UserWallet 用户钱包模型
type UserWallet struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Num       int64      `gorm:"not null;default:0" json:"num"`
	Type      WalletType `gorm:"size:20;not null" json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// TableName 指定表名
func (UserWallet) TableName() string {
	return "user_wallets"
}

// BeforeCreate 创建前的钩子
func (uw *UserWallet) BeforeCreate(scope *gorm.Scope) error {
	return nil
}
