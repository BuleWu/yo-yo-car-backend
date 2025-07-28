package applications

import (
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

func NewReservationApplication(reservationRepository repositories.ReservationRepository, userRepository repositories.UserRepository, rideRepository repositories.RideRepository) *Reservation {
	return &Reservation{
		reservationRepository: reservationRepository,
		userRepository:        userRepository,
		rideRepository:        rideRepository,
	}
}

type Reservation struct {
	Application
	reservationRepository repositories.ReservationRepository
	userRepository        repositories.UserRepository
	rideRepository        repositories.RideRepository
}

type CreateReservationRequest struct {
	UserID string                   `json:"user_id"`
	RideID string                   `json:"ride_id"`
	Status models.ReservationStatus `json:"status"`
}

func (a *Reservation) CreateReservation(request *CreateReservationRequest) (*models.Reservation, Exception) {
	if _, err := a.userRepository.GetById(request.UserID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	if _, err := a.rideRepository.GetById(request.RideID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	reservation, err := a.reservationRepository.Persist(models.NewReservation(request.UserID, request.RideID, request.Status))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return reservation, nil
}

func (a *Reservation) GetReservationById(id string) (*models.Reservation, Exception) {
	reservation, err := a.reservationRepository.GetById(id)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}
	return reservation, nil
}

func (a *Reservation) GetAllReservations() ([]*models.Reservation, Exception) {
	reservations, err := a.reservationRepository.GetAll()
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return reservations, nil
}

type UpdateReservationRequest struct {
	ReservationID string                   `json:"-"`
	Status        models.ReservationStatus `json:"status"`
}

func (a *Reservation) UpdateReservation(request *UpdateReservationRequest) (*models.Reservation, Exception) {
	reservation, err := a.reservationRepository.GetById(request.ReservationID)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	reservation.Status = request.Status

	updated, err := a.reservationRepository.Update(reservation)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return updated, nil
}

func (a *Reservation) DeleteReservation(id string) Exception {
	reservation, err := a.reservationRepository.GetById(id)
	if err != nil {
		return NewApplicationException(http.StatusNotFound, err)
	}

	if err = a.reservationRepository.Delete(reservation); err != nil {
		return NewApplicationException(http.StatusInternalServerError, err)
	}

	return nil
}
