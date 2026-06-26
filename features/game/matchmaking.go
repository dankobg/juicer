package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/gameplay"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type Pool struct {
	ClockMS     int32
	IncrementMS int32
	Rated       bool
}

func (p Pool) Name() string {
	rated := "unrated"
	if p.Rated {
		rated = "rated"
	}

	return fmt.Sprintf("%d_%d_%s", p.ClockMS, p.IncrementMS, rated)
}

func (g *GameService) StartMatchmaking(ctx context.Context) {
	g.log.Info("game matchmaking started")

	matchmakingInterval := time.Second * 15
	ticker := time.NewTicker(matchmakingInterval)

loop:
	for {
		select {
		case <-ticker.C:
			g.tryMatchPoolPlayers(ctx)
		case <-ctx.Done():
			break loop
		}
	}
}

func (g *GameService) tryMatchPoolPlayers(ctx context.Context) {
	g.log.Debug("matchmaking trying to match pool players")

	for _, quickGame := range quickGames {
		g.tryMatchPoolPlayersForPool(ctx, quickGame.ClockSecs*1000, quickGame.IncrementSecs*1000, true)
		g.tryMatchPoolPlayersForPool(ctx, quickGame.ClockSecs*1000, quickGame.IncrementSecs*1000, false)
	}
}

func (g *GameService) tryMatchPoolPlayersForPool(ctx context.Context, clockMS, incrementMS int32, rated bool) {
	pool := Pool{ClockMS: clockMS, IncrementMS: incrementMS, Rated: rated}

	// mm.log.Debug("try match pool", slog.String("pool", pool.Name()))

	matchedPairs := make([][2]string, 0)

	for {
		res, err := g.pst.Pool.MatchPair(ctx, pool)
		if err != nil {
			break
		}

		if len(res) < 2 {
			break
		}

		matchedPairs = append(matchedPairs, [2]string{res[0], res[1]})
	}

	if len(matchedPairs) == 0 {
		return
	}

	var wg sync.WaitGroup

	for _, pair := range matchedPairs {
		wg.Go(func() {
			g.processPoolMatchFound(ctx, pair, pool)
		})
	}

	wg.Wait()
}

func (g *GameService) processPoolMatchFound(ctx context.Context, pair [2]string, pool Pool) {
	g.log.Debug("processing matched pool pair", slog.String("pool", pool.Name()), slog.Any("pair", pair))

	userID1, err1 := uuid.Parse(pair[0])

	userID2, err2 := uuid.Parse(pair[1])
	if err1 != nil || err2 != nil {
		g.log.Error("invalid matched pair user id", slog.Any("error", errors.Join(err1, err2)))
	}

	username1, username2 := "guest", "guest"

	if pool.Rated {
		uname1, err5 := g.usrRdr.GetUsername(ctx, userID1)

		uname2, err6 := g.usrRdr.GetUsername(ctx, userID2)
		if err5 != nil || err6 != nil {
			g.log.Error("failed to get usernames", slog.Any("error", errors.Join(err5, err6)))
		}

		username1, username2 = uname1, uname2
	}

	color1, color2 := pb.Color_COLOR_WHITE, pb.Color_COLOR_BLACK
	if rand.IntN(2) == 1 {
		color1, color2 = pb.Color_COLOR_BLACK, pb.Color_COLOR_WHITE
	}

	players := [2]gameplay.Player{
		{ID: userID1, Username: username1, Color: color1, Guest: !pool.Rated},
		{ID: userID2, Username: username2, Color: color2, Guest: !pool.Rated},
	}

	gtc := &pb.GameTimeControl{ClockMs: pool.ClockMS, IncrementMs: pool.IncrementMS}

	thresholds := []gameplay.CategoryThreshold{}
	for _, x := range g.categoryThresholds {
		thresholds = append(thresholds, gameplay.CategoryThreshold{
			UpperLimit:   x.upperLimit,
			TimeCategory: x.timeCategory,
		})
	}

	gs, err := gameplay.NewGameState(-1, players, gtc, thresholds, g.gameEvent, gameplay.WithRated(pool.Rated))
	if err != nil {
		g.log.Error("gameplay.NewGameState", slog.Any("error", err))
		return
	}

	gameRemainingSecs := gs.GameTimeControl.GetClockMs() / 1000
	gameRemainingNs := (int64(gs.GameTimeControl.GetClockMs()) % 1000) * int64(time.Millisecond)

	gameSetter := models.GameSetter{
		GameVariantID:          omit.From(g.gameVariantProtoToID(gs.GameVariant)),
		GameTimeKindID:         omit.From(g.gameTimeKindProtoToID(gs.GameTimeKind)),
		GameTimeCategoryID:     omit.From(g.gameTimeCategoryProtoToID(gs.GameTimeCategory)),
		GameStateID:            omit.From(g.gameStateProtoToID(gs.GameState)),
		TimeControlClockMS:     omit.From(gs.GameTimeControl.ClockMs),
		TimeControlIncrementMS: omit.From(gs.GameTimeControl.IncrementMs),
		FirstMoveTimeoutMS:     omit.From(int32(gs.FirstMoveTimeout.Milliseconds())),
		ReconnectTimeoutMS:     omit.From(int32(gs.ReconnectTimeout.Milliseconds())),
		WhiteGameRemainingSecs: omit.From(gameRemainingSecs),
		WhiteGameRemainingNS:   omit.From(gameRemainingNs),
		BlackGameRemainingSecs: omit.From(gameRemainingSecs),
		BlackGameRemainingNS:   omit.From(gameRemainingNs),
		Rated:                  omit.From(gs.Rated),
		StartTime:              omitnull.FromPtr(gs.StartTime),
		EndTime:                omitnull.FromPtr(gs.EndTime),
		LastMove:               omitnull.FromPtr(gs.LastMove),
		Fen:                    omit.From(gs.Chess.Position.Fen()),
		Repetitions:            omit.From(int32(gs.Chess.Repetitions)),
	}
	if pool.Rated {
		gameSetter.WhiteID = omitnull.From(gs.White.ID)
		gameSetter.BlackID = omitnull.From(gs.Black.ID)
	} else {
		gameSetter.GuestWhiteID = omitnull.From(gs.White.ID)
		gameSetter.GuestBlackID = omitnull.From(gs.Black.ID)
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		gameSetter.GameResultID = omitnull.From(g.gameResultProtoToID(gs.GameResult))
	}

	if gs.GameResultStatus != pb.GameResultStatus_GAME_RESULT_STATUS_UNSPECIFIED {
		gameSetter.GameResultStatusID = omitnull.From(g.gameResultStatusProtoToID(gs.GameResultStatus))
	}

	moveSetters := make([]models.GameMoveSetter, len(gs.GameMoves))
	for i, move := range gs.GameMoves {
		moveSetter := models.GameMoveSetter{
			Fen: omit.From(move.GetFen()),
			Uci: omit.From(move.GetUci()),
			San: omit.From(move.GetSan()),
			Lan: omit.From(move.GetLan()),
		}
		if move.GetPlayedAt() != nil {
			moveSetter.PlayedAt = omitnull.From(move.GetPlayedAt().AsTime())
		}

		moveSetters[i] = moveSetter
	}

	hashSetters := make([]models.GameHistoryHashSetter, len(gs.Chess.HistoryHashes))
	for i, hash := range gs.Chess.HistoryHashes {
		hashSetters[i] = models.GameHistoryHashSetter{
			Hash: omit.From(int64(hash)),
		}
	}

	gs.Start(ctx)

	game, err := g.pst.Game.CreateGame(ctx, gameSetter, moveSetters, hashSetters)
	if err != nil {
		g.log.Error("CreateGame", slog.Any("error", err))
		return
	}

	gs.GameID = game.ID

	g.gamestates[game.ID] = gs

	if err := g.pst.ActiveGame.CreateActiveGame(ctx, gs); err != nil {
		g.log.Error("CreateActiveGame", slog.Any("error", err))
	}

	gameFoundMsg := &pb.Message{Event: &pb.Message_GameFound{GameFound: &pb.GameFound{GameId: int32(gs.GameID)}}}

	gameFoundMsgBytes, err := protojson.Marshal(gameFoundMsg)
	if err != nil {
		g.log.Error("protojson marshal Message_GameFound", slog.Any("error", err))
	} else {
		if err := g.bus.Publish(ctx, "user."+userID1.String(), gameFoundMsgBytes); err != nil {
			g.log.Error("publish Message_GameFound", slog.Any("error", err))
		}

		if err := g.bus.Publish(ctx, "user."+userID2.String(), gameFoundMsgBytes); err != nil {
			g.log.Error("publish Message_GameFound", slog.Any("error", err))
		}
	}
}
