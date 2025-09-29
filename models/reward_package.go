package models

import (
	"time"
)

type RewardPackage struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (RewardPackage) TableName() string {
	return "reward_package"
}

// 奖励类型常量
const (
	ItemTypeGoods    uint = 0 // 商品货物
	ItemTypeCurrency uint = 1 // 货币
	// 预留扩展空间
	// ItemTypeEquipment uint = 2 // 装备
	// ItemTypeBadge     uint = 3 // 徽章
	// ItemTypeVIP       uint = 4 // VIP特权
)

type RewardPackageItem struct {
	ID        uint      `gorm:"primaryKey"`
	PackageID uint      `gorm:"column:package_id"`
	ItemType  uint      `gorm:"column:item_type" json:"item_type"` // 0:商品货物, 1:货币, 2+:预留扩展
	ItemID    uint      `gorm:"column:item_id" json:"item_id"`
	Num       uint      `gorm:"column:num" json:"num"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (RewardPackageItem) TableName() string {
	return "reward_package_item"
}

// GetItemTypeName 获取奖励类型名称
func (item *RewardPackageItem) GetItemTypeName() string {
	switch item.ItemType {
	case ItemTypeGoods:
		return "商品"
	case ItemTypeCurrency:
		return "货币"
	// 预留扩展空间
	// case ItemTypeEquipment:
	//     return "装备"
	// case ItemTypeBadge:
	//     return "徽章"
	// case ItemTypeVIP:
	//     return "VIP特权"
	default:
		return "未知类型"
	}
}

type RewardRecord struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"column:user_id"`
	PackageID uint      `gorm:"column:package_id"`
	Source    string    `gorm:"column:source"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (RewardRecord) TableName() string {
	return "reward_record"
}
