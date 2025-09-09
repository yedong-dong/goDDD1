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

type StoreDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int64     `json:"stock"`
	StoreType StoreType `json:"store_type"`
	Status    int       `json:"status"`
	CostType  CostType  `json:"cost_type"`
}

func (s *Store) ToStoreDTO() *StoreDTO {
	return &StoreDTO{
		ID:        s.ID,
		Name:      s.Name,
		Price:     s.Price,
		Stock:     s.Stock,
		StoreType: s.StoreType,
		Status:    s.Status,
		CostType:  s.CostType,
	}
}

func StoreToDTO(store *Store) *StoreDTO {
	if store == nil {
		return nil
	}
	return store.ToStoreDTO()
}
