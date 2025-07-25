package models

// NewUser user constructor
func NewUser(firstName string, lastName string, email string, password string, provider string, profilePicture string) *User {
	return &User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Password:       password,
		Provider:       provider,
		ProfilePicture: profilePicture,
	}
}

type User struct {
	Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email" gorm:"uniqueIndex"`
	Password       string `json:"-" gorm:"column:password"`
	Provider       string `json:"provider"`
	ProfilePicture string `json:"profile_picture"`
	Vehicle        string `json:"vehicle"`
}
