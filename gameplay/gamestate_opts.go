package gameplay

import (
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
)

type gameOpts struct {
	gameID           int64
	fen              string
	white            *Player
	black            *Player
	players          map[uuid.UUID]*Player
	gameVariant      pb.GameVariant
	gameTimeKind     pb.GameTimeKind
	gameTimeCategory pb.GameTimeCategory
	gameTimeControl  *pb.GameTimeControl
	gameResult       pb.GameResult
	gameResultStatus pb.GameResultStatus
	gameState        pb.GameState
	reconnectTimeout time.Duration
	firstMoveTimeout time.Duration
	rated            bool
	lastMove         *time.Time
	startTime        *time.Time
	endTime          *time.Time
	version          int
	gameMoves        []*pb.GameMove
}

type GameOption interface {
	apply(*gameOpts)
}

type gameOptionFunc func(*gameOpts)

func (f gameOptionFunc) apply(g *gameOpts) {
	f(g)
}

type GameOptions []GameOption

func (o GameOptions) apply(g *gameOpts) {
	for _, opt := range o {
		opt.apply(g)
	}
}

func WithGameID(gameID int64) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameID = gameID
	})
}

func WithFEN(fen string) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.fen = fen
	})
}

func WithGameVariant(gameVariant pb.GameVariant) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameVariant = gameVariant
	})
}

func WithGameTimeKind(gameTimeKind pb.GameTimeKind) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameTimeKind = gameTimeKind
	})
}

func WithGameTimeCategory(gameTimeCategory pb.GameTimeCategory) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameTimeCategory = gameTimeCategory
	})
}

func WithGameTimeControl(gtc *pb.GameTimeControl) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameTimeControl = gtc
	})
}

func WithGameResult(gameResult pb.GameResult) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameResult = gameResult
	})
}

func WithGameResultStatus(gameResultStatus pb.GameResultStatus) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameResultStatus = gameResultStatus
	})
}

func WithGameState(gameState pb.GameState) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameState = gameState
	})
}

func WithReconnectTimeout(dur time.Duration) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.reconnectTimeout = dur
	})
}

func WithFirstMoveTimeout(dur time.Duration) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.firstMoveTimeout = dur
	})
}

func WithRated(rated bool) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.rated = rated
	})
}

func WithLastMove(lastMove *time.Time) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.lastMove = lastMove
	})
}

func WithStartTime(startTime *time.Time) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.startTime = startTime
	})
}

func WithEndTime(endTime *time.Time) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.endTime = endTime
	})
}

func WithGameMoves(gameMoves []*pb.GameMove) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.gameMoves = gameMoves
	})
}

func WithWhitePlayer(p *Player) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.white = p
	})
}

func WithBlackPlayer(p *Player) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.black = p
	})
}

func WithPlayers(players map[uuid.UUID]*Player) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.players = players
	})
}

func WithVersion(version int) GameOption {
	return gameOptionFunc(func(g *gameOpts) {
		g.version = version
	})
}
