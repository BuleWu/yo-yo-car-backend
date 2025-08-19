package applications

import (
	"fmt"
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/pusher"
	"zavrsni/yo-yo-car/repositories"
)

func NewChatApplication(chatRepository repositories.ChatRepository, userRepository repositories.UserRepository, messageRepository repositories.MessageRepository, pusherService *pusher.PusherService) *Chat {
	return &Chat{
		chatRepository:    chatRepository,
		userRepository:    userRepository,
		messageRepository: messageRepository,
		pusherService:     pusherService,
	}
}

type Chat struct {
	Application
	chatRepository    repositories.ChatRepository
	userRepository    repositories.UserRepository
	messageRepository repositories.MessageRepository
	pusherService     *pusher.PusherService
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
	if _, err := a.chatRepository.GetById(chatID); err != nil {
		return NewApplicationException(http.StatusNotFound, fmt.Errorf("chat with id %s doesn't exist", chatID))
	}

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

func (a *Chat) SendMessage(request *SendMessageRequest) (*models.Message, Exception) {
	if _, err := a.chatRepository.GetById(request.ChatID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("chat with id %s doesn't exist", request.ChatID))
	}

	if _, err := a.userRepository.GetById(request.SenderID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with id %s doesn't exist", request.SenderID))
	}

	if _, err := a.userRepository.GetById(request.ReceiverID); err != nil {
		return nil, NewApplicationException(http.StatusNotFound, fmt.Errorf("user with id %s doesn't exist", request.ReceiverID))
	}

	message, err := a.messageRepository.Persist(models.NewMessage(request.Content, request.ChatID, request.SenderID, request.ReceiverID, false))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, fmt.Errorf("failed to create message"))
	}

	data := map[string]interface{}{
		"message_id":  message.ID,
		"message":     message.Content,
		"sender_id":   message.SenderID,
		"receiver_id": message.ReceiverID,
		"timestamp":   message.CreatedAt,
	}
	err = a.pusherService.Client.Trigger("chat-"+request.ChatID, "send-message", data)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, fmt.Errorf("pusher trigger failed for chatID=%s, senderID=%s, receiverID=%s: %v", request.ChatID, request.SenderID, request.ReceiverID, err))
	}

	return message, nil
}
