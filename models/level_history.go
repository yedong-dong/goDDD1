package models

import (
	"time"
)

// LevelHistory 用户等级历史记录
type LevelHistory struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	UserID          uint      `gorm:"not null" json:"user_id"`     // 用户ID
	OldLevel        uint      `gorm:"not null" json:"old_level"`   // 旧等级
	NewLevel        uint      `gorm:"not null" json:"new_level"`   // 新等级
	ExpGained       uint      `json:"exp_gained"`                  // 获得的经验值
	Experience      uint      `json:"experience"`                  // 升级后的经验值
	CoinRewarded    uint      `json:"coin_rewarded"`               // 奖励的金币
	DiamondRewarded uint      `json:"diamond_rewarded"`            // 奖励的钻石
	Description     string    `gorm:"size:255" json:"description"` // 升级描述
	CreatedAt       time.Time `json:"created_at"`
}

// TableName 指定表名
func (LevelHistory) TableName() string {
	return "level_histories"
}
