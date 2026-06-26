package game

import (
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/gameplay"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob/types"
)

type ListGameResultStatusesFilters struct {
	// Page Page number (1-based)
	Page *api.PaginationPage `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *api.PaginationPageSize `form:"page_size,omitempty" json:"page_size,omitempty"`

	// ID Filter game result statuses by ids
	ID *[]int64 `form:"id,omitempty" json:"id,omitempty"`

	// Name Filter game result statuses by name (partial match)
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Sort Sort by fields (add prefix `-` for descending e.g. -created_at)
	Sort *[]string `form:"sort,omitempty" json:"sort,omitempty"`
}

type ListGameResultsFilters struct {
	// Page Page number (1-based)
	Page *api.PaginationPage `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *api.PaginationPageSize `form:"page_size,omitempty" json:"page_size,omitempty"`

	// ID Filter game results by ids
	ID *[]int64 `form:"id,omitempty" json:"id,omitempty"`

	// Name Filter game results by name (partial match)
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Sort Sort by fields (add prefix `-` for descending e.g. -created_at)
	Sort *[]string `form:"sort,omitempty" json:"sort,omitempty"`
}

type ListGameStatesFilters struct {
	// Page Page number (1-based)
	Page *api.PaginationPage `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *api.PaginationPageSize `form:"page_size,omitempty" json:"page_size,omitempty"`

	// ID Filter game states by ids
	ID *[]int64 `form:"id,omitempty" json:"id,omitempty"`

	// Name Filter game states by name (partial match)
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Sort Sort by fields (add prefix `-` for descending e.g. -created_at)
	Sort *[]string `form:"sort,omitempty" json:"sort,omitempty"`
}

type ListGameTimeCategoriesFilters struct {
	api.ListGameTimeCategoriesParams
}

type ListGameTimeKindsFilters struct {
	api.ListGameTimeKindsParams
}

type ListGameVariantsFilters struct {
	api.ListGameVariantsParams
}

type GameStats struct {
	Win   int32 `json:"win"`
	Loss  int32 `json:"loss"`
	Draw  int32 `json:"draw"`
	Total int32 `json:"total"`
}

type ListGamesFilters struct {
	api.ListGamesParams
	WithGameHashes bool
}

type GetGameByIDFilters struct {
	api.GetGameParams
	WithGameHashes bool
}

type GameStatsForUserFilters struct {
	From *time.Time
	To   *time.Time
}

type GameDetails struct {
	models.Game
	GameMoves         types.JSON[*[]models.GameMove]
	GameHistoryHashes types.JSON[*[]models.GameHistoryHash]
}

type GetActiveGameFilters struct {
	WithGameMoves  bool
	WithGameHashes bool
}

type ListActiveGameFilters struct {
	WithGameMoves  bool
	WithGameHashes bool
}

type ActiveGame struct {
	GameID  int64  `json:"game_id"`
	WhiteID string `json:"white_id"`
	BlackID string `json:"black_id"`
	// WhiteID                *uuid.UUID          `json:"white_id"`
	// BlackID                *uuid.UUID          `json:"black_id"`
	// GuestWhiteID           *uuid.UUID          `json:"guest_white_id"`
	// GuestBlackID           *uuid.UUID          `json:"guest_black_id"`
	GameVariant            pb.GameVariant                    `json:"game_variant"`
	GameTimeKind           pb.GameTimeKind                   `json:"game_time_kind"`
	GameTimeCategory       pb.GameTimeCategory               `json:"game_time_category"`
	TimeControlClockMs     int32                             `json:"time_control_clock_ms"`
	TimeControlIncrementMs int32                             `json:"time_control_increment_ms"`
	GameResult             pb.GameResult                     `json:"game_result"`
	GameResultStatus       pb.GameResultStatus               `json:"game_result_status"`
	GameState              pb.GameState                      `json:"game_state"`
	ReconnectTimeoutMs     int32                             `json:"reconnect_timeout_ms"`
	FirstMoveTimeoutMs     int32                             `json:"first_move_timeout_ms"`
	LastMove               *time.Time                        `json:"last_move"`
	StartTime              *time.Time                        `json:"start_time"`
	EndTime                *time.Time                        `json:"end_time"`
	Rated                  bool                              `json:"rated"`
	GameMoves              []ActiveGameMove                  `json:"game_moves"`
	GameHistoryHashes      []ActiveGameHistoryHash           `json:"game_history_hashes"`
	Version                int32                             `json:"version"`
	PendingDrawOffers      map[uuid.UUID]*gameplay.DrawOffer `json:"pending_draw_offers" `
	WhiteGameRemainingSecs int32                             `json:"white_game_remaining_secs" `
	WhiteGameRemainingNS   int64                             `json:"white_game_remaining_ns" `
	BlackGameRemainingSecs int32                             `json:"black_game_remaining_secs" `
	BlackGameRemainingNS   int64                             `json:"black_game_remaining_ns" `
}

type ActiveGameMove struct {
	ID       int64      `json:"id"`
	Fen      string     `json:"fen"`
	Uci      string     `json:"uci"`
	San      string     `json:"san"`
	PlayedAt *time.Time `json:"played_at"`
}

type ActiveGameHistoryHash struct {
	ID   int64 `json:"id"`
	Hash int64 `json:"hash"`
}

type ListRatingsFilters struct {
	api.ListRatingsParams
}
