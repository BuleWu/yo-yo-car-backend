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
	Get() (*[]models.Ride, error)
	GetById(ID string) (*models.Ride, error)
	Persist(ride *models.Ride) (*models.Ride, error)
	Update(ride *models.Ride) (*models.Ride, error)
	Delete(ride *models.Ride) error
}

type Ride struct {
	conn *database.Connection
}

func (repo *Ride) Get() (*[]models.Ride, error) {
	db := repo.conn.GetConnection()
	var record []models.Ride

	if err := db.Preload("Driver").Find(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (repo *Ride) GetById(ID string) (*models.Ride, error) {
	db := repo.conn.GetConnection()
	var record models.Ride

	if err := db.Preload("Driver").First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (repo *Ride) Persist(ride *models.Ride) (*models.Ride, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(ride).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("Driver").First(ride, "id = ?", ride.ID).Error; err != nil {
		return nil, err
	}

	return ride, nil
}

func (repo *Ride) Update(record *models.Ride) (*models.Ride, error) {
	db := repo.conn.GetConnection()

	if err := db.Save(record).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("Driver").First(record, "id = ?", record.ID).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (repo *Ride) Delete(ride *models.Ride) error {
	db := repo.conn.GetConnection()
	if err := db.Delete(ride).Error; err != nil {
		return err
	}
	return nil
}
