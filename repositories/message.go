package repositories

import (
	"fmt"
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewMessageRepository(db *database.Connection) MessageRepository {
	return &Message{
		conn: db,
	}
}

type MessageRepository interface {
	Persist(message *models.Message) (*models.Message, error)
	GetById(ID string) (*models.Message, error)
	Update(message *models.Message) (*models.Message, error)
	GetMessagesByChat(chatID string, limit, offset int) ([]models.Message, error)
}

type Message struct {
	conn *database.Connection
}

func (repo *Message) Persist(record *models.Message) (*models.Message, error) {
	db := repo.conn.GetConnection()

	if err := db.Create(&record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (repo *Message) GetById(ID string) (*models.Message, error) {
	db := repo.conn.GetConnection()
	var record models.Message

	if err := db.First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (repo *Message) Update(record *models.Message) (*models.Message, error) {
	db := repo.conn.GetConnection()

	if err := db.Save(record).Error; err != nil {
		return nil, err
	}
	if err := db.First(record, "id = ?", record.ID).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (repo *Message) GetMessagesByChat(chatID string, limit, offset int) ([]models.Message, error) {
	db := repo.conn.GetConnection()
	var messages []models.Message

	if err := db.Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	fmt.Println("Chat id: ", chatID)

	return messages, nil
}
