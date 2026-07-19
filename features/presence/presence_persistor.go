package presence

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PresencePersistor interface {
	SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channels []string) (PresenceChannelsDiff, error)
	ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) (PresenceChannelsDiff, error)
	RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) error
	UserLastSeen(ctx context.Context, userID uuid.UUID) (time.Time, error)
	ListUsersInChannel(ctx context.Context, channel string) ([]UserPresenceInfo, error)
	ListChannelsForUser(ctx context.Context, userID uuid.UUID) ([]string, error)
	UsersCountInChannel(ctx context.Context, channel string) (int64, error)
	TotalActiveConnsCount(ctx context.Context) (int64, error)
}
