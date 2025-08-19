package models

func NewChat(user1Id string, user2Id string) *Chat {
	return &Chat{
		User1ID: user1Id,
		User2ID: user2Id,
	}
}

type Chat struct {
	Model
	User1ID string `json:"user_1_id"`
	User2ID string `json:"user_2_id"`
}
