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
	Get() (*[]models.Rating, error)
	GetById(Id string) (*models.Rating, error)
	GetByUserId(userId string) (*[]models.Rating, error)
	Persist(rating *models.Rating) (*models.Rating, error)
	Update(rating *models.Rating) (*models.Rating, error)
	Delete(rating *models.Rating) error
}

type Rating struct {
	conn *database.Connection
}

func (repo *Rating) Get() (*[]models.Rating, error) {
	db := repo.conn.GetConnection()
	var ratings []models.Rating

	if err := db.Find(&ratings).Error; err != nil {
		return nil, err
	}

	return &ratings, nil
}

func (repo *Rating) GetById(id string) (*models.Rating, error) {
	db := repo.conn.GetConnection()
	var record models.Rating

	if err := db.Preload("RatedUser").Preload("Rater").Preload("Ride").First(&record, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (repo *Rating) GetByUserId(userId string) (*[]models.Rating, error) {
	db := repo.conn.GetConnection()
	var ratings []models.Rating

	if err := db.Where("rated_user_id = ?", userId).
		Preload("RatedUser").
		Preload("Rater").
		Preload("Ride").
		Find(&ratings).Error; err != nil {
		return nil, err
	}

	return &ratings, nil
}

func (repo *Rating) Persist(rating *models.Rating) (*models.Rating, error) {
	db := repo.conn.GetConnection()

	if err := db.Create(rating).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("RatedUser").Preload("Rater").Preload("Ride").First(rating, "id = ?", rating.ID).Error; err != nil {
		return nil, err
	}

	return rating, nil
}

func (repo *Rating) Update(rating *models.Rating) (*models.Rating, error) {
	db := repo.conn.GetConnection()

	if err := db.Save(rating).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("RatedUser").Preload("Rater").Preload("Ride").First(rating, "id = ?", rating.ID).Error; err != nil {
		return nil, err
	}

	return rating, nil
}

func (repo *Rating) Delete(rating *models.Rating) error {
	db := repo.conn.GetConnection()

	if err := db.Delete(rating, "id = ?", rating.ID).Error; err != nil {
		return err
	}

	return nil
}
