package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type GameMoveChangeset struct {
	GameID   opt.Val[uuid.UUID] `json:"game_id,omitempty"`
	Fen      opt.Val[string]    `json:"fen,omitempty"`
	Uci      opt.Val[string]    `json:"uci,omitempty"`
	San      opt.Val[string]    `json:"san,omitempty"`
	PlayedAt opt.Val[time.Time] `json:"played_at,omitempty"`
}

func (gmc GameMoveChangeset) ToModel() (p.ColumnList, model.GameMove) {
	var cols p.ColumnList
	var m model.GameMove

	if gmc.GameID.IsSpecified() {
		cols = append(cols, t.GameMove.GameID)
		if !gmc.GameID.IsNull() {
			m.GameID = gmc.GameID.MustGet()
		}
	}
	if gmc.Fen.IsSpecified() {
		cols = append(cols, t.GameMove.Fen)
		if !gmc.Fen.IsNull() {
			m.Fen = gmc.Fen.MustGet()
		}
	}
	if gmc.Uci.IsSpecified() {
		cols = append(cols, t.GameMove.Uci)
		if !gmc.Uci.IsNull() {
			m.Uci = gmc.Uci.MustGet()
		}
	}
	if gmc.San.IsSpecified() {
		cols = append(cols, t.GameMove.San)
		if !gmc.San.IsNull() {
			m.San = gmc.San.MustGet()
		}
	}
	if gmc.PlayedAt.IsSpecified() {
		cols = append(cols, t.GameMove.PlayedAt)
		if !gmc.PlayedAt.IsNull() {
			m.PlayedAt = gmc.PlayedAt.MustGet()
		}
	}

	return cols, m
}
