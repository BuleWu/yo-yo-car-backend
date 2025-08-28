package models

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "pending"
	ReservationConfirmed ReservationStatus = "confirmed"
	ReservationCancelled ReservationStatus = "cancelled"
)

// NewReservation constructor for creating a new reservation instance
func NewReservation(userId string, rideId string, status ReservationStatus) *Reservation {
	return &Reservation{
		UserID: userId,
		RideID: rideId,
		Status: status,
	}
}

type Reservation struct {
	Model
	UserID string            `json:"user_id"`
	RideID string            `json:"ride_id"`
	Status ReservationStatus `json:"status"`
}
