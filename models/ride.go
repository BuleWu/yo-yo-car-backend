package models

import "zavrsni/yo-yo-car/core/utils"

// NewRide Ride constructor
func NewRide(startingPoint string, destination string, price utils.EUR, driverId string, driver *User, finished bool, passengers []*User, maxPassengers int) *Ride {
	return &Ride{
		StartingPoint: startingPoint,
		Destination:   destination,
		Price:         price,
		DriverID:      driverId,
		Driver:        driver,
		Finished:      finished,
		Passengers:    passengers,
		MaxPassengers: maxPassengers,
	}
}

type Ride struct {
	Model
	StartingPoint string    `json:"starting_point"`
	Destination   string    `json:"destination"`
	Price         utils.EUR `json:"price"`

	DriverID string `json:"driver_id" gorm:"not null"`
	Driver   *User  `json:"driver" gorm:"foreignKey:DriverID"`

	Passengers []*User `json:"passengers" gorm:"many2many:ride_passengers"`

	Finished      bool `json:"finished"`
	MaxPassengers int  `json:"max_passengers"`
}
