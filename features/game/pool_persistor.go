package game

import (
	"context"

	"github.com/google/uuid"
)

type PoolPersistor interface {
	JoinPool(ctx context.Context, userID uuid.UUID, pool Pool) error
	LeavePool(ctx context.Context, userID uuid.UUID) error
	ListPoolPlayers(ctx context.Context, pool Pool) ([]string, error)
	GetPoolUserMeta(ctx context.Context, userID uuid.UUID) (map[string]any, error)
	MatchPair(ctx context.Context, pool Pool) ([2]string, error)
}
