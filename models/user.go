package models

// NewUser user constructor
func NewUser(userId string) *User {
	return &User{
		Model: Model{
			ID: userId,
		},
	}
}

type User struct {
	Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"password" gorm:"column:password"`
}
