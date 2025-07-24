package models

// NewRide Ride constructor
func NewRide(startingPoint string, destination string, driverId string, driver *User, finished bool, passengers []*User, maxPassengers int) *Ride {
	return &Ride{
		StartingPoint: startingPoint,
		Destination:   destination,
		DriverID:      driverId,
		Driver:        driver,
		Finished:      finished,
		Passengers:    passengers,
		MaxPassengers: maxPassengers,
	}
}

type Ride struct {
	Model
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`

	DriverID string `json:"driver_id" gorm:"not null"`
	Driver   *User  `json:"driver" gorm:"foreignKey:DriverID"`

	Passengers []*User `json:"passengers" gorm:"many2many:ride_passengers"`

	Finished      bool `json:"finished"`
	MaxPassengers int  `json:"max_passengers"`
}
