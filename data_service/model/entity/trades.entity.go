package entity

import (
	"time"

	"gorm.io/gorm"
)

type Trades struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Price     int            `json:"price" gorm:"index"`
	CoinID    int            `json:"coin_id" gorm:"index"`
	UserID    int            `json:"user_id" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
