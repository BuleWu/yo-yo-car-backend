package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewRatingRepository(db *database.Connection) RatingRepository {
	return &Rating{
		conn: db,
	}
}

type RatingRepository interface {
	Persist(rating *models.Rating) (*models.Rating, error)
}

type Rating struct {
	conn *database.Connection
}

func (repo *Rating) Persist(rating *models.Rating) (*models.Rating, error) {}
