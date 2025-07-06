package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
)

type GameVariantChangeset struct {
	Name      opt.Val[string]    `json:"name,omitempty"`
	Enabled   opt.Val[bool]      `json:"enabled,omitempty"`
	CreatedAt opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (gs GameVariantChangeset) ToModel() (p.ColumnList, model.GameVariant) {
	var cols p.ColumnList
	var m model.GameVariant

	cols = append(cols, t.GameVariant.UpdatedAt)
	m.UpdatedAt = time.Now()

	if gs.Name.IsSpecified() {
		cols = append(cols, t.GameVariant.Name)
		if !gs.Name.IsNull() {
			m.Name = gs.Name.MustGet()
		}
	}
	if gs.Enabled.IsSpecified() {
		cols = append(cols, t.GameVariant.Enabled)
		if !gs.Enabled.IsNull() {
			m.Enabled = gs.Enabled.MustGet()
		}
	}

	return cols, m
}
