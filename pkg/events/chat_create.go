package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	m "github.com/PRYVT/chats/pkg/models/command"
	"github.com/google/uuid"
)

type ChatCreatedEvent struct {
	Id           string
	Name         string
	UserIds      []uuid.UUID
	CreationDate time.Time
}

func NewChatCreateEvent(cp m.CreateChat) *models.ChangeTrackedEvent {

	userIds := make([]uuid.UUID, 0, len(cp.UserIds))
	for k := range cp.UserIds {
		userIds = append(userIds, k)
	}

	b := UnsafeSerializeAny(ChatCreatedEvent{
		Id:           cp.Id,
		Name:         cp.Name,
		UserIds:      userIds,
		CreationDate: time.Now(),
	})
	return &models.ChangeTrackedEvent{
		Event: models.Event{
			Name: "ChatCreatedEvent",
			Data: b,
		},
	}
}
