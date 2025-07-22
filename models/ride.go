package models

// NewRide Ride constructor
func NewRide(rideID string) *Ride {
	return &Ride{
		Model: Model{
			ID: rideID,
		},
	}
}

type Ride struct {
	Model
	StartingPoint string `json:"starting_point"`
	Destination   string `json:"destination"`

	DriverID string `json:"driver_id" gorm:"not null"`
	Driver   User   `json:"driver" gorm:"foreignKey:DriverID"`

	Passengers []User `json:"passengers" gorm:"many2many:ride_passengers"`

	Finished      bool `json:"finished"`
	MaxPassengers int  `json:"max_passengers"`
}
