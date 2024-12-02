package aggregates

import "github.com/PRYVT/chats/pkg/models/query"

func GetChatModelFromAggregate(userAggregate *ChatRoomAggregate) *query.Chat {
	return &query.Chat{}
}
