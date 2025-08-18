package applications

import (
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

func NewMessageApplication(
	messageRepository repositories.MessageRepository,
) *Message {
	return &Message{
		messageRepository: messageRepository,
	}
}

type Message struct {
	Application
	messageRepository repositories.MessageRepository
}

type UpdateMessageRequest struct {
	MessageID string `json:"-"`
	Read      bool   `json:"read"`
}

func (a *Message) UpdateMessage(request *UpdateMessageRequest) (*models.Message, Exception) {
	message, err := a.messageRepository.GetById(request.MessageID)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	message.Read = request.Read

	message, err = a.messageRepository.Update(message)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	return message, nil
}
