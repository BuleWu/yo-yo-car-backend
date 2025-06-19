package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

type UserRepository interface {
	Persist(user *models.User) (*models.User, error)
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
