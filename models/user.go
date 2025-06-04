package models

import (
	"github.com/google/uuid"
)

type User struct {
	/*gorm.Model*/
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password"`
}
