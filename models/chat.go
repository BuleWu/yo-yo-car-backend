package models

import "time"

func NewChat(user1Id string, user2Id string, rideId string) *Chat {
	return &Chat{
		User1ID:       user1Id,
		User2ID:       user2Id,
		RideID:        rideId,
		LastMessageAt: time.Time{},
	}
}

type Chat struct {
	Model

	User1ID string `json:"user_1_id"`
	User2ID string `json:"user_2_id"`

	User1 *User `json:"user_1" gorm:"foreignKey:User1ID"`
	User2 *User `json:"user_2" gorm:"foreignKey:User2ID"`

	RideID string `json:"ride_id"`
	Ride   *Ride  `json:"ride" gorm:"foreignKey:RideID"`

	LastMessageAt time.Time `json:"last_message_at"`
}
