package store

import "time"

type GameStatsForUserFilters struct {
	From *time.Time
	To   *time.Time
}
