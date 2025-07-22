package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewRideRepository(db *database.Connection) RideRepository {
	return &Ride{
		conn: db,
	}
}

type RideRepository interface {
	Persist(ride *models.Ride) (*models.Ride, error)
}

type Ride struct {
	conn *database.Connection
}

func (repo *Ride) Persist(record *models.Ride) (*models.Ride, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
