package chat

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dankobg/juicer/pagination"
	"github.com/google/uuid"
)

type ChatService struct {
	lobbyChatPst  ChatPersistor
	gameChatPst   ChatPersistor
	gametvChatPst ChatPersistor
	pmChatPst     ChatPersistor
	log           *slog.Logger
}

func NewChatService(lobbyPst, gamePst, gametvPst, pmPst ChatPersistor, l *slog.Logger) *ChatService {
	return &ChatService{
		lobbyChatPst:  lobbyPst,
		gameChatPst:   gamePst,
		gametvChatPst: gametvPst,
		pmChatPst:     pmPst,
		log:           l,
	}
}

type ChatMessage struct {
	Channel   string
	UserID    uuid.UUID
	Username  string
	MessageID string
	Message   string
	PostedAt  time.Time
}

type ChatFilters struct {
	PageSize *int
	Cursor   *string
}

func (chat *ChatService) persistor(channel string) (ChatPersistor, error) {
	switch {
	case channel == "lobby.chat":
		return chat.lobbyChatPst, nil

	case strings.HasPrefix(channel, "gametv."):
		return chat.gametvChatPst, nil

	case strings.HasPrefix(channel, "game."):
		return chat.gameChatPst, nil

	case strings.HasPrefix(channel, "pm."):
		return chat.pmChatPst, nil

	default:
		return nil, fmt.Errorf("invalid channel: %s", channel)
	}
}

func (chat *ChatService) AddChatMessage(ctx context.Context, channel string, msg ChatMessage) (ChatMessage, error) {
	pst, err := chat.persistor(channel)
	if err != nil {
		return ChatMessage{}, err
	}

	return pst.AddChatMessage(ctx, channel, msg)
}

func (chat *ChatService) ListChatMessages(ctx context.Context, channel string, filters ChatFilters) (pagination.WithHasMore[ChatMessage], error) {
	pst, err := chat.persistor(channel)
	if err != nil {
		return pagination.WithHasMore[ChatMessage]{}, err
	}

	return pst.ListChatMessages(ctx, channel, filters)
}
