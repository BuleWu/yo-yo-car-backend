package models

import (
	"time"
	"zavrsni/yo-yo-car/core/utils"
)

// NewRide Ride constructor
func NewRide(startingPoint string, destination string, startTime time.Time, endTime time.Time, price utils.EUR, driverId string, driver *User, finished bool, passengers []*User, maxPassengers int, date time.Time) *Ride {
	return &Ride{
		StartingPoint: startingPoint,
		Destination:   destination,
		StartTime:     startTime,
		EndTime:       endTime,
		Price:         price,
		DriverID:      driverId,
		Driver:        driver,
		Finished:      finished,
		Passengers:    passengers,
		MaxPassengers: maxPassengers,
		Date:          date,
	}
}

type Ride struct {
	Model
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`

	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`

	Price utils.EUR `json:"price"`

	DriverID string `json:"driver_id" gorm:"not null"`
	Driver   *User  `json:"driver" gorm:"foreignKey:DriverID"`

	Passengers []*User `json:"passengers" gorm:"many2many:ride_passengers"`

	Finished      bool `json:"finished"`
	MaxPassengers int  `json:"max_passengers"`

	Date time.Time `json:"date"`
}
