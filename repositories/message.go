package repositories

import (
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/models"
)

func NewMessageRepository(db *database.Connection) MessageRepository {
	return &Message{
		conn: db,
	}
}

type MessageRepository interface {
	GetById(ID string) (*models.Message, error)
	Update(message *models.Message) (*models.Message, error)
}

type Message struct {
	conn *database.Connection
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
