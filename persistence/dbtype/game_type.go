package dbtype

import (
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/stephenafamo/bob/types"
)

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

type GameWithJoinData struct {
	models.Game
	GameMoves         types.JSON[*[]models.GameMove]
	GameHistoryHashes types.JSON[*[]models.GameHistoryHash]
}
