package entity

import (
	"time"

	"gorm.io/gorm"
)

type Trades struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Price     int            `json:"name"`
	UserID    int            `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
