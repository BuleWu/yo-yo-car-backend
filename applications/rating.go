package applications

import (
	"net/http"
	"strings"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

var (
	MinRatingValue = 1
	MaxRatingValue = 5
)

func NewRatingApplication(
	ratingRepository repositories.RatingRepository,
	rideRepository repositories.RideRepository,
	userRepository repositories.UserRepository,
) *Rating {
	return &Rating{
		ratingRepository: ratingRepository,
		rideRepository:   rideRepository,
		userRepository:   userRepository,
	}
}

type Rating struct {
	Application
	ratingRepository repositories.RatingRepository
	rideRepository   repositories.RideRepository
	userRepository   repositories.UserRepository
}

func (a *Rating) GetRatings() (*[]models.Rating, Exception) {
	ratings, err := a.ratingRepository.Get()
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return ratings, nil
}

func (a *Rating) GetRatingById(id string) (*models.Rating, Exception) {
	rating, err := a.ratingRepository.GetById(id)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return rating, nil
}

func (a *Rating) GetRatingsByUserId(userId string) (*[]models.Rating, Exception) {
	rating, err := a.ratingRepository.GetByUserId(userId)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return rating, nil
}

type CreateRatingRequest struct {
	Value       int    `json:"value" validate:"required,min=1,max=5"`
	RaterID     string `json:"rater_id" validate:"required,uuid"`
	RatedUserID string `json:"rated_user_id" validate:"required,uuid"`
	RideID      string `json:"ride_id" validate:"required,uuid"`
	Comment     string `json:"comment"`
}

func (a *Rating) CreateRating(request *CreateRatingRequest) (*models.Rating, Exception) {
	rating, err := a.ratingRepository.Persist(models.NewRating(request.Value, request.RaterID, request.RatedUserID, request.RideID, request.Comment))

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return rating, nil
}

type UpdateRatingRequest struct {
	RatingID string `json:"-"`
	Value    int    `json:"value" validate:"required,min=1,max=5"`
	Comment  string `json:"comment"`
}

func (a *Rating) UpdateRating(request *UpdateRatingRequest) (*models.Rating, Exception) {
	rating, err := a.ratingRepository.GetById(request.RatingID)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	if request.Value >= MinRatingValue && request.Value <= MaxRatingValue {
		rating.Value = request.Value
	}

	if strings.TrimSpace(request.Comment) != "" {
		rating.Comment = request.Comment
	}

	rating, err = a.ratingRepository.Update(rating)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return rating, nil
}

func (a *Rating) DeleteRating(ratingId string) Exception {
	rating, err := a.ratingRepository.GetById(ratingId)
	if err != nil {
		return NewApplicationException(http.StatusNotFound, err)
	}

	if err = a.ratingRepository.Delete(rating); err != nil {
		return NewApplicationException(http.StatusInternalServerError, err)
	}

	return nil
}
