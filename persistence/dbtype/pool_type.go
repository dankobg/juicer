package dbtype

import "fmt"

type Pool struct {
	ClockMS     int32
	IncrementMS int32
	Rated       bool
}

func (p Pool) Name() string {
	rated := "unrated"
	if p.Rated {
		rated = "rated"
	}

	return fmt.Sprintf("%d_%d_%s", p.ClockMS, p.IncrementMS, rated)
}
