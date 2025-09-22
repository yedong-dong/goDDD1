package models

import (
	"time"
)

// LevelConfig 等级配置模型
type LevelConfig struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	Level           uint      `gorm:"not null;unique" json:"level"`           // 等级
	RequiredExp     uint      `gorm:"not null" json:"required_exp"`           // 所需经验值
	CoinReward      uint      `gorm:"default:0" json:"coin_reward"`           // 金币奖励
	DiamondReward   uint      `gorm:"default:0" json:"diamond_reward"`        // 钻石奖励
	DiscountPercent uint      `gorm:"default:100" json:"discount_percent"`    // 折扣百分比(100表示无折扣)
	Description     string    `gorm:"size:255" json:"description"`            // 等级描述
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (LevelConfig) TableName() string {
	return "level_configs"
}
