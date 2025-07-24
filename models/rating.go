package models

// NewRating constructor for new rating
func NewRating(value int, raterId string, ratedUserId string, rideId string, comment string) *Rating {
	return &Rating{
		Value:       value,
		RaterID:     raterId,
		RatedUserID: ratedUserId,
		RideID:      rideId,
		Comment:     comment,
	}
}

type Rating struct {
	Model
	Value int `json:"value"`

	RaterID string `json:"rater_id" gorm:"not null"`
	Rater   *User  `json:"rater" gorm:"foreignKey:RaterID;references:ID"`

	RatedUserID string `json:"rated_user_id" gorm:"not null"`
	RatedUser   *User  `json:"rated_user" gorm:"foreignKey:RatedUserID;references:ID"`

	RideID string `json:"ride_id" gorm:"not null"`
	Ride   *Ride  `json:"ride" gorm:"foreignKey:RideID;references:ID"`

	Comment string `json:"comment"`
}
