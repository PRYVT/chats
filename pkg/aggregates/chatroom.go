package aggregates

import (
	"fmt"
	"slices"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/PRYVT/chats/pkg/events"
	"github.com/PRYVT/chats/pkg/models/command"
	"github.com/PRYVT/chats/pkg/models/common"
	"github.com/google/uuid"
)

type ChatRoomAggregate struct {
	UserIds       []uuid.UUID
	Name          string
	ChangeDate    time.Time
	Messages      []common.ChatMessage
	Events        []models.ChangeTrackedEvent
	aggregateType string
	AggregateId   uuid.UUID
	client        *client.EventSourcingHttpClient
}

func NewChatAggregate(id uuid.UUID) (*ChatRoomAggregate, error) {

	c, err := client.NewEventSourcingHttpClient(client.RetrieveEventSourcingClientUrl())
	if err != nil {
		panic(err)
	}
	iter, err := c.GetEventsOrdered(id.String())
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve events")
	}
	ua := &ChatRoomAggregate{
		client:        c,
		Events:        []models.ChangeTrackedEvent{},
		aggregateType: "chat",
		AggregateId:   id,
		ChangeDate:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
	}

	for {
		ev, ok := iter.Next()
		if !ok {
			break
		}
		changeTrackedEv := models.ChangeTrackedEvent{
			Event: *ev,
			IsNew: false,
		}
		ua.addEvent(&changeTrackedEv)
	}
	return ua, nil
}

func (pa *ChatRoomAggregate) apply_ChatCreatedEvent(e *events.ChatCreatedEvent) {
	pa.UserIds = e.UserIds
	pa.Name = e.Name
	pa.ChangeDate = e.CreationDate
}

func (pa *ChatRoomAggregate) apply_ChatMessageAddedEvent(e *events.ChatMessageAddedEvent) {
	pa.Messages = append(pa.Messages, common.ChatMessage{
		Id:           e.Id,
		Text:         e.Text,
		ImageBase64:  e.ImageBase64,
		UserId:       e.UserId,
		CreationDate: e.CreationDate,
	})
	pa.ChangeDate = e.CreationDate
}

func (ua *ChatRoomAggregate) addEvent(ev *models.ChangeTrackedEvent) {
	switch ev.Name {
	case "ChatCreatedEvent":
		e := events.UnsafeDeserializeAny[events.ChatCreatedEvent](ev.Data)
		ua.apply_ChatCreatedEvent(e)
	case "ChatMessageAddedEvent":
		e := events.UnsafeDeserializeAny[events.ChatMessageAddedEvent](ev.Data)
		ua.apply_ChatMessageAddedEvent(e)
	default:
		panic(fmt.Errorf("NO KNOWN EVENT %v", ev))
	}
	if ev.Version == 0 {
		ev.IsNew = true
	}
	v := len(ua.Events) + 1 //for validation we need to start at 1
	ev.Version = int64(v)
	ev.AggregateType = ua.aggregateType
	ev.AggregateId = ua.AggregateId.String()
	ua.Events = append(ua.Events, *ev)
}

func (ua *ChatRoomAggregate) saveChanges() error {
	return ua.client.AddEvents(ua.AggregateId.String(), ua.Events)
}

func (ua *ChatRoomAggregate) CreateChat(chatCreate command.CreateChat) error {

	if len(ua.Events) != 0 {
		return fmt.Errorf("chat already exists")
	}
	if chatCreate.Id == "" {
		return fmt.Errorf("chat id is empty") // should be aggregateId
	}
	if chatCreate.Name == "" {
		return fmt.Errorf("chat name is empty")
	}
	if len(chatCreate.UserIds) < 2 {
		return fmt.Errorf("chat must have at least 2 users")
	}

	ua.addEvent(events.NewChatCreateEvent(chatCreate))
	err := ua.saveChanges()
	if err != nil {
		return fmt.Errorf("error")
	}
	return nil
}

func (ua *ChatRoomAggregate) AddChatMessage(chatMessage command.AddChatMessage, userId uuid.UUID) error {

	if len(ua.Events) == 0 {
		return fmt.Errorf("chat doesn't exist")
	}
	if chatMessage.Text == "" && chatMessage.ImageBase64 == "" {
		return fmt.Errorf("neither text nor image is provided")
	}

	if userId == uuid.Nil {
		return fmt.Errorf("user id is empty")
	}

	containsUser := slices.ContainsFunc[[]uuid.UUID](ua.UserIds, func(x uuid.UUID) bool {
		return x == userId
	})
	if !containsUser {
		return fmt.Errorf("user is not in the chat")
	}

	ua.addEvent(events.NewChatMessageAddedEvent(chatMessage, userId))
	err := ua.saveChanges()
	if err != nil {
		return fmt.Errorf("error")
	}
	return nil
}
