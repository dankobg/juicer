package redis

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"time"

	"github.com/dankobg/juicer/features/chat"
	"github.com/dankobg/juicer/pagination"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var _ chat.ChatPersistor = (*RedisChatPersistor)(nil)

type RedisChatPersistor struct {
	rdb *redis.Client
}

func NewRedisChatPersistor(rdb *redis.Client) *RedisChatPersistor {
	return &RedisChatPersistor{rdb: rdb}
}

func (pst *RedisChatPersistor) AddChatMessage(ctx context.Context, channel string, msg chat.ChatMessage) (chat.ChatMessage, error) {
	values := map[string]any{
		"channel":   msg.Channel,
		"user_id":   msg.UserID.String(),
		"username":  msg.Username,
		"message":   msg.Message,
		"posted_at": msg.PostedAt.UnixMilli(),
	}

	res, err := pst.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: channel,
		Values: values,
		MaxLen: 5_000,
	}).Result()
	if err != nil {
		return chat.ChatMessage{}, err
	}

	out := chat.ChatMessage{
		Channel:   channel,
		UserID:    msg.UserID,
		Username:  msg.Username,
		MessageID: res,
		Message:   msg.Message,
		PostedAt:  msg.PostedAt,
	}

	return out, nil
}

func (pst *RedisChatPersistor) ListChatMessages(ctx context.Context, channel string, filters chat.ChatFilters) (pagination.WithHasMore[chat.ChatMessage], error) {
	pageSize := pagination.DefaultPageSize
	if filters.PageSize != nil {
		pageSize = min(max(*filters.PageSize, pagination.MinPageSize), pagination.MaxPageSize)
	}

	// fetching + 1 and discarding to get hasMore (or could do separate xlen...)
	fetchingSize := pageSize + 1

	var (
		msgs []redis.XMessage
		err  error
	)

	var isForward bool

	if isForward {
		start := "-"
		if filters.Cursor != nil {
			start = "(" + *filters.Cursor
		}

		msgs, err = pst.rdb.XRangeN(ctx, channel, start, "+", int64(fetchingSize)).Result()
	} else {
		start := "+"
		if filters.Cursor != nil {
			start = "(" + *filters.Cursor
		}

		msgs, err = pst.rdb.XRevRangeN(ctx, channel, start, "-", int64(fetchingSize)).Result()
	}

	if err != nil {
		return pagination.WithHasMore[chat.ChatMessage]{}, err
	}

	hasMore := len(msgs) > pageSize
	if hasMore {
		msgs = msgs[:pageSize]
	}

	data := make([]chat.ChatMessage, len(msgs))
	for i, m := range msgs {
		chatMsg, err := chatMessageFromStream(m, channel)
		if err != nil {
			return pagination.WithHasMore[chat.ChatMessage]{}, err
		}

		data[i] = chatMsg
	}

	if !isForward {
		slices.Reverse(data)
	}

	out := pagination.NewWithHasMore(data, hasMore)

	return out, nil
}

func chatMessageFromStream(m redis.XMessage, channel string) (chat.ChatMessage, error) {
	userIDRaw, ok1 := m.Values["user_id"].(string)
	username, ok2 := m.Values["username"].(string)
	message, ok3 := m.Values["message"].(string)

	postedAtRaw, ok4 := m.Values["posted_at"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return chat.ChatMessage{}, errors.New("invalid message data")
	}

	postedAtMs, err := strconv.ParseInt(postedAtRaw, 10, 64)
	if err != nil {
		return chat.ChatMessage{}, errors.New("invalid posted_at timestamp")
	}

	postedAt := time.UnixMilli(postedAtMs)

	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return chat.ChatMessage{}, errors.New("invalid user_id format")
	}

	return chat.ChatMessage{
		Channel:   channel,
		UserID:    userID,
		Username:  username,
		MessageID: m.ID,
		Message:   message,
		PostedAt:  postedAt,
	}, nil
}
