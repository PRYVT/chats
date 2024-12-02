package repository

import (
	"database/sql"
	"time"

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

func (repo *ChatRepository) GetAllChats(limit, offset int) ([]models.ChatReduced, error) {
	stmt, err := repo.db.Prepare(`
		SELECT id, name
		FROM Chats
		LIMIT ? OFFSET ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatReduced
	for rows.Next() {
		var chat models.ChatReduced
		if err := rows.Scan(&chat.Id, &chat.Name); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (repo *ChatRepository) AddOrReplaceChat(chat *models.Chat) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO Chats (id, name, change_date) VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET 
		name=excluded.name, 
		change_date=excluded.change_date
	`, chat.Id.String(), chat.Name, chat.ChangeDate.Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, message := range chat.Messages {
		_, err = tx.Exec(`
			INSERT INTO ChatMessages (id, chat_id, user_id, text, image_base64, creation_date) VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET 
			chat_id=excluded.chat_id, 
			user_id=excluded.user_id, 
			text=excluded.text, 
			image_base64=excluded.image_base64, 
			creation_date=excluded.creation_date 
		`, message.Id.String(), chat.Id.String(), message.UserId.String(), message.Text, message.ImageBase64, message.CreationDate.Format(time.RFC3339))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, userId := range chat.UserIds {
		_, err = tx.Exec(`
			INSERT INTO Users (user_id, chat_id) VALUES (?, ?)
			ON CONFLICT(user_id, chat_id) DO NOTHING
		`, userId.String(), chat.Id.String())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
