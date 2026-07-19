package postgres

import (
	"context"

	"github.com/dankobg/juicer/features/chat"
	"github.com/dankobg/juicer/pagination"
	"github.com/dankobg/juicer/postgres"
)

var _ chat.ChatPersistor = (*PostgresChatPersistor)(nil)

type PostgresChatPersistor struct {
	*postgres.PgPersistor
}

func NewPostgresChatPersistor(pst *postgres.PgPersistor) *PostgresChatPersistor {
	return &PostgresChatPersistor{pst}
}

func (pst *PostgresChatPersistor) AddChatMessage(ctx context.Context, channel string, msg chat.ChatMessage) (chat.ChatMessage, error) {
	panic("")
}

func (pst *PostgresChatPersistor) ListChatMessages(ctx context.Context, channel string, filters chat.ChatFilters) (pagination.WithHasMore[chat.ChatMessage], error) {
	panic("")
}
