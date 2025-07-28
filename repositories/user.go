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
	GetMany(IDs []string) ([]*models.User, error)
	Persist(user *models.User) (*models.User, error)
	GetById(ID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(ID string) error
	UpdateProfilePicture(ID string, url string) error
}

type User struct {
	conn *database.Connection
}

func (repo *User) GetMany(IDs []string) ([]*models.User, error) {
	users := make([]*models.User, len(IDs))
	if len(IDs) == 0 {
		return users, nil
	}
	db := repo.conn.GetConnection()
	if err := db.Find(&users, IDs).Error; err != nil {
		return nil, err
	}
	return users, nil
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

func (repo *User) Update(record *models.User) (*models.User, error) {
	db := repo.conn.GetConnection()

	if err := db.Save(record).Error; err != nil {
		return nil, err
	}

	if err := db.First(record, "id = ?", record.ID).Error; err != nil {
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

func (repo *User) UpdateProfilePicture(ID string, url string) error {
	db := repo.conn.GetConnection()

	if err := db.Model(&models.User{}).
		Where("id = ?", ID).
		Update("profile_picture", url).Error; err != nil {
		return fmt.Errorf("failed to update profile picture for user %s: %w", ID, err)
	}

	return nil
}
