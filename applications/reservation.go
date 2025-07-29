package applications

import (
	"fmt"
	"net/http"
	"zavrsni/yo-yo-car/email"
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
	UserID string `json:"user_id"`
	RideID string `json:"ride_id"`
}

func (a *Reservation) CreateReservation(request *CreateReservationRequest) (*models.Reservation, Exception) {
	if request.UserID != "" {
		if _, err := a.userRepository.GetById(request.UserID); err != nil {
			return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with ID %s not found: %w", request.UserID, err))
		}
	}

	var ride *models.Ride
	var err error

	if request.RideID != "" {
		if ride, err = a.rideRepository.GetById(request.RideID); err != nil {
			return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("ride with ID %s not found: %w", request.RideID, err))
		}
	}

	reservation, err := a.reservationRepository.Persist(models.NewReservation(request.UserID, request.RideID, models.Pending))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	var driver *models.User

	if driver, err = a.userRepository.GetById(ride.DriverID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with ID %s not found: %w", ride.DriverID, err))
	}

	passenger, _ := a.userRepository.GetById(request.UserID)
	emailBody := fmt.Sprintf("You received a new reservation from %s %s. Click here to cofirm.", passenger.FirstName, passenger.LastName)

	go func() {
		if err = email.SendEmail(driver.Email, email.ReservationMadeSubject, emailBody); err != nil {
			fmt.Printf("Failed to send email: %v\n", err)
		}
	}()

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

	if updated.Status == models.Confirmed {
		var passenger *models.User
		if passenger, err = a.userRepository.GetById(updated.UserID); err != nil {
			fmt.Printf("user with id %s not found: %v", updated.UserID, err)
		}

		var ride *models.Ride
		if ride, err = a.rideRepository.GetById(updated.RideID); err != nil {
			return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("ride with ID %s not found: %w", updated.RideID, err))
		}

		var driver *models.User

		if driver, err = a.userRepository.GetById(ride.DriverID); err != nil {
			fmt.Printf("user with id %s not found", ride.DriverID)
		}

		formattedDate := ride.StartTime.Format("02 Jan 2006 at 15:04")

		emailBody := fmt.Sprintf("Great news! %s %s has confirmed your reservation on %s from %s to %s.", driver.FirstName, driver.LastName, formattedDate, ride.StartingPoint, ride.Destination)

		go func() {
			if err = email.SendEmail(passenger.Email, email.ReservationConfirmedSubject, emailBody); err != nil {
				fmt.Printf("Failed to send email: %v\n", err)
			}
		}()
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
