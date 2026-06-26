package chat

import (
	"context"

	"github.com/dankobg/juicer/pagination"
)

type ChatPersistor interface {
	AddChatMessage(ctx context.Context, channel string, msg ChatMessage) (ChatMessage, error)
	ListChatMessages(ctx context.Context, channel string, filters ChatFilters) (pagination.WithHasMore[ChatMessage], error)
}
