package applications

import (
	"errors"
	"fmt"
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
	if request.UserID != "" {
		if _, err := a.userRepository.GetById(request.UserID); err != nil {
			return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with ID %s not found: %w", request.UserID, err))
		}
	}

	if request.RideID != "" {
		if _, err := a.rideRepository.GetById(request.RideID); err != nil {
			return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("ride with ID %s not found: %w", request.RideID, err))
		}
	}

	if !isValidReservationStatus(request.Status) {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("not a valid status"))
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

	if isValidReservationStatus(request.Status) {
		reservation.Status = request.Status
	}

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

func isValidReservationStatus(status models.ReservationStatus) bool {
	switch status {
	case models.Pending, models.Confirmed, models.Cancelled, models.Completed:
		return true
	default:
		return false
	}
}
