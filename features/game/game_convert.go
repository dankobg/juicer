package game

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
		ID:           g.ID,
		WhiteID:      g.WhiteID.Ptr(),
		BlackID:      g.BlackID.Ptr(),
		WhiteIsGuest: whiteIsGuest,
		BlackIsGuest: blackIsGuest,
		Rated:        g.Rated,
		WhiteGuestID: g.GuestWhiteID.Ptr(),
		BlackGuestID: g.GuestBlackID.Ptr(),
		// WhiteGameClock:         int64(g.WhiteGameClock), // @TODO: check lejtaaaaaaaaa
		// BlackGameClock:         int64(g.BlackGameClock), // @TODO: check lejtaaaaaaaaa
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

func GameDetailsToResponse(g GameDetails) api.Game {
	whiteIsGuest, blackIsGuest := true, true
	if g.WhiteID.IsNull() {
		whiteIsGuest = false
	}

	if g.BlackID.IsNull() {
		blackIsGuest = false
	}

	game := api.Game{
		ID:           g.ID,
		WhiteID:      g.WhiteID.Ptr(),
		BlackID:      g.BlackID.Ptr(),
		WhiteIsGuest: whiteIsGuest,
		BlackIsGuest: blackIsGuest,
		Rated:        g.Rated,
		WhiteGuestID: g.GuestWhiteID.Ptr(),
		BlackGuestID: g.GuestBlackID.Ptr(),
		// WhiteGameClock:         int64(g.WhiteGameClock), // @TODO: check lejtaaaaaaaaa
		// BlackGameClock:         int64(g.BlackGameClock), // @TODO: check lejtaaaaaaaaa
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

	if g.GameMoves.Val != nil && len(*g.GameMoves.Val) > 0 {
		moves := make([]api.GameMove, len(*g.GameMoves.Val))
		for i, m := range *g.GameMoves.Val {
			moves[i] = api.GameMove{
				ID:       m.ID,
				GameID:   m.GameID,
				Fen:      m.Fen,
				San:      m.San,
				Uci:      m.Uci,
				PlayedAt: m.PlayedAt.Ptr(),
			}
		}

		game.Moves = &moves
	}

	return game
}

func GameTimeCategoryToResponse(tc models.GameTimeCategory) api.GameTimeCategory {
	return api.GameTimeCategory{
		ID:                 tc.ID,
		Name:               tc.Name,
		UpperTimeLimitSecs: tc.UpperTimeLimitSecs.Ptr(),
		CreatedAt:          tc.CreatedAt,
		UpdatedAt:          tc.UpdatedAt,
	}
}

func GameTimeKindToResponse(tk models.GameTimeKind) api.GameTimeKind {
	return api.GameTimeKind{
		ID:        tk.ID,
		Name:      tk.Name,
		Enabled:   tk.Enabled,
		CreatedAt: tk.CreatedAt,
		UpdatedAt: tk.UpdatedAt,
	}
}

func GameVariantToResponse(gv models.GameVariant) api.GameVariant {
	return api.GameVariant{
		ID:        gv.ID,
		Enabled:   gv.Enabled,
		Name:      gv.Name,
		CreatedAt: gv.CreatedAt,
		UpdatedAt: gv.UpdatedAt,
	}
}

func RatingToResponse(r models.Rating) api.Rating {
	return api.Rating{
		ID:                 r.ID,
		UserID:             r.UserID,
		GameTimeCategoryID: r.GameTimeCategoryID,
		Glicko:             int64(r.Glicko),
		Glicko2:            int64(r.Glicko2),
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

func GameStatsToResponse(gs GameStats) api.GameStats {
	return api.GameStats{
		All: api.GameStat{
			Draw:        gs.All.Draw,
			Interrupted: gs.All.Interrupted,
			Loss:        gs.All.Loss,
			Total:       new(gs.All.Total),
			Win:         gs.All.Win,
		},
		Blitz: api.GameStat{
			Draw:        gs.Blitz.Draw,
			Interrupted: gs.Blitz.Interrupted,
			Loss:        gs.Blitz.Loss,
			Total:       new(gs.Blitz.Total),
			Win:         gs.Blitz.Win,
		},
		Bullet: api.GameStat{
			Draw:        gs.Bullet.Draw,
			Interrupted: gs.Bullet.Interrupted,
			Loss:        gs.Bullet.Loss,
			Total:       new(gs.Bullet.Total),
			Win:         gs.Bullet.Win,
		},
		Classical: api.GameStat{
			Draw:        gs.Classical.Draw,
			Interrupted: gs.Classical.Interrupted,
			Loss:        gs.Classical.Loss,
			Total:       new(gs.Classical.Total),
			Win:         gs.Classical.Win,
		},
		Hyperbullet: api.GameStat{
			Draw:        gs.Hyperbullet.Draw,
			Interrupted: gs.Hyperbullet.Interrupted,
			Loss:        gs.Hyperbullet.Loss,
			Total:       new(gs.Hyperbullet.Total),
			Win:         gs.Hyperbullet.Win,
		},
		Rapid: api.GameStat{
			Draw:        gs.Rapid.Draw,
			Interrupted: gs.Rapid.Interrupted,
			Loss:        gs.Rapid.Loss,
			Total:       new(gs.Rapid.Total),
			Win:         gs.Rapid.Win,
		},
	}
}
