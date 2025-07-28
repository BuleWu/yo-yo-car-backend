package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewReservationRepository(db *database.Connection) ReservationRepository {
	return &Reservation{
		conn: db,
	}
}

type ReservationRepository interface {
	Persist(reservation *models.Reservation) (*models.Reservation, error)
	GetById(id string) (*models.Reservation, error)
	GetAll() ([]*models.Reservation, error)
	Update(reservation *models.Reservation) (*models.Reservation, error)
	Delete(reservation *models.Reservation) error
}

type Reservation struct {
	conn *database.Connection
}

func (repo *Reservation) Persist(reservation *models.Reservation) (*models.Reservation, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(reservation).Error; err != nil {
		return nil, err
	}
	return reservation, nil
}

func (repo *Reservation) GetById(id string) (*models.Reservation, error) {
	db := repo.conn.GetConnection()
	var reservation models.Reservation
	if err := db.First(&reservation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (repo *Reservation) GetAll() ([]*models.Reservation, error) {
	db := repo.conn.GetConnection()
	var reservations []*models.Reservation
	if err := db.Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (repo *Reservation) Update(reservation *models.Reservation) (*models.Reservation, error) {
	db := repo.conn.GetConnection()
	if err := db.Save(reservation).Error; err != nil {
		return nil, err
	}
	return reservation, nil
}

func (repo *Reservation) Delete(reservation *models.Reservation) error {
	db := repo.conn.GetConnection()
	if err := db.Delete(reservation).Error; err != nil {
		return err
	}
	return nil
}
