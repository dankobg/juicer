package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type RatingChangeset struct {
	UserID             opt.Val[uuid.UUID] `json:"user_id,omitempty"`
	GameTimeCategoryID opt.Val[uuid.UUID] `json:"game_time_category_id,omitempty"`
	Glicko             opt.Val[int32]     `json:"glicko,omitempty"`
	Glicko2            opt.Val[int32]     `json:"glicko2,omitempty"`
	CreatedAt          opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt          opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (rt RatingChangeset) ToModel() (p.ColumnList, model.Rating) {
	var cols p.ColumnList
	var m model.Rating

	cols = append(cols, t.Rating.UpdatedAt)
	m.UpdatedAt = time.Now()

	if rt.UserID.IsSpecified() {
		cols = append(cols, t.Rating.UserID)
		if !rt.UserID.IsNull() {
			m.UserID = rt.UserID.MustGet()
		}
	}
	if rt.GameTimeCategoryID.IsSpecified() {
		cols = append(cols, t.Rating.GameTimeCategoryID)
		if !rt.GameTimeCategoryID.IsNull() {
			m.GameTimeCategoryID = rt.GameTimeCategoryID.MustGet()
		}
	}
	if rt.Glicko.IsSpecified() {
		cols = append(cols, t.Rating.Glicko)
		if !rt.Glicko.IsNull() {
			m.Glicko = rt.Glicko.MustGet()
		}
	}
	if rt.Glicko2.IsSpecified() {
		cols = append(cols, t.Rating.Glicko2)
		if !rt.Glicko2.IsNull() {
			m.Glicko2 = rt.Glicko2.MustGet()
		}
	}

	return cols, m
}
