package repository

import (
	"database/sql"

	models "github.com/PRYVT/chats/pkg/models/query"
	"github.com/google/uuid"
)

type ChatRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *ChatRepository {
	if db == nil {
		return nil
	}
	return &ChatRepository{db: db}
}

func (repo *ChatRepository) GetChatById(ChatId uuid.UUID) (*models.Chat, error) {

	return nil, nil
}

func (repo *ChatRepository) GetAllChats(limit, offset int) ([]models.Chat, error) {
	return nil, nil
}

func (repo *ChatRepository) AddOrReplaceChat(Chat *models.Chat) error {
	return nil
}
