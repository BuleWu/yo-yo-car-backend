package applications

import (
	"github.com/google/uuid"
	"net/http"
	"zavrsni/yo-yo-car/models"
)

func NewRideApplication(
	rideRepository repositories.RideRepository,
) *Ride {
	return &Ride{
		rideRepository: rideRepository,
	}
}

type Ride struct {
	Application
	rideRepository repositories.RideRepository
}

type CreateRideRequest struct {
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`
	DriverID      string `json:"driver_id"`
	MaxPassengers int    `json:"max_passengers"`
}

func (a *Ride) CreateRide(request *CreateRideRequest) (*models.Ride, Exception) {
	ride := models.NewRide(uuid.NewString())
	ride.StartingPoint = request.StartingPoint
	ride.Destination = request.Destination
	ride.DriverID = request.DriverID
	ride.Passengers = []models.User{}
	ride.Finished = false
	ride.MaxPassengers = request.MaxPassengers
	ride, err := a.rideRepository.Persist(ride)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return ride, nil
}
