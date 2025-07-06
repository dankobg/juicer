package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
)

type GameTimeKindChangeset struct {
	Name      opt.Val[string]    `json:"name,omitempty"`
	Enabled   opt.Val[bool]      `json:"enabled,omitempty"`
	CreatedAt opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (grs GameTimeKindChangeset) ToModel() (p.ColumnList, model.GameTimeKind) {
	var cols p.ColumnList
	var m model.GameTimeKind

	cols = append(cols, t.GameTimeKind.UpdatedAt)
	m.UpdatedAt = time.Now()

	if grs.Name.IsSpecified() {
		cols = append(cols, t.GameTimeKind.Name)
		if !grs.Name.IsNull() {
			m.Name = grs.Name.MustGet()
		}
	}
	if grs.Enabled.IsSpecified() {
		cols = append(cols, t.GameTimeKind.Enabled)
		if !grs.Enabled.IsNull() {
			m.Enabled = grs.Enabled.MustGet()
		}
	}

	return cols, m
}
