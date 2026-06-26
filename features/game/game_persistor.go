package game

import (
	"context"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/gameplay"
	"github.com/dankobg/juicer/pagination"
	"github.com/google/uuid"
)

type GameResultPersistor interface {
	ListGameResults(ctx context.Context, filters ListGameResultsFilters) (pagination.WithTotal[models.GameResult], error)
	GetGameResultByID(ctx context.Context, gameResultID int64) (models.GameResult, error)
	CreateGameResult(ctx context.Context, in models.GameResultSetter) (models.GameResult, error)
	UpdateGameResult(ctx context.Context, gameResultID int64, in models.GameResultSetter) (models.GameResult, error)
	DeleteGameResultByID(ctx context.Context, gameResultID int64) (int64, error)
	BulkDeleteGameResults(ctx context.Context, gameResultIDs []int64) error
}

type GameResultStatusPersistor interface {
	ListGameResultStatuses(ctx context.Context, filters ListGameResultStatusesFilters) (pagination.WithTotal[models.GameResultStatus], error)
	GetGameResultStatusByID(ctx context.Context, gameResultStatusID int64) (models.GameResultStatus, error)
	CreateGameResultStatus(ctx context.Context, in models.GameResultStatusSetter) (models.GameResultStatus, error)
	UpdateGameResultStatus(ctx context.Context, gameResultStatusID int64, in models.GameResultStatusSetter) (models.GameResultStatus, error)
	DeleteGameResultStatusByID(ctx context.Context, gameResultStatusID int64) (int64, error)
	BulkDeleteGameResultStatuses(ctx context.Context, gameResultStatusIDs []int64) error
}

type GameStatePersistor interface {
	ListGameStates(ctx context.Context, filters ListGameStatesFilters) (pagination.WithTotal[models.GameState], error)
	GetGameStateByID(ctx context.Context, gameStateID int64) (models.GameState, error)
	CreateGameState(ctx context.Context, in models.GameStateSetter) (models.GameState, error)
	UpdateGameState(ctx context.Context, gameStateID int64, in models.GameStateSetter) (models.GameState, error)
	DeleteGameStateByID(ctx context.Context, gameStateID int64) (int64, error)
	BulkDeleteGameStates(ctx context.Context, gameStateIDs []int64) error
}

type GameTimeCategoryPersistor interface {
	ListGameTimeCategories(ctx context.Context, filters ListGameTimeCategoriesFilters) (pagination.WithTotal[models.GameTimeCategory], error)
	GetGameTimeCategoryByID(ctx context.Context, gameTimeCategoryID int64) (models.GameTimeCategory, error)
	CreateGameTimeCategory(ctx context.Context, in models.GameTimeCategorySetter) (models.GameTimeCategory, error)
	UpdateGameTimeCategory(ctx context.Context, gameTimeCategoryID int64, in models.GameTimeCategorySetter) (models.GameTimeCategory, error)
	DeleteGameTimeCategoryByID(ctx context.Context, gameTimeCategoryID int64) (int64, error)
	BulkDeleteGameTimeCategories(ctx context.Context, gameTimeCategoryIDs []int64) error
}

type GameTimeKindPersistor interface {
	ListGameTimeKinds(ctx context.Context, filters ListGameTimeKindsFilters) (pagination.WithTotal[models.GameTimeKind], error)
	GetGameTimeKindByID(ctx context.Context, gameTimeKindID int64) (models.GameTimeKind, error)
	CreateGameTimeKind(ctx context.Context, in models.GameTimeKindSetter) (models.GameTimeKind, error)
	UpdateGameTimeKind(ctx context.Context, gameTimeKindID int64, in models.GameTimeKindSetter) (models.GameTimeKind, error)
	DeleteGameTimeKindByID(ctx context.Context, gameTimeKindID int64) (int64, error)
	BulkDeleteGameTimeKinds(ctx context.Context, gameTimeKindIDs []int64) error
}

type GameVariantPersistor interface {
	ListGameVariants(ctx context.Context, filters ListGameVariantsFilters) (pagination.WithTotal[models.GameVariant], error)
	GetGameVariantByID(ctx context.Context, gameVariantID int64) (models.GameVariant, error)
	CreateGameVariant(ctx context.Context, in models.GameVariantSetter) (models.GameVariant, error)
	UpdateGameVariant(ctx context.Context, gameVariantID int64, in models.GameVariantSetter) (models.GameVariant, error)
	DeleteGameVariantByID(ctx context.Context, gameVariantID int64) (int64, error)
	BulkDeleteGameVariants(ctx context.Context, gameVariantIDs []int64) error
}

type GamePersistor interface {
	GetGameByID(ctx context.Context, gameID int64, filters GetGameByIDFilters) (GameDetails, error)
	ListGames(ctx context.Context, filters ListGamesFilters) (pagination.WithTotal[GameDetails], error)
	CreateGame(ctx context.Context, in models.GameSetter, inMoves []models.GameMoveSetter, inHashes []models.GameHistoryHashSetter) (models.Game, error)
	UpdateGame(ctx context.Context, gameID int64, in models.GameSetter, newMove *models.GameMoveSetter, newHash *models.GameHistoryHashSetter) (models.Game, error)
	DeleteGameByID(ctx context.Context, gameID int64) (int64, error)
	BulkDeleteGames(ctx context.Context, gameIDs []int64) error
	GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *GameStatsForUserFilters) (GameStats, error)
}

type ActiveGamePersistor interface {
	GetActiveGameByID(ctx context.Context, gameID int64, filters GetActiveGameFilters) (GameDetails, error)
	ListActiveGames(ctx context.Context, filters ListActiveGameFilters) (pagination.WithTotal[GameDetails], error)
	ListUserActiveGames(ctx context.Context, userID uuid.UUID, filters ListActiveGameFilters) (pagination.WithTotal[GameDetails], error)
	IsGameActive(ctx context.Context, gameID int64) (bool, error)
	IsUserInActiveGame(ctx context.Context, userID uuid.UUID, gameID int64) (bool, error)
	ListActiveGameUsers(ctx context.Context, gameID int64) ([2]uuid.UUID, error)
	CreateActiveGame(ctx context.Context, gs *gameplay.GameState) error
	UpdateActiveGame(ctx context.Context, gameID int64, in models.GameSetter) error
	DeleteActiveGameByID(ctx context.Context, gameID int64) error
}

// ########################

type RatingPersistor interface {
	ListRatings(ctx context.Context, filters ListRatingsFilters) (pagination.WithTotal[models.Rating], error)
	GetRatingByID(ctx context.Context, ratingID int64) (models.Rating, error)
	CreateRating(ctx context.Context, in models.RatingSetter) (models.Rating, error)
	UpdateRating(ctx context.Context, ratingID int64, in models.RatingSetter) (models.Rating, error)
	DeleteRatingByID(ctx context.Context, ratingID int64) (int64, error)
	BulkCreateRatings(ctx context.Context, in []models.RatingSetter) ([]models.Rating, error)
	BulkDeleteRatings(ctx context.Context, ratingIDs []int64) error
}
