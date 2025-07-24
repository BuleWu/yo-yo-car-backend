package applications

import "zavrsni/yo-yo-car/repositories"

func NewRatingApplication(
	ratingRepository *repositories.RatingRepository,
	rideRepository *repositories.RideRepository,
	userRepository *repositories.UserRepository,
) *Rating {
	return &Rating{
		ratingRepository: ratingRepository,
		rideRepository:   rideRepository,
		userRepository:   userRepository,
	}
}

type Rating struct {
	Application
	repositories.RatingRepository
	repositories.RideRepository
	repositories.UserRepository
}
