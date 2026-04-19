package persistence

import (
	"context"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
)

type UserPersistor interface {
	ListUsers(ctx context.Context, filters dbtype.ListUsersFilters) (dbtype.PagedResult[models.User], error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	CreateUser(ctx context.Context, in models.UserSetter) (models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, in models.UserSetter) (models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)

	CreateFriendRequest(ctx context.Context, in models.FriendshipSetter) (dbtype.FriendRequest, error)
	ListFriendRequests(ctx context.Context, filters dbtype.ListFriendRequestsFilters) (dbtype.PagedResult[dbtype.FriendRequest], error)
	AcceptFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	DeclineFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	CancelFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	ListFriends(ctx context.Context, userID uuid.UUID, filters dbtype.ListFriendsFilters) (dbtype.PagedResult[models.User], error)
	GetFriend(ctx context.Context, userID, friendID uuid.UUID) (models.User, error)

	ListFollowings(ctx context.Context, userID uuid.UUID, filters dbtype.ListFollowingsFilters) (dbtype.PagedResult[models.User], error)
	GetFollowing(ctx context.Context, userID, followingID uuid.UUID) (models.User, error)
	FollowUser(ctx context.Context, userID uuid.UUID, in models.FollowingSetter) error
	ListBlockedUsers(ctx context.Context, userID uuid.UUID, filters dbtype.ListBlockedUsersFilters) (dbtype.PagedResult[models.User], error)
	GetBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) (models.User, error)
	BlockUser(ctx context.Context, userID uuid.UUID, in models.BlocklistSetter) error
}

type GameResultPersistor interface {
	ListGameResults(ctx context.Context, filters dbtype.ListGameResultsFilters) (dbtype.PagedResult[models.GameResult], error)
	GetGameResultByID(ctx context.Context, gameResultID int64) (models.GameResult, error)
	CreateGameResult(ctx context.Context, in models.GameResultSetter) (models.GameResult, error)
	UpdateGameResult(ctx context.Context, gameResultID int64, in models.GameResultSetter) (models.GameResult, error)
	DeleteGameResultByID(ctx context.Context, gameResultID int64) (int64, error)
	BulkDeleteGameResults(ctx context.Context, gameResultIDs []int64) error
}

type GameResultStatusPersistor interface {
	ListGameResultStatuses(ctx context.Context, filters dbtype.ListGameResultStatusesFilters) (dbtype.PagedResult[models.GameResultStatus], error)
	GetGameResultStatusByID(ctx context.Context, gameResultStatusID int64) (models.GameResultStatus, error)
	CreateGameResultStatus(ctx context.Context, in models.GameResultStatusSetter) (models.GameResultStatus, error)
	UpdateGameResultStatus(ctx context.Context, gameResultStatusID int64, in models.GameResultStatusSetter) (models.GameResultStatus, error)
	DeleteGameResultStatusByID(ctx context.Context, gameResultStatusID int64) (int64, error)
	BulkDeleteGameResultStatuses(ctx context.Context, gameResultStatusIDs []int64) error
}

type GameStatePersistor interface {
	ListGameStates(ctx context.Context, filters dbtype.ListGameStatesFilters) (dbtype.PagedResult[models.GameState], error)
	GetGameStateByID(ctx context.Context, gameStateID int64) (models.GameState, error)
	CreateGameState(ctx context.Context, in models.GameStateSetter) (models.GameState, error)
	UpdateGameState(ctx context.Context, gameStateID int64, in models.GameStateSetter) (models.GameState, error)
	DeleteGameStateByID(ctx context.Context, gameStateID int64) (int64, error)
	BulkDeleteGameStates(ctx context.Context, gameStateIDs []int64) error
}

type GameTimeCategoryPersistor interface {
	ListGameTimeCategories(ctx context.Context, filters dbtype.ListGameTimeCategoriesFilters) (dbtype.PagedResult[models.GameTimeCategory], error)
	GetGameTimeCategoryByID(ctx context.Context, gameTimeCategoryID int64) (models.GameTimeCategory, error)
	CreateGameTimeCategory(ctx context.Context, in models.GameTimeCategorySetter) (models.GameTimeCategory, error)
	UpdateGameTimeCategory(ctx context.Context, gameTimeCategoryID int64, in models.GameTimeCategorySetter) (models.GameTimeCategory, error)
	DeleteGameTimeCategoryByID(ctx context.Context, gameTimeCategoryID int64) (int64, error)
	BulkDeleteGameTimeCategories(ctx context.Context, gameTimeCategoryIDs []int64) error
}

type GameTimeKindPersistor interface {
	ListGameTimeKinds(ctx context.Context, filters dbtype.ListGameTimeKindsFilters) (dbtype.PagedResult[models.GameTimeKind], error)
	GetGameTimeKindByID(ctx context.Context, gameTimeKindID int64) (models.GameTimeKind, error)
	CreateGameTimeKind(ctx context.Context, in models.GameTimeKindSetter) (models.GameTimeKind, error)
	UpdateGameTimeKind(ctx context.Context, gameTimeKindID int64, in models.GameTimeKindSetter) (models.GameTimeKind, error)
	DeleteGameTimeKindByID(ctx context.Context, gameTimeKindID int64) (int64, error)
	BulkDeleteGameTimeKinds(ctx context.Context, gameTimeKindIDs []int64) error
}

type GameVariantPersistor interface {
	ListGameVariants(ctx context.Context, filters dbtype.ListGameVariantsFilters) (dbtype.PagedResult[models.GameVariant], error)
	GetGameVariantByID(ctx context.Context, gameVariantID int64) (models.GameVariant, error)
	CreateGameVariant(ctx context.Context, in models.GameVariantSetter) (models.GameVariant, error)
	UpdateGameVariant(ctx context.Context, gameVariantID int64, in models.GameVariantSetter) (models.GameVariant, error)
	DeleteGameVariantByID(ctx context.Context, gameVariantID int64) (int64, error)
	BulkDeleteGameVariants(ctx context.Context, gameVariantIDs []int64) error
}

type GamePersistor interface {
	ListGames(ctx context.Context, filters dbtype.ListGamesFilters) (dbtype.PagedResult[models.Game], error)
	GetGameByID(ctx context.Context, gameID int64) (models.Game, error)
	CreateGame(ctx context.Context, in models.GameSetter, inMoves []models.GameMoveSetter) (models.Game, error)
	UpdateGame(ctx context.Context, gameID int64, in models.GameSetter) (models.Game, error)
	DeleteGameByID(ctx context.Context, gameID int64) (int64, error)
	BulkDeleteGames(ctx context.Context, gameIDs []int64) error
	GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *dbtype.GameStatsForUserFilters) (dbtype.GameStats, error)
}

type RatingPersistor interface {
	ListRatings(ctx context.Context, filters dbtype.ListRatingsFilters) (dbtype.PagedResult[models.Rating], error)
	GetRatingByID(ctx context.Context, ratingID int64) (models.Rating, error)
	CreateRating(ctx context.Context, in models.RatingSetter) (models.Rating, error)
	UpdateRating(ctx context.Context, ratingID int64, in models.RatingSetter) (models.Rating, error)
	DeleteRatingByID(ctx context.Context, ratingID int64) (int64, error)
	BulkCreateRatings(ctx context.Context, in []models.RatingSetter) ([]models.Rating, error)
	BulkDeleteRatings(ctx context.Context, ratingIDs []int64) error
}

type PresencePersistor interface {
	SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channel string) ([]string, []string, error)
	ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, []string, error)
	GetPresence(ctx context.Context, userID uuid.UUID) ([]string, error)
	RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, error)
	CountInChannel(ctx context.Context, channel string) (int64, error)
	LastSeen(ctx context.Context, userID uuid.UUID) (int64, error)
	GetChannelsForUser(userID uuid.UUID) ([]string, error)
	GetInChannel(ctx context.Context, channel string) ([]string, error) // later some struct with user data
	// UpdateFollower(ctx context.Context, followee, follower *entity.User, following bool) error
	// UpdateActiveGame(ctx context.Context, activeGameEntry *pb.ActiveGameEntry) ([][][]string, error)

}

type Persistor interface {
	User() UserPersistor
	GameResult() GameResultPersistor
	GameResultStatus() GameResultStatusPersistor
	GameState() GameStatePersistor
	GameTimeCategory() GameTimeCategoryPersistor
	GameTimeKind() GameTimeKindPersistor
	GameVariant() GameVariantPersistor
	Game() GamePersistor
	Rating() RatingPersistor
	Presence() PresencePersistor
}
