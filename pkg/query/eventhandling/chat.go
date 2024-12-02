package eventhandling

import (
	"sync"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/PRYVT/chats/pkg/aggregates"
	"github.com/PRYVT/chats/pkg/query/store/repository"
	"github.com/PRYVT/utils/pkg/interfaces"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ChatEventHandler struct {
	ChatRepo      *repository.ChatRepository
	wsConnections []interfaces.WebsocketConnecter
	mu            sync.Mutex
}

func NewChatEventHandler(ChatRepo *repository.ChatRepository) *ChatEventHandler {
	return &ChatEventHandler{
		ChatRepo:      ChatRepo,
		wsConnections: []interfaces.WebsocketConnecter{},
	}
}

func (eh *ChatEventHandler) AddWebsocketConnection(conn interfaces.WebsocketConnecter) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	eh.wsConnections = append(eh.wsConnections, conn)
}

func removeDisconnectedSockets(slice []interfaces.WebsocketConnecter) []interfaces.WebsocketConnecter {
	output := []interfaces.WebsocketConnecter{}
	for _, element := range slice {
		if element.IsConnected() {
			output = append(output, element)
		}
	}
	return output
}

func (eh *ChatEventHandler) HandleEvent(event models.Event) error {
	if event.AggregateType == "chat" {
		log.Debug().Msg("Handling Chat event")
		ua, err := aggregates.NewChatAggregate(uuid.MustParse(event.AggregateId))
		if err != nil {
			return err
		}
		p := aggregates.GetChatModelFromAggregate(ua)
		err = eh.ChatRepo.AddOrReplaceChat(p)
		if err != nil {
			log.Err(err).Msg("Error while processing user event")
			return err
		}
		for _, conn := range eh.wsConnections {
			if !conn.IsAuthenticated() {
				continue
			}
			err := conn.WriteJSON(p)
			if err != nil {
				log.Warn().Err(err).Msg("Error while writing to websocket connection")
			}
		}
		eh.mu.Lock()
		defer eh.mu.Unlock()
		eh.wsConnections = removeDisconnectedSockets(eh.wsConnections)
		log.Trace().Msgf("Number of active connections: %d", len(eh.wsConnections))
	}
	return nil
}
