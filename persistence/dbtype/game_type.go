package dbtype

import (
	"time"

	api "github.com/dankobg/juicer/api/gen"
)

type ListGamesFilters struct {
	api.ListGamesParams
}

type GameStatsForUserFilters struct {
	From *time.Time
	To   *time.Time
}
