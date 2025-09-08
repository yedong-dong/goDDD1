package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type CostType string

const (
	CostTypeDiamond CostType = "diamond"
	CostTypeCoin    CostType = "coin"
)

type StoreType string

const (
	StoreTypeGood StoreType = "good"
	StoreTypeGift StoreType = "gift"
)

type Store struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `gorm:"size:50;not null;unique" json:"name"`
	Price     int64      `gorm:"not null" json:"price"`
	Stock     int64      `gorm:"not null" json:"stock"`
	StoreType StoreType  `gorm:"size:20;not null" json:"store_type"`
	Status    int        `gorm:"not null;default:1" json:"status"`
	CostType  CostType   `gorm:"size:20;not null" json:"cost_type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (Store) TableName() string {
	return "stores"
}

func (s *Store) BeforeCreate(scope *gorm.Scope) error {
	return nil
}
