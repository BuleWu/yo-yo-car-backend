package models

func NewMessage(content string, chatId string, senderId string, receiverId string, read bool) *Message {
	return &Message{
		Content:    content,
		ChatID:     chatId,
		SenderID:   senderId,
		ReceiverID: receiverId,
		Read:       read,
	}
}

type Message struct {
	Model

	Content string `json:"content"`

	ChatID string `json:"chatID"`

	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`

	/*Sender   User `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver User `json:"receiver" gorm:"foreignKey:ReceiverID"`*/

	Read bool `json:"read"`
}
