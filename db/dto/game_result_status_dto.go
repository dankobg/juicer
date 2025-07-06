package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
)

type GameResultStatusChangeset struct {
	Name      opt.Val[string]    `json:"name,omitempty"`
	CreatedAt opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (grs GameResultStatusChangeset) ToModel() (p.ColumnList, model.GameResultStatus) {
	var cols p.ColumnList
	var m model.GameResultStatus

	cols = append(cols, t.GameResultStatus.UpdatedAt)
	m.UpdatedAt = time.Now()

	if grs.Name.IsSpecified() {
		cols = append(cols, t.GameResultStatus.Name)
		if !grs.Name.IsNull() {
			m.Name = grs.Name.MustGet()
		}
	}

	return cols, m
}
