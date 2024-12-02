package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	m "github.com/PRYVT/chats/pkg/models/command"
	"github.com/google/uuid"
)

type ChatMessageAddedEvent struct {
	Id           uuid.UUID
	Text         string
	ImageBase64  string
	UserId       uuid.UUID
	CreationDate time.Time
}

func NewChatMessageAddedEvent(chatMessage m.AddChatMessage, userId uuid.UUID) *models.ChangeTrackedEvent {
	b := UnsafeSerializeAny(ChatMessageAddedEvent{
		Id:           uuid.New(),
		Text:         chatMessage.Text,
		ImageBase64:  chatMessage.ImageBase64,
		UserId:       userId,
		CreationDate: time.Now(),
	})
	return &models.ChangeTrackedEvent{
		Event: models.Event{
			Name: "ChatMessageAddedEvent",
			Data: b,
		},
	}
}
