package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func GameToResponse(g models.Game) api.Game {
	whiteIsGuest, blackIsGuest := true, true
	if g.WhiteID.IsNull() {
		whiteIsGuest = false
	}
	if g.BlackID.IsNull() {
		blackIsGuest = false
	}

	return api.Game{
		ID:                     g.ID,
		WhiteID:                g.WhiteID.Ptr(),
		BlackID:                g.BlackID.Ptr(),
		WhiteIsGuest:           whiteIsGuest,
		BlackIsGuest:           blackIsGuest,
		Rated:                  g.Rated,
		WhiteGuestID:           g.GuestWhiteID.Ptr(),
		BlackGuestID:           g.GuestBlackID.Ptr(),
		WhiteGameClock:         int64(g.WhiteGameClock),
		BlackGameClock:         int64(g.BlackGameClock),
		GameVariantID:          g.GameVariantID,
		GameStateID:            g.GameStateID,
		GameTimeKindID:         g.GameTimeKindID,
		GameTimeCategoryID:     g.GameTimeCategoryID,
		GameResultID:           g.GameResultID.Ptr(),
		GameResultStatusID:     g.GameResultStatusID.Ptr(),
		FirstMoveTimeoutMs:     int64(g.FirstMoveTimeoutMS),
		ReconnectTimeoutMs:     int64(g.ReconnectTimeoutMS),
		Fen:                    g.Fen,
		Pgn:                    g.PGN.Ptr(),
		StartTime:              g.StartTime.Ptr(),
		EndTime:                g.EndTime.Ptr(),
		LastMove:               g.LastMove.Ptr(),
		TimeControlClockMs:     int64(g.TimeControlClockMS),
		TimeControlIncrementMs: int64(g.TimeControlIncrementMS),
		CreatedAt:              g.CreatedAt,
		UpdatedAt:              g.UpdatedAt,
	}
}
