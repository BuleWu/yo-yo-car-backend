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
	Get(ID string) (*models.User, error)
}

type User struct {
	conn *database.Connection
}

func (repo *User) Persist(record *models.User) (*models.User, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// Get user by id
func (repo *User) Get(ID string) (*models.User, error) {
	var record models.User
	db := repo.conn.GetConnection()
	err := db.First(&record, "id = ?", ID).Error
	if err != nil {
		return nil, fmt.Errorf("user with id: %s was not found", ID)
	}
	return &record, nil
}
