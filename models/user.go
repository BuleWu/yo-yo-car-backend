package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password"`
}
