package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	"github.com/dankobg/juicer/ptr"
	p "github.com/go-jet/jet/v2/postgres"
)

type GameTimeCategoryChangeset struct {
	Name               opt.Val[string]    `json:"name,omitempty"`
	UpperTimeLimitSecs opt.Val[int32]     `json:"upper_time_limit_secs,omitempty"`
	CreatedAt          opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt          opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (grs GameTimeCategoryChangeset) ToModel() (p.ColumnList, model.GameTimeCategory) {
	var cols p.ColumnList
	var m model.GameTimeCategory

	cols = append(cols, t.GameTimeCategory.UpdatedAt)
	m.UpdatedAt = time.Now()

	if grs.Name.IsSpecified() {
		cols = append(cols, t.GameTimeCategory.Name)
		if !grs.Name.IsNull() {
			m.Name = grs.Name.MustGet()
		}
	}
	if grs.UpperTimeLimitSecs.IsSpecified() {
		cols = append(cols, t.GameTimeCategory.UpperTimeLimitSecs)
		if !grs.UpperTimeLimitSecs.IsNull() {
			m.UpperTimeLimitSecs = ptr.Of(grs.UpperTimeLimitSecs.MustGet())
		}
	}

	return cols, m
}
