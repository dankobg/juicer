package core

import (
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
)

type gameOpts struct {
	fen          string
	variant      pb.Variant
	timeKind     pb.GameTimeKind
	timeCategory pb.GameTimeCategory
	timeControl  *pb.GameTimeControl
	// result           pb.GameResult
	// resultStatus     pb.GameResultStatus
	// state            pb.GameState
	reconnectTimeout time.Duration
	firstMoveTimeout time.Duration
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

type fenOpt string

func (o fenOpt) apply(g *gameOpts)  { g.fen = string(o) }
func WithFEN(fen string) GameOption { return fenOpt(fen) }

type variantOpt pb.Variant

func (o variantOpt) apply(g *gameOpts)          { g.variant = pb.Variant(o) }
func WithVariant(variant pb.Variant) GameOption { return variantOpt(variant) }

type timeKindOpt pb.GameTimeKind

func (o timeKindOpt) apply(g *gameOpts)                { g.timeKind = pb.GameTimeKind(o) }
func WithTimeKind(timeKind pb.GameTimeKind) GameOption { return timeKindOpt(timeKind) }

type timeCategoryOpt pb.GameTimeCategory

func (o timeCategoryOpt) apply(g *gameOpts) { g.timeCategory = pb.GameTimeCategory(o) }
func WithTimeCategory(timeCategory pb.GameTimeCategory) GameOption {
	return timeCategoryOpt(timeCategory)
}

type timeControlOpt struct{ tc *pb.GameTimeControl }

func (o timeControlOpt) apply(g *gameOpts) { g.timeControl = o.tc }
func WithTimeControl(timeControl *pb.GameTimeControl) GameOption {
	return timeControlOpt{tc: timeControl}
}

// type resultOpt pb.GameResult

// func (o resultOpt) apply(g *gameOpts) { g.result = pb.GameResult(o) }
// func WithGameResult(result pb.GameResult) GameOption {
// 	return resultOpt(result)
// }

// type resultStatusOpt pb.GameResultStatus

// func (o resultStatusOpt) apply(g *gameOpts) { g.resultStatus = pb.GameResultStatus(o) }
// func WithGameResultStatus(resultStatus pb.GameResultStatus) GameOption {
// 	return resultStatusOpt(resultStatus)
// }

// type stateOpt pb.GameState

// func (o stateOpt) apply(g *gameOpts) { g.state = pb.GameState(o) }
// func WithGameState(state pb.GameState) GameOption {
// 	return stateOpt(state)
// }

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
