package persistence

import (
	"context"
	"time"

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
	AddFriend(ctx context.Context, in models.FriendshipSetter) error
	DeleteFriend(ctx context.Context, userID, friendID uuid.UUID) error
	DeleteFriends(ctx context.Context, userID uuid.UUID, friendIDs []uuid.UUID) error

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

type PoolPersistor interface {
	JoinPool(ctx context.Context, userID uuid.UUID, pool dbtype.Pool) error
	LeavePool(ctx context.Context, userID uuid.UUID) error
	ListPoolPlayers(ctx context.Context, pool dbtype.Pool) ([]string, error)
	GetPoolUserMeta(ctx context.Context, userID uuid.UUID) (map[string]any, error)
	MatchPair(ctx context.Context, pool dbtype.Pool) ([2]string, error)
}

// type GameSeek struct{}
// type GameSeekSetter struct{}

// type GameSeekPersistor interface {
// 	CreateGameSeek(ctx context.Context, in GameSeekSetter) (GameSeek, error)
// 	GetGameSeekByID(ctx context.Context, gameSeekID int64) (GameSeek, error)
// 	DeleteGameSeekByID(ctx context.Context, gameSeekID int64) error

// 	// 	Get(ctx context.Context, id string) (*entity.SoughtGame, error)
// 	// 	GetBySeekerConnID(ctx context.Context, connID string) (*entity.SoughtGame, error)
// 	// 	New(context.Context, *entity.SoughtGame) error
// 	// 	Delete(ctx context.Context, id string) error
// 	// 	ListOpenSeeks(ctx context.Context, receiverID, tourneyID string) ([]*entity.SoughtGame, error)
// 	// 	ListCorrespondenceSeeksForUser(ctx context.Context, userID string) ([]*entity.SoughtGame, error)
// 	// 	ExistsForUser(ctx context.Context, userID string) (bool, error)
// 	// 	CanCreateSeek(ctx context.Context, userID string, gameMode pb.GameMode, receiverID string) (bool, string, error)
// 	// 	DeleteForUser(ctx context.Context, userID string) (*entity.SoughtGame, error)
// 	// 	UpdateForReceiver(ctx context.Context, userID string) (*entity.SoughtGame, error)
// 	// 	DeleteForSeekerConnID(ctx context.Context, connID string) (*entity.SoughtGame, error)
// 	// 	UpdateForReceiverConnID(ctx context.Context, connID string) (*entity.SoughtGame, error)
// 	// 	UserMatchedBy(ctx context.Context, userID, matcher string) (bool, error)
// 	// 	ExpireOld(ctx context.Context) error
// }

type RatingPersistor interface {
	ListRatings(ctx context.Context, filters dbtype.ListRatingsFilters) (dbtype.PagedResult[models.Rating], error)
	GetRatingByID(ctx context.Context, ratingID int64) (models.Rating, error)
	CreateRating(ctx context.Context, in models.RatingSetter) (models.Rating, error)
	UpdateRating(ctx context.Context, ratingID int64, in models.RatingSetter) (models.Rating, error)
	DeleteRatingByID(ctx context.Context, ratingID int64) (int64, error)
	BulkCreateRatings(ctx context.Context, in []models.RatingSetter) ([]models.Rating, error)
	BulkDeleteRatings(ctx context.Context, ratingIDs []int64) error
}

type UserPresenceInfo struct {
	ID       string
	Username string
	Guest    bool
}

type PresenceChannelsDiff struct {
	ConnJoined []string
	ConnLeft   []string
	UserJoined []string
	UserLeft   []string
}

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
	Pool() PoolPersistor
}
