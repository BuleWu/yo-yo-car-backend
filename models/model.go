package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
