package applications

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"zavrsni/yo-yo-car/core/utils"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

type RideDTO struct {
	ID            string         `json:"id"`
	StartingPoint string         `json:"starting_point"`
	Destination   string         `json:"destination"`
	StartTime     time.Time      `json:"startTime"`
	EndTime       time.Time      `json:"endTime"`
	Price         float64        `json:"price"`
	Date          time.Time      `json:"date"`
	DriverID      string         `json:"driver_id"`
	Driver        *models.User   `json:"driver"`
	Passengers    []*models.User `json:"passengers"`
	MaxPassengers int            `json:"max_passengers"`
	Finished      bool           `json:"finished"`
}

func NewRideApplication(
	rideRepository repositories.RideRepository,
	userRepository repositories.UserRepository,
	reservationRepository repositories.ReservationRepository,
) *Ride {
	return &Ride{
		rideRepository:        rideRepository,
		userRepository:        userRepository,
		reservationRepository: reservationRepository,
	}
}

type Ride struct {
	Application
	rideRepository        repositories.RideRepository
	userRepository        repositories.UserRepository
	reservationRepository repositories.ReservationRepository
}

func (a *Ride) GetRides() ([]*RideDTO, Exception) {
	rides, err := a.rideRepository.Get()

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	var rideDtos []*RideDTO

	for _, ride := range rides {
		rideDtos = append(rideDtos, ToRideDTO(ride))
	}

	return rideDtos, nil
}

func (a *Ride) GetUserRides(userId string) ([]*RideDTO, Exception) {
	rides, err := a.rideRepository.GetByUserId(userId)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	var rideDtos []*RideDTO

	for _, ride := range rides {
		rideDtos = append(rideDtos, ToRideDTO(ride))
	}

	return rideDtos, nil
}

func (a *Ride) GetRideById(ID string) (*RideDTO, Exception) {
	ride, err := a.rideRepository.GetById(ID)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	rideDto := ToRideDTO(ride)

	return rideDto, nil
}

type CreateRideRequest struct {
	StartingPoint string    `json:"starting_point"`
	Destination   string    `json:"destination"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Price         float64   `json:"price"`
	DriverID      string    `json:"driver_id"`
	MaxPassengers int       `json:"max_passengers"`
	Date          time.Time `json:"date"`
}

func (a *Ride) CreateRide(request *CreateRideRequest) (*RideDTO, Exception) {
	var driver *models.User
	var err error
	if request.DriverID != "" {
		driver, err = a.userRepository.GetById(request.DriverID)
		if err != nil {
			return nil, NewApplicationException(http.StatusNotFound, err)
		}
	}

	price := utils.ToEUR(request.Price)

	ride, err := a.rideRepository.Persist(models.NewRide(request.StartingPoint, request.Destination, request.StartTime, request.EndTime, price, request.DriverID, driver, false, nil, request.MaxPassengers, request.Date))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	rideDto := ToRideDTO(ride)

	return rideDto, nil
}

type UpdateRideRequest struct {
	RideID        string    `json:"-"`
	StartingPoint string    `json:"starting_point"`
	Destination   string    `json:"destination"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Price         float64   `json:"price"`
	DriverID      string    `json:"driver_id"`
	Finished      bool      `json:"finished"`
	PassengerIDs  []string  `json:"passenger_ids"`
	MaxPassengers int       `json:"max_passengers"`
	Date          time.Time `json:"date"`
}

func (a *Ride) UpdateRide(request *UpdateRideRequest) (*RideDTO, Exception) {
	ride, err := a.rideRepository.GetById(request.RideID)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	if request.StartingPoint != "" {
		ride.StartingPoint = request.StartingPoint
	}

	if request.Destination != "" {
		ride.Destination = request.Destination
	}

	if request.Price > 0 {
		ride.Price = utils.ToEUR(request.Price)
	}

	if request.Date.Before(time.Now()) {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("date must be today or in the future"))
	}

	if request.StartTime.Before(time.Now()) {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("start time must be in the future"))
	}

	if !sameDay(request.StartTime, request.Date) {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("start time must be on the same calendar day as the date"))
	}

	if request.EndTime.Before(request.StartTime) {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("end time must be after start time"))
	}

	ride.Date = request.Date
	ride.StartTime = request.StartTime
	ride.EndTime = request.EndTime

	var newDriver *models.User

	if request.DriverID != "" {
		newDriver, err = a.userRepository.GetById(request.DriverID)
		if err != nil {
			return nil, NewApplicationException(http.StatusBadRequest, err)
		}
		ride.DriverID = request.DriverID
		ride.Driver = newDriver
	}

	if len(request.PassengerIDs) > 0 {
		passengers, appErr := a.checkPassengerExistence(request.PassengerIDs)
		if appErr != nil {
			return nil, appErr
		}
		ride.Passengers = passengers
	}

	if request.MaxPassengers > 0 {
		ride.MaxPassengers = request.MaxPassengers
	}

	ride, err = a.rideRepository.Update(ride)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	rideDto := ToRideDTO(ride)

	return rideDto, nil
}

func (a *Ride) DeleteRide(rideID string) Exception {
	ride, err := a.rideRepository.GetById(rideID)
	if err != nil {
		return NewApplicationException(http.StatusNotFound, err)
	}

	if err = a.rideRepository.Delete(ride); err != nil {
		return NewApplicationException(http.StatusInternalServerError, err)
	}

	return nil
}

func (a *Ride) SearchRides(queries []repositories.SearchQuery) ([]*RideDTO, Exception) {
	rides, err := a.rideRepository.Search(queries)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	var rideDtos []*RideDTO

	for _, ride := range rides {
		rideDtos = append(rideDtos, ToRideDTO(ride))
	}

	return rideDtos, nil
}

func (a *Ride) GetRideReservations(rideID string) ([]*models.Reservation, Exception) {
	if rideID == "" {
		return nil, NewApplicationException(http.StatusBadRequest, errors.New("rideID cannot be empty"))
	}

	if _, err := a.rideRepository.GetById(rideID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}
	reservations, err := a.reservationRepository.GetByRideId(rideID)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return reservations, nil
}

func (a *Ride) checkPassengerExistence(PassengerIDs []string) ([]*models.User, *ApplicationException) {
	passengers, err := a.userRepository.GetMany(PassengerIDs)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	if len(passengers) < len(PassengerIDs) {
		var missingPassengers []string
		passengerMap := make(map[string]bool)
		for _, passenger := range passengers {
			passengerMap[fmt.Sprint(passenger.ID)] = true
		}
		for _, passengerID := range PassengerIDs {
			if !passengerMap[passengerID] {
				missingPassengers = append(missingPassengers, passengerID)
			}
		}
		return nil, NewApplicationException(http.StatusUnprocessableEntity, fmt.Errorf("could not find passengers with IDs: %s", strings.Join(missingPassengers, ",")))
	}
	return passengers, nil
}

func ToRideDTO(ride *models.Ride) *RideDTO {
	return &RideDTO{
		ID:            ride.ID,
		StartingPoint: ride.StartingPoint,
		Destination:   ride.Destination,
		StartTime:     ride.StartTime,
		EndTime:       ride.EndTime,
		Price:         ride.Price.Float64(),
		DriverID:      ride.DriverID,
		Driver:        ride.Driver,
		Passengers:    ride.Passengers,
		MaxPassengers: ride.MaxPassengers,
		Finished:      ride.Finished,
		Date:          ride.Date,
	}
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
