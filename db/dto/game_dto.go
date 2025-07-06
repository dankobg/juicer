package dto

import (
	"time"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/opt"
	"github.com/dankobg/juicer/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type GameChangeset struct {
	ID                   opt.Val[uuid.UUID] `json:"id,omitempty"`
	WhiteID              opt.Val[uuid.UUID] `json:"white_id,omitempty"`
	BlackID              opt.Val[uuid.UUID] `json:"black_id,omitempty"`
	GuestWhiteID         opt.Val[uuid.UUID] `json:"guest_white_id,omitempty"`
	GuestBlackID         opt.Val[uuid.UUID] `json:"guest_black_id,omitempty"`
	VariantID            opt.Val[uuid.UUID] `json:"variant_id,omitempty"`
	TimeKindID           opt.Val[uuid.UUID] `json:"time_kind_id,omitempty"`
	TimeCategoryID       opt.Val[uuid.UUID] `json:"time_category_id,omitempty"`
	IsGuest              opt.Val[bool]      `json:"is_guest,omitempty"`
	TimeControlClock     opt.Val[int32]     `json:"time_control_clock,omitempty"`
	TimeControlIncrement opt.Val[int32]     `json:"time_control_increment,omitempty"`
	ReconnectTimeout     opt.Val[int32]     `json:"reconnect_timer,omitempty"`
	FirstMoveTimeout     opt.Val[int32]     `json:"first_move_timer,omitempty"`
	WhiteGameClock       opt.Val[int32]     `json:"white_game_clock,omitempty"`
	BlackGameClock       opt.Val[int32]     `json:"black_game_clock,omitempty"`
	ResultID             opt.Val[uuid.UUID] `json:"result_id,omitempty"`
	ResultStatusID       opt.Val[uuid.UUID] `json:"result_status_id,omitempty"`
	StateID              opt.Val[uuid.UUID] `json:"state_id,omitempty"`
	StartTime            opt.Val[time.Time] `json:"start_time,omitempty"`
	EndTime              opt.Val[time.Time] `json:"end_time,omitempty"`
	LastMove             opt.Val[time.Time] `json:"last_move,omitempty"`
	Fen                  opt.Val[string]    `json:"fen,omitempty"`
	Pgn                  opt.Val[string]    `json:"pgn,omitempty"`
	CreatedAt            opt.Val[time.Time] `json:"created_at,omitempty"`
	UpdatedAt            opt.Val[time.Time] `json:"updated_at,omitempty"`
}

func (gs GameChangeset) ToModel() (p.ColumnList, model.Game) {
	var cols p.ColumnList
	var m model.Game

	cols = append(cols, t.Game.UpdatedAt)
	m.UpdatedAt = time.Now()

	if gs.ID.IsSpecified() {
		cols = append(cols, t.Game.ID)
		if !gs.ID.IsNull() {
			m.ID = gs.ID.MustGet()
		}
	}
	if gs.WhiteID.IsSpecified() {
		cols = append(cols, t.Game.WhiteID)
		if !gs.WhiteID.IsNull() {
			m.WhiteID = ptr.Of(gs.WhiteID.MustGet())
		}
	}
	if gs.BlackID.IsSpecified() {
		cols = append(cols, t.Game.BlackID)
		if !gs.BlackID.IsNull() {
			m.BlackID = ptr.Of(gs.BlackID.MustGet())
		}
	}
	if gs.IsGuest.IsSpecified() {
		cols = append(cols, t.Game.IsGuest)
		if !gs.IsGuest.IsNull() {
			m.IsGuest = gs.IsGuest.MustGet()
		}
	}
	if gs.GuestWhiteID.IsSpecified() {
		cols = append(cols, t.Game.GuestWhiteID)
		if !gs.GuestWhiteID.IsNull() {
			m.GuestWhiteID = ptr.Of(gs.GuestWhiteID.MustGet())
		}
	}
	if gs.GuestBlackID.IsSpecified() {
		cols = append(cols, t.Game.GuestBlackID)
		if !gs.GuestBlackID.IsNull() {
			m.GuestBlackID = ptr.Of(gs.GuestBlackID.MustGet())
		}
	}
	if gs.VariantID.IsSpecified() {
		cols = append(cols, t.Game.VariantID)
		if !gs.VariantID.IsNull() {
			m.VariantID = gs.VariantID.MustGet()
		}
	}
	if gs.TimeKindID.IsSpecified() {
		cols = append(cols, t.Game.TimeKindID)
		if !gs.TimeKindID.IsNull() {
			m.TimeKindID = gs.TimeKindID.MustGet()
		}
	}
	if gs.TimeCategoryID.IsSpecified() {
		cols = append(cols, t.Game.TimeCategoryID)
		if !gs.TimeCategoryID.IsNull() {
			m.TimeCategoryID = gs.TimeCategoryID.MustGet()
		}
	}
	if gs.ResultID.IsSpecified() {
		cols = append(cols, t.Game.ResultID)
		if !gs.ResultID.IsNull() {
			m.ResultID = ptr.Of(gs.ResultID.MustGet())
		}
	}
	if gs.ResultStatusID.IsSpecified() {
		cols = append(cols, t.Game.ResultStatusID)
		if !gs.ResultStatusID.IsNull() {
			m.ResultStatusID = ptr.Of(gs.ResultStatusID.MustGet())
		}
	}
	if gs.StateID.IsSpecified() {
		cols = append(cols, t.Game.StateID)
		if !gs.StateID.IsNull() {
			m.StateID = gs.StateID.MustGet()
		}
	}
	if gs.TimeControlClock.IsSpecified() {
		cols = append(cols, t.Game.TimeControlClock)
		if !gs.TimeControlClock.IsNull() {
			m.TimeControlClock = gs.TimeControlClock.MustGet()
		}
	}
	if gs.TimeControlIncrement.IsSpecified() {
		cols = append(cols, t.Game.TimeControlIncrement)
		if !gs.TimeControlIncrement.IsNull() {
			m.TimeControlIncrement = gs.TimeControlIncrement.MustGet()
		}
	}
	if gs.ReconnectTimeout.IsSpecified() {
		cols = append(cols, t.Game.ReconnectTimeout)
		if !gs.ReconnectTimeout.IsNull() {
			m.ReconnectTimeout = gs.ReconnectTimeout.MustGet()
		}
	}
	if gs.FirstMoveTimeout.IsSpecified() {
		cols = append(cols, t.Game.FirstMoveTimeout)
		if !gs.FirstMoveTimeout.IsNull() {
			m.FirstMoveTimeout = gs.FirstMoveTimeout.MustGet()
		}
	}
	if gs.WhiteGameClock.IsSpecified() {
		cols = append(cols, t.Game.WhiteGameClock)
		if !gs.WhiteGameClock.IsNull() {
			m.WhiteGameClock = gs.WhiteGameClock.MustGet()
		}
	}
	if gs.BlackGameClock.IsSpecified() {
		cols = append(cols, t.Game.BlackGameClock)
		if !gs.BlackGameClock.IsNull() {
			m.BlackGameClock = gs.BlackGameClock.MustGet()
		}
	}
	if gs.StartTime.IsSpecified() {
		cols = append(cols, t.Game.StartTime)
		if !gs.StartTime.IsNull() {
			m.StartTime = gs.StartTime.MustGet()
		}
	}
	if gs.EndTime.IsSpecified() {
		cols = append(cols, t.Game.EndTime)
		if !gs.EndTime.IsNull() {
			m.EndTime = ptr.Of(gs.EndTime.MustGet())
		}
	}
	if gs.LastMove.IsSpecified() {
		cols = append(cols, t.Game.LastMove)
		if !gs.LastMove.IsNull() {
			m.LastMove = ptr.Of(gs.LastMove.MustGet())
		}
	}
	if gs.Fen.IsSpecified() {
		cols = append(cols, t.Game.Fen)
		if !gs.Fen.IsNull() {
			m.Fen = gs.Fen.MustGet()
		}
	}
	if gs.Pgn.IsSpecified() {
		cols = append(cols, t.Game.Pgn)
		if !gs.Pgn.IsNull() {
			m.Pgn = ptr.Of(gs.Pgn.MustGet())
		}
	}

	return cols, m
}
