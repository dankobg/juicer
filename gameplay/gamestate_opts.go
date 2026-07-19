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

type GameOptions []GameOption

func (o GameOptions) apply(g *gameOpts) {
	for _, opt := range o {
		opt.apply(g)
	}
}

type gameIDOpt int64

func (o gameIDOpt) apply(g *gameOpts)    { g.gameID = int64(o) }
func WithGameID(gameID int64) GameOption { return gameIDOpt(gameID) }

type fenOpt string

func (o fenOpt) apply(g *gameOpts)  { g.fen = string(o) }
func WithFEN(fen string) GameOption { return fenOpt(fen) }

type gameVariantOpt pb.GameVariant

func (o gameVariantOpt) apply(g *gameOpts)                  { g.gameVariant = pb.GameVariant(o) }
func WithGameVariant(gameVariant pb.GameVariant) GameOption { return gameVariantOpt(gameVariant) }

type gameTimeKindOpt pb.GameTimeKind

func (o gameTimeKindOpt) apply(g *gameOpts)                    { g.gameTimeKind = pb.GameTimeKind(o) }
func WithGameTimeKind(gameTimeKind pb.GameTimeKind) GameOption { return gameTimeKindOpt(gameTimeKind) }

type gameTimeCategoryOpt pb.GameTimeCategory

func (o gameTimeCategoryOpt) apply(g *gameOpts) { g.gameTimeCategory = pb.GameTimeCategory(o) }
func WithGameTimeCategory(gameTimeCategory pb.GameTimeCategory) GameOption {
	return gameTimeCategoryOpt(gameTimeCategory)
}

type gameTimeControlOpt struct{ tc *pb.GameTimeControl }

func (o gameTimeControlOpt) apply(g *gameOpts) { g.gameTimeControl = o.tc }
func WithGameTimeControl(gameTimeControl *pb.GameTimeControl) GameOption {
	return gameTimeControlOpt{tc: gameTimeControl}
}

type gameResultOpt pb.GameResult

func (o gameResultOpt) apply(g *gameOpts) { g.gameResult = pb.GameResult(o) }
func WithGameResult(gameResult pb.GameResult) GameOption {
	return gameResultOpt(gameResult)
}

type gameResultStatusOpt pb.GameResultStatus

func (o gameResultStatusOpt) apply(g *gameOpts) { g.gameResultStatus = pb.GameResultStatus(o) }
func WithGameResultStatus(gameResultStatus pb.GameResultStatus) GameOption {
	return gameResultStatusOpt(gameResultStatus)
}

type gameStateOpt pb.GameState

func (o gameStateOpt) apply(g *gameOpts) { g.gameState = pb.GameState(o) }
func WithGameState(gameState pb.GameState) GameOption {
	return gameStateOpt(gameState)
}

type reconnectTimeoutOpt time.Duration

func (o reconnectTimeoutOpt) apply(g *gameOpts) { g.reconnectTimeout = time.Duration(o) }
func WithReconnectTimeout(reconnectTimeout time.Duration) GameOption {
	return reconnectTimeoutOpt(reconnectTimeout)
}

type firstMoveTimeoutOpt time.Duration

func (o firstMoveTimeoutOpt) apply(g *gameOpts) { g.firstMoveTimeout = time.Duration(o) }
func WithFirstMoveTimeoutOpt(firstMoveTimeout time.Duration) GameOption {
	return firstMoveTimeoutOpt(firstMoveTimeout)
}

type ratedOpt bool

func (o ratedOpt) apply(g *gameOpts)  { g.rated = bool(o) }
func WithRated(rated bool) GameOption { return ratedOpt(rated) }

type lastMoveOpt struct{ t *time.Time }

func (o lastMoveOpt) apply(g *gameOpts)           { g.lastMove = o.t }
func WithLastMove(lastMove *time.Time) GameOption { return lastMoveOpt{t: lastMove} }

type startTimeOpt struct{ t *time.Time }

func (o startTimeOpt) apply(g *gameOpts)            { g.startTime = o.t }
func WithStartTime(startTime *time.Time) GameOption { return startTimeOpt{t: startTime} }

type endTimeOpt struct{ t *time.Time }

func (o endTimeOpt) apply(g *gameOpts)          { g.endTime = o.t }
func WithEndTime(endTime *time.Time) GameOption { return endTimeOpt{t: endTime} }

type gameMovesOpt struct{ moves []*pb.GameMove }

func (o gameMovesOpt) apply(g *gameOpts)            { g.gameMoves = o.moves }
func WithGameMoves(moves []*pb.GameMove) GameOption { return gameMovesOpt{moves: moves} }

type whitePlayerOpt struct{ p *Player }

func (o whitePlayerOpt) apply(g *gameOpts) { g.white = o.p }
func WithWhitePlayer(p *Player) GameOption { return whitePlayerOpt{p: p} }

type blackPlayerOpt struct{ p *Player }

func (o blackPlayerOpt) apply(g *gameOpts) { g.white = o.p }
func WithBlackPlayer(p *Player) GameOption { return blackPlayerOpt{p: p} }

type playersOpt struct{ players map[uuid.UUID]*Player }

func (o playersOpt) apply(g *gameOpts)                     { g.players = o.players }
func WithPlayers(players map[uuid.UUID]*Player) GameOption { return playersOpt{players: players} }

type versionOpt int

func (o versionOpt) apply(g *gameOpts) { g.version = int(o) }
func WithVersion(rated int) GameOption { return versionOpt(rated) }
