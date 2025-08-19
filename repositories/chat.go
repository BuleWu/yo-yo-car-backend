package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewChatRepository(db *database.Connection) ChatRepository {
	return &Chat{
		conn: db,
	}
}

// ChatRepository defines all necessary chat operations
type ChatRepository interface {
	Persist(chat *models.Chat) (*models.Chat, error)
	GetById(ID string) (*models.Chat, error)
	Update(chat *models.Chat) (*models.Chat, error)
	Delete(ID string) error
	GetUserChats(userID string) ([]models.Chat, error)
	GetChatBetweenUsers(user1ID, user2ID string) (*models.Chat, error)
}

// Chat repository implementation
type Chat struct {
	conn *database.Connection
}

// Persist creates a new chat
func (repo *Chat) Persist(chat *models.Chat) (*models.Chat, error) {
	db := repo.conn.GetConnection()
	if err := db.Create(&chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

// GetById fetches a chat by ID
func (repo *Chat) GetById(ID string) (*models.Chat, error) {
	db := repo.conn.GetConnection()
	var chat models.Chat
	if err := db.First(&chat, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

// Update modifies an existing chat
func (repo *Chat) Update(chat *models.Chat) (*models.Chat, error) {
	db := repo.conn.GetConnection()
	if err := db.Save(chat).Error; err != nil {
		return nil, err
	}
	if err := db.First(chat, "id = ?", chat.ID).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

// Delete removes a chat by ID
func (repo *Chat) Delete(ID string) error {
	db := repo.conn.GetConnection()
	return db.Delete(&models.Chat{}, "id = ?", ID).Error
}

// GetUserChats returns all chats that involve a user
func (repo *Chat) GetUserChats(userID string) ([]models.Chat, error) {
	db := repo.conn.GetConnection()
	var chats []models.Chat
	if err := db.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// GetChatBetweenUsers returns a 1-on-1 chat between two users (or nil if it doesn't exist)
func (repo *Chat) GetChatBetweenUsers(user1ID, user2ID string) (*models.Chat, error) {
	db := repo.conn.GetConnection()
	var chat models.Chat
	if err := db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}
