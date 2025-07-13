package repositories

import (
	"fmt"
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewUserRepository(db *database.Connection) UserRepository {
	return &User{
		conn: db,
	}
}

type UserRepository interface {
	Persist(user *models.User) (*models.User, error)
	GetById(ID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Delete(ID string) error
}

type User struct {
	conn *database.Connection
}

// GetById Get user by id
func (repo *User) GetById(ID string) (*models.User, error) {
	var record models.User
	db := repo.conn.GetConnection()
	err := db.First(&record, "id = ?", ID).Error
	if err != nil {
		return nil, fmt.Errorf("user with id: %s was not found", ID)
	}
	return &record, nil
}

// GetByEmail Get user by email
func (repo *User) GetByEmail(email string) (*models.User, error) {
	var record models.User
	db := repo.conn.GetConnection()
	err := db.First(&record, "email = ?", email).Error
	if err != nil {
		return nil, fmt.Errorf("user with email %s was not found", email)
	}
	return &record, nil
}

func (repo *User) Persist(record *models.User) (*models.User, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (repo *User) Delete(ID string) error {
	var record models.User
	db := repo.conn.GetConnection()
	err := db.Delete(&record, "id = ?", ID).Error
	if err != nil {
		return fmt.Errorf("user with id: %s couldn't be deleted", ID)
	}
	return nil
}
