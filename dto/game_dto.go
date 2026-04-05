package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func GameToResponse(g models.Game) api.Game {
	return api.Game{
		ID:                   g.ID,
		WhiteID:              g.WhiteID.Ptr(),
		BlackID:              g.BlackID.Ptr(),
		WhiteGuestID:         g.GuestWhiteID.Ptr(),
		BlackGuestID:         g.GuestBlackID.Ptr(),
		IsGuest:              g.IsGuest,
		WhiteGameClock:       int64(g.WhiteGameClock),
		BlackGameClock:       int64(g.BlackGameClock),
		VariantID:            g.VariantID,
		FirstMoveTimeout:     int64(g.FirstMoveTimeout),
		ReconnectTimeout:     int64(g.ReconnectTimeout),
		Fen:                  g.Fen,
		Pgn:                  g.PGN.Ptr(),
		ResultID:             g.ResultID.Ptr(),
		ResultStatusID:       g.ResultStatusID.Ptr(),
		StartTime:            g.StartTime,
		EndTime:              g.EndTime.Ptr(),
		LastMove:             g.LastMove.Ptr(),
		StateID:              g.StateID,
		TimeCategoryID:       g.TimeCategoryID,
		TimeControlClock:     int64(g.TimeControlClock),
		TimeControlIncrement: int64(g.TimeControlIncrement),
		TimeKindID:           g.TimeKindID,
		CreatedAt:            g.CreatedAt,
		UpdatedAt:            g.UpdatedAt,
	}
}
