package repository

import (
	"database/sql"
	"time"

	"github.com/PRYVT/chats/pkg/models/common"
	models "github.com/PRYVT/chats/pkg/models/query"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func (repo *ChatRepository) GetChatById(chatId uuid.UUID) (*models.Chat, error) {
	var chat models.Chat

	stmt, err := repo.db.Prepare(`
		SELECT id, name, change_date
		FROM Chats
		WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var changeDate string
	err = stmt.QueryRow(chatId.String()).Scan(&chat.Id, &chat.Name, &changeDate)
	if err != nil {
		return nil, err
	}

	changeDateT, err := time.Parse(time.RFC3339, changeDate)
	if err != nil {
		log.Warn().Err(err).Msgf("Error while parsing creation date of chat %v", chatId.String())
	}
	chat.ChangeDate = changeDateT

	msgStmt, err := repo.db.Prepare(`
		SELECT id, user_id, text, image_base64, creation_date
		FROM ChatMessages
		WHERE chat_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer msgStmt.Close()

	rows, err := msgStmt.Query(chatId.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message common.ChatMessage
		var creationDate string
		if err := rows.Scan(&message.Id, &message.UserId, &message.Text, &message.ImageBase64, &creationDate); err != nil {
			return nil, err
		}
		creationDateT, err := time.Parse(time.RFC3339, creationDate)
		if err != nil {
			log.Warn().Err(err).Msgf("Error while parsing creation date of message %v", message.Id.String())
		}
		message.CreationDate = creationDateT
		chat.Messages = append(chat.Messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	userStmt, err := repo.db.Prepare(`
		SELECT user_id
		FROM Users
		WHERE chat_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer userStmt.Close()

	userRows, err := userStmt.Query(chatId.String())
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	for userRows.Next() {
		var userId uuid.UUID
		if err := userRows.Scan(&userId); err != nil {
			return nil, err
		}
		chat.UserIds = append(chat.UserIds, userId)
	}

	if err := userRows.Err(); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (repo *ChatRepository) GetAllChats(limit, offset int, userId uuid.UUID) ([]models.ChatReduced, error) {
	stmt, err := repo.db.Prepare(`
		SELECT id, name
		FROM Chats
		WHERE id IN (SELECT chat_id FROM Users WHERE user_id = ?)
		LIMIT ? OFFSET ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId.String(), limit, offset)
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

		userStmt, err := repo.db.Prepare(`
			SELECT user_id
			FROM Users
			WHERE chat_id = ?
		`)
		if err != nil {
			return nil, err
		}
		defer userStmt.Close()

		userRows, err := userStmt.Query(chat.Id.String())
		if err != nil {
			return nil, err
		}
		defer userRows.Close()

		for userRows.Next() {
			var userId uuid.UUID
			if err := userRows.Scan(&userId); err != nil {
				return nil, err
			}
			chat.UserIds = append(chat.UserIds, userId)
		}

		if err := userRows.Err(); err != nil {
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
