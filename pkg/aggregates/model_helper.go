package aggregates

import "github.com/PRYVT/chats/pkg/models/query"

func GetChatModelFromAggregate(userAggregate *ChatRoomAggregate) *query.Chat {
	return &query.Chat{
		Id:         userAggregate.AggregateId,
		Name:       userAggregate.Name,
		UserIds:    userAggregate.UserIds,
		ChangeDate: userAggregate.ChangeDate,
		Messages:   userAggregate.Messages,
	}
}
