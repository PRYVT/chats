package common

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Id           uuid.UUID `json:"id" binding:"required"`
	Text         string    `json:"text" binding:"required"`
	ImageBase64  string    `json:"image_base64" binding:"required"`
	UserId       uuid.UUID `json:"user_id" binding:"required"`
	CreationDate time.Time `json:"creation_date" binding:"required"`
}
