package store

import (
	"context"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/db/gen/test/public/model"
	"github.com/google/uuid"
)

type UserStore interface {
	Create(ctx context.Context, in dto.UserChangeset) (model.User, error)
	Update(ctx context.Context, in dto.UserChangeset) (model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, userID uuid.UUID) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
}

type GameResultStore interface {
	Create(ctx context.Context, in dto.GameResultChangeset) (model.GameResult, error)
	Update(ctx context.Context, in dto.GameResultChangeset) (model.GameResult, error)
	Delete(ctx context.Context, gameResultID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameResultID uuid.UUID) (model.GameResult, error)
	List(ctx context.Context) ([]model.GameResult, error)
}

type GameResultStatusStore interface {
	Create(ctx context.Context, in dto.GameResultStatusChangeset) (model.GameResultStatus, error)
	Update(ctx context.Context, in dto.GameResultStatusChangeset) (model.GameResultStatus, error)
	Delete(ctx context.Context, gameResultStatusID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameResultStatusID uuid.UUID) (model.GameResultStatus, error)
	List(ctx context.Context) ([]model.GameResultStatus, error)
}

type GameStateStore interface {
	Create(ctx context.Context, in dto.GameStateChangeset) (model.GameState, error)
	Update(ctx context.Context, in dto.GameStateChangeset) (model.GameState, error)
	Delete(ctx context.Context, gameStateID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameStateID uuid.UUID) (model.GameState, error)
	List(ctx context.Context) ([]model.GameState, error)
}

type GameTimeCategoryStore interface {
	Create(ctx context.Context, in dto.GameTimeCategoryChangeset) (model.GameTimeCategory, error)
	Update(ctx context.Context, in dto.GameTimeCategoryChangeset) (model.GameTimeCategory, error)
	Delete(ctx context.Context, gameTimeCategoryID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameTimeCategoryID uuid.UUID) (model.GameTimeCategory, error)
	List(ctx context.Context) ([]model.GameTimeCategory, error)
}

type GameTimeKindStore interface {
	Create(ctx context.Context, in dto.GameTimeKindChangeset) (model.GameTimeKind, error)
	Update(ctx context.Context, in dto.GameTimeKindChangeset) (model.GameTimeKind, error)
	Delete(ctx context.Context, gameTimeKindID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameTimeKindID uuid.UUID) (model.GameTimeKind, error)
	List(ctx context.Context) ([]model.GameTimeKind, error)
}

type GameVariantStore interface {
	Create(ctx context.Context, in dto.GameVariantChangeset) (model.GameVariant, error)
	Update(ctx context.Context, in dto.GameVariantChangeset) (model.GameVariant, error)
	Delete(ctx context.Context, gameVariantID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameVariantID uuid.UUID) (model.GameVariant, error)
	List(ctx context.Context) ([]model.GameVariant, error)
}

type GameStore interface {
	Create(ctx context.Context, in dto.GameChangeset, inMoves []dto.GameMoveChangeset) (model.Game, error)
	Update(ctx context.Context, gameID uuid.UUID, in dto.GameChangeset) (model.Game, error)
	Delete(ctx context.Context, gameID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, gameID uuid.UUID) (model.Game, error)
	List(ctx context.Context) ([]model.Game, error)
	GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *GameStatsForUserFilters) (api.GameStats, error)
}

type PresenceStore interface {
	SetActiveGame(ctx context.Context, in model.Game) (model.Game, error)
	GetActiveGame(ctx context.Context, gameID uuid.UUID) (model.Game, error)
	RemoveActiveGame(ctx context.Context, gameID uuid.UUID) error
	SetPlayerGameID(ctx context.Context, clientID uuid.UUID, gameID uuid.UUID) error
	GetPlayerGameID(ctx context.Context, clientID uuid.UUID) (uuid.UUID, error)
	DelPlayerGameID(ctx context.Context, clientID uuid.UUID) error
}

type RatingStore interface {
	BatchCreate(ctx context.Context, in []dto.RatingChangeset) ([]model.Rating, error)
	Create(ctx context.Context, in dto.RatingChangeset) (model.Rating, error)
	Update(ctx context.Context, in dto.RatingChangeset) (model.Rating, error)
	Delete(ctx context.Context, ratingID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, ratingID uuid.UUID) (model.Rating, error)
	List(ctx context.Context) ([]model.Rating, error)
}

type Store interface {
	User() UserStore
	GameResult() GameResultStore
	GameResultStatus() GameResultStatusStore
	GameState() GameStateStore
	GameTimeCategory() GameTimeCategoryStore
	GameTimeKind() GameTimeKindStore
	GameVariant() GameVariantStore
	Game() GameStore
	Presence() PresenceStore
	Rating() RatingStore
}
