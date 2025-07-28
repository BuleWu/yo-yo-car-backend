package models

// NewReservation constructor for creating a new reservation instance
func NewReservation(userId string, rideId string, status string) *Reservation {
	return &Reservation{
		UserID: userId,
		RideID: rideId,
		Status: status,
	}
}

type Reservation struct {
	Model
	UserID string `json:"user_id"`
	RideID string `json:"ride_id"`
	Status string `json:"status"`
}
