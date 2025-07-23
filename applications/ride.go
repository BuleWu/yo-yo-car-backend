package applications

import (
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
	RideID        string `json:"-"`
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`
	DriverID      string `json:"driver_id"`
	Finished      bool   `json:"finished"`
	MaxPassengers int    `json:"max_passengers"`
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
