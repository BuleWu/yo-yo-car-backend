package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"zavrsni/yo-yo-car/models"
)

func NewConnection(host, port, user, password, database string) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	fmt.Println("Successfully connected!")

	return &Connection{db: db}, nil
}

type Connection struct {
	db *gorm.DB
}

// GetConnection returns new gorm.DB connection.
func (r *Connection) GetConnection() *gorm.DB {
	return r.db
}
