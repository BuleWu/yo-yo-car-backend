package applications

import (
	"fmt"
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
)

func NewChatApplication(chatRepository repositories.ChatRepository, userRepository repositories.UserRepository, messageRepository repositories.MessageRepository) *Chat {
	return &Chat{
		chatRepository:    chatRepository,
		userRepository:    userRepository,
		messageRepository: messageRepository,
	}
}

type Chat struct {
	Application
	chatRepository    repositories.ChatRepository
	userRepository    repositories.UserRepository
	messageRepository repositories.MessageRepository
}

type CreateChatRequest struct {
	User1ID string `json:"user_1_id"`
	User2ID string `json:"user_2_id"`
}

func (a *Chat) CreateChat(request *CreateChatRequest) (*models.Chat, Exception) {
	var err error

	if _, err = a.userRepository.GetById(request.User1ID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with id %s not found", request.User1ID))
	}

	if _, err = a.userRepository.GetById(request.User2ID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with id %s not found", request.User2ID))
	}

	chat, err := a.chatRepository.Persist(models.NewChat(request.User1ID, request.User2ID))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return chat, nil
}

func (a *Chat) GetUserChats(userID string) ([]models.Chat, Exception) {
	if _, err := a.userRepository.GetById(userID); err != nil {
		return nil, NewApplicationException(
			http.StatusNotFound,
			fmt.Errorf("user with id %s not found", userID),
		)
	}

	chats, err := a.chatRepository.GetUserChats(userID)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return chats, nil
}

func (a *Chat) GetChatMessages(chatID string, limit, offset int) ([]models.Message, Exception) {
	messages, err := a.messageRepository.GetMessagesByChat(chatID, limit, offset)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return messages, nil
}

func (a *Chat) DeleteChat(chatID string) Exception {
	if err := a.chatRepository.Delete(chatID); err != nil {
		return NewApplicationException(http.StatusInternalServerError, err)
	}
	return nil
}

type SendMessageRequest struct {
	ChatID     string `json:"chat_id"`
	Content    string `json:"content"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
}
