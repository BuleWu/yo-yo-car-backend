package applications

import (
	"fmt"
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

func NewRideApplication(
	rideRepository repositories.RideRepository,
	userRepository repositories.UserRepository,
) *Ride {
	return &Ride{
		rideRepository: rideRepository,
		userRepository: userRepository,
	}
}

type Ride struct {
	Application
	rideRepository repositories.RideRepository
	userRepository repositories.UserRepository
}

func (a *Ride) GetRides() (*[]models.Ride, Exception) {
	rides, err := a.rideRepository.Get()

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return rides, nil
}

func (a *Ride) GetRideById(ID string) (*models.Ride, Exception) {
	ride, err := a.rideRepository.GetById(ID)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return ride, nil
}

type CreateRideRequest struct {
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`
	DriverID      string `json:"driver_id"`
	MaxPassengers int    `json:"max_passengers"`
}

func (a *Ride) CreateRide(request *CreateRideRequest) (*models.Ride, Exception) {
	var driver *models.User
	var err error
	if request.DriverID != "" {
		driver, err = a.userRepository.GetById(request.DriverID)
		if err != nil {
			return nil, NewApplicationException(http.StatusNotFound, err)
		}
	}

	ride, err := a.rideRepository.Persist(models.NewRide(request.StartingPoint, request.Destination, request.DriverID, driver, false, nil, request.MaxPassengers))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return ride, nil
}

type UpdateRideRequest struct {
	RideID        string   `json:"-"`
	StartingPoint string   `json:"starting_point"`
	Destination   string   `json:"destination"`
	DriverID      string   `json:"driver_id"`
	Finished      bool     `json:"finished"`
	PassengerIDs  []string `json:"passenger_ids"`
	MaxPassengers int      `json:"max_passengers"`
}

func (a *Ride) UpdateRide(request *UpdateRideRequest) (*models.Ride, Exception) {
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

	return ride, nil
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
		return nil, NewApplicationException(http.StatusUnprocessableEntity, fmt.Errorf("could not find passengers with IDs: %s", missingPassengers))
	}
	return passengers, nil
}
