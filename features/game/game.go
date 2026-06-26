package game

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/bus"
	"github.com/dankobg/juicer/features/chat"
	"github.com/dankobg/juicer/features/presence"
	"github.com/dankobg/juicer/gameplay"
	"github.com/dankobg/juicer/pagination"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
)

type categoryThreshold struct {
	upperLimit   time.Duration
	timeCategory pb.GameTimeCategory
}

// for now this does not change, so i keep it static
type protoMappingsCache struct {
	gameVariantsProtoToDB       map[pb.GameVariant]int64
	gameTimeKindsProtoToDB      map[pb.GameTimeKind]int64
	gameTimeCategoriesProtoToDB map[pb.GameTimeCategory]int64
	gameResultsProtoToDB        map[pb.GameResult]int64
	gameResultStatusesProtoToDB map[pb.GameResultStatus]int64
	gameStatesProtoToDB         map[pb.GameState]int64

	gameVariantsDBToProto       map[int64]pb.GameVariant
	gameTimeKindsDBToProto      map[int64]pb.GameTimeKind
	gameTimeCategoriesDBToProto map[int64]pb.GameTimeCategory
	gameResultsDBToProto        map[int64]pb.GameResult
	gameResultStatusesDBToProto map[int64]pb.GameResultStatus
	gameStatesDBToProto         map[int64]pb.GameState
}

func newProtoMappingsCache() protoMappingsCache {
	return protoMappingsCache{
		gameVariantsProtoToDB:       make(map[pb.GameVariant]int64),
		gameTimeKindsProtoToDB:      make(map[pb.GameTimeKind]int64),
		gameTimeCategoriesProtoToDB: make(map[pb.GameTimeCategory]int64),
		gameResultsProtoToDB:        make(map[pb.GameResult]int64),
		gameResultStatusesProtoToDB: make(map[pb.GameResultStatus]int64),
		gameStatesProtoToDB:         make(map[pb.GameState]int64),

		gameVariantsDBToProto:       make(map[int64]pb.GameVariant),
		gameTimeKindsDBToProto:      make(map[int64]pb.GameTimeKind),
		gameTimeCategoriesDBToProto: make(map[int64]pb.GameTimeCategory),
		gameResultsDBToProto:        make(map[int64]pb.GameResult),
		gameResultStatusesDBToProto: make(map[int64]pb.GameResultStatus),
		gameStatesDBToProto:         make(map[int64]pb.GameState),
	}
}

type GameService struct {
	presenceSvc        *presence.PresenceService
	chatSvc            *chat.ChatService
	usrRdr             UserReader
	bus                *bus.Bus
	pst                Persistor
	log                *slog.Logger
	categoryThresholds []categoryThreshold
	protoMappingsCache protoMappingsCache
	gamestates         map[int64]*gameplay.GameState
	gameEvent          chan gameplay.GameEvent
}

type UserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

type UserReader interface {
	GetUserInfo(ctx context.Context, userID uuid.UUID) (UserInfo, error)
	GetUsername(ctx context.Context, userID uuid.UUID) (string, error)
}

type Persistor struct {
	Game             GamePersistor
	GameVariant      GameVariantPersistor
	GameTimeCategory GameTimeCategoryPersistor
	GameTimeKind     GameTimeKindPersistor
	ActiveGame       ActiveGamePersistor
	GameResult       GameResultPersistor
	GameResultStatus GameResultStatusPersistor
	GameState        GameStatePersistor
	Rating           RatingPersistor
	Pool             PoolPersistor
}

func NewGameService(
	presenceSvc *presence.PresenceService,
	chatSvc *chat.ChatService,
	bus *bus.Bus,
	usrRdr UserReader,
	pst Persistor,
	l *slog.Logger,
) *GameService {
	return &GameService{
		presenceSvc:        presenceSvc,
		chatSvc:            chatSvc,
		usrRdr:             usrRdr,
		bus:                bus,
		pst:                pst,
		log:                l,
		categoryThresholds: make([]categoryThreshold, 0),
		protoMappingsCache: newProtoMappingsCache(),
		gamestates:         make(map[int64]*gameplay.GameState),
		gameEvent:          make(chan gameplay.GameEvent, 100),
	}
}

func (g *GameService) ListGameVariants(ctx context.Context, request api.ListGameVariantsRequestObject) (pagination.Result[api.GameVariant], error) {
	filters := ListGameVariantsFilters{ListGameVariantsParams: request.Params}
	np := pagination.GetNormalized(request.Params.Page, request.Params.PageSize)
	filters.Page = &np.Page
	filters.PageSize = &np.PageSize

	gameVariants, err := g.pst.GameVariant.ListGameVariants(ctx, filters)
	if err != nil {
		return pagination.Result[api.GameVariant]{}, fmt.Errorf("failed to list gamevariants: %w", err)
	}

	gameVariantsData := make([]api.GameVariant, len(gameVariants.Data))
	for i, gameVariant := range gameVariants.Data {
		gameVariantsData[i] = GameVariantToResponse(gameVariant)
	}

	meta := pagination.NewMeta(np, gameVariants.TotalCount)
	out := pagination.NewRes(gameVariantsData, meta)

	return out, nil
}

func (g *GameService) ListGameTimeCategories(ctx context.Context, request api.ListGameTimeCategoriesRequestObject) (pagination.Result[api.GameTimeCategory], error) {
	filters := ListGameTimeCategoriesFilters{ListGameTimeCategoriesParams: request.Params}
	np := pagination.GetNormalized(request.Params.Page, request.Params.PageSize)
	filters.Page = &np.Page
	filters.PageSize = &np.PageSize

	gameTimeCategories, err := g.pst.GameTimeCategory.ListGameTimeCategories(ctx, filters)
	if err != nil {
		return pagination.Result[api.GameTimeCategory]{}, fmt.Errorf("failed to list gametimecategories: %w", err)
	}

	gameTimeCategoriesData := make([]api.GameTimeCategory, len(gameTimeCategories.Data))
	for i, gameTimeCategory := range gameTimeCategories.Data {
		gameTimeCategoriesData[i] = GameTimeCategoryToResponse(gameTimeCategory)
	}

	meta := pagination.NewMeta(np, gameTimeCategories.TotalCount)
	out := pagination.NewRes(gameTimeCategoriesData, meta)

	return out, nil
}

func (g *GameService) ListGameTimeKinds(ctx context.Context, request api.ListGameTimeKindsRequestObject) (pagination.Result[api.GameTimeKind], error) {
	filters := ListGameTimeKindsFilters{ListGameTimeKindsParams: request.Params}
	np := pagination.GetNormalized(request.Params.Page, request.Params.PageSize)
	filters.Page = &np.Page
	filters.PageSize = &np.PageSize

	gameTimeKinds, err := g.pst.GameTimeKind.ListGameTimeKinds(ctx, filters)
	if err != nil {
		return pagination.Result[api.GameTimeKind]{}, fmt.Errorf("failed to list gametimekinds: %w", err)
	}

	gameTimeKindsData := make([]api.GameTimeKind, len(gameTimeKinds.Data))
	for i, gameTimeKind := range gameTimeKinds.Data {
		gameTimeKindsData[i] = GameTimeKindToResponse(gameTimeKind)
	}

	meta := pagination.NewMeta(np, gameTimeKinds.TotalCount)
	out := pagination.NewRes(gameTimeKindsData, meta)

	return out, nil
}

func (g *GameService) GetGame(ctx context.Context, request api.GetGameRequestObject) (*api.Game, error) {
	filters := GetGameByIDFilters{GetGameParams: request.Params}

	gameDetails, err := g.pst.Game.GetGameByID(ctx, request.ID, filters)
	if err != nil {
		// if errors.Is(err, postgres.ErrGameNotFound) {
		// 	return nil, fmt.Errorf("game not found: %w", err)
		// }
		return nil, fmt.Errorf("failed to get game by id: %w", err)
	}

	out := GameDetailsToResponse(gameDetails)

	return &out, nil
}

func (g *GameService) ListGames(ctx context.Context, request api.ListGamesRequestObject) (pagination.Result[api.Game], error) {
	filters := ListGamesFilters{ListGamesParams: request.Params}
	np := pagination.GetNormalized(request.Params.Page, request.Params.PageSize)
	filters.Page = &np.Page
	filters.PageSize = &np.PageSize

	games, err := g.pst.Game.ListGames(ctx, filters)
	if err != nil {
		return pagination.Result[api.Game]{}, fmt.Errorf("failed to list games: %w", err)
	}

	gamesData := make([]api.Game, len(games.Data))
	for i, gameDetailsData := range games.Data {
		gamesData[i] = GameDetailsToResponse(gameDetailsData)
	}

	meta := pagination.NewMeta(np, games.TotalCount)
	out := pagination.NewRes(gamesData, meta)

	return out, nil
}

func (g *GameService) ListRatings(ctx context.Context, request api.ListRatingsRequestObject) (pagination.Result[api.Rating], error) {
	filters := ListRatingsFilters{ListRatingsParams: request.Params}
	np := pagination.GetNormalized(request.Params.Page, request.Params.PageSize)
	filters.Page = &np.Page
	filters.PageSize = &np.PageSize

	ratings, err := g.pst.Rating.ListRatings(ctx, filters)
	if err != nil {
		return pagination.Result[api.Rating]{}, fmt.Errorf("failed to list ratings: %w", err)
	}

	ratingsData := make([]api.Rating, len(ratings.Data))
	for i, rating := range ratings.Data {
		ratingsData[i] = RatingToResponse(rating)
	}

	meta := pagination.NewMeta(np, ratings.TotalCount)
	out := pagination.NewRes(ratingsData, meta)

	return out, nil
}

var quickGames = []api.QuickGame{
	{Name: "Hyperbullet", ClockSecs: 30, IncrementSecs: 0},
	{Name: "Bullet", ClockSecs: 60, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 180, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 180, IncrementSecs: 1},
	{Name: "Blitz", ClockSecs: 300, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 300, IncrementSecs: 2},
	{Name: "Rapid", ClockSecs: 600, IncrementSecs: 0},
	{Name: "Rapid", ClockSecs: 600, IncrementSecs: 5},
	{Name: "Rapid", ClockSecs: 900, IncrementSecs: 0},
	{Name: "Rapid", ClockSecs: 900, IncrementSecs: 5},
	{Name: "Classical", ClockSecs: 1800, IncrementSecs: 0},
	{Name: "Classical", ClockSecs: 2700, IncrementSecs: 10},
}

func (g *GameService) ListQuickGames(ctx context.Context, request api.ListQuickGamesRequestObject) ([]api.QuickGame, error) {
	return quickGames, nil
}

func (g *GameService) GetGameStats(ctx context.Context, request api.GetGameStatsRequestObject) (api.GameStats, error) {
	return api.GameStats{}, nil
	// 	gameStats, err := a.gamePst.GetGameStatsForUser(ctx, request.UserID, nil)
	// 	if err != nil {
	// 		return api.GetGameStats404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Message: "stats not found for user", Code: 404}}, nil
	// 	}
	// 	resp := api.GetGameStats200JSONResponse(gameStats)
	// 	return resp, nil
}
