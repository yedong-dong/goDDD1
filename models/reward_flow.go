package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// RewardFlowType 奖励类型
type RewardFlowType string

const (
	RewardTypeItem    RewardFlowType = "item"    // 物品奖励
	RewardTypeCoin    RewardFlowType = "coin"    // 金币奖励
	RewardTypeDiamond RewardFlowType = "diamond" // 钻石奖励
)

// RewardFlow 奖励流水记录模型
type RewardFlow struct {
	ID          uint         `gorm:"primary_key" json:"id"`
	UserID      uint         `gorm:"not null;index" json:"user_id"`           // 用户ID
	ItemType    RewardFlowType `gorm:"size:20;not null" json:"item_type"`      // 商品类型
	ItemID      uint         `gorm:"index" json:"item_id"`                     // 商品ID（物品奖励时使用）
	Quantity    int64        `gorm:"not null" json:"quantity"`                 // 获得数量
	Source      string       `gorm:"size:50;not null" json:"source"`           // 奖励来源
	Ctime       time.Time    `gorm:"not null" json:"ctime"`                    // 创建时间
	Utime       time.Time    `gorm:"not null" json:"utime"`                    // 更新时间
}

// TableName 指定表名
func (RewardFlow) TableName() string {
	return "reward_flows"
}

// BeforeCreate 创建前的钩子
func (rf *RewardFlow) BeforeCreate(scope *gorm.Scope) error {
	if rf.Ctime.IsZero() {
		rf.Ctime = time.Now()
	}
	rf.Utime = rf.Ctime
	return nil
}

// BeforeUpdate 更新前的钩子
func (rf *RewardFlow) BeforeUpdate(scope *gorm.Scope) error {
	rf.Utime = time.Now()
	return nil
}