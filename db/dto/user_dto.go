package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type UserChangeset struct {
	ID        opt.Val[uuid.UUID] `json:"id,omitempty"`
	CreatedAt opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (gs UserChangeset) ToModel() (p.ColumnList, model.User) {
	var cols p.ColumnList
	var m model.User

	cols = append(cols, t.User.UpdatedAt)
	m.UpdatedAt = time.Now()

	if gs.ID.IsSpecified() {
		cols = append(cols, t.User.ID)
		if !gs.ID.IsNull() {
			m.ID = gs.ID.MustGet()
		}
	}

	return cols, m
}
