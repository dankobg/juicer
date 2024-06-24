package server

import (
	"time"
)

type timerAction int

const (
	timerActionStart timerAction = iota
	timerActionStop
)

func (ta timerAction) String() string {
	switch ta {
	case timerActionStart:
		return "start"
	case timerActionStop:
		return "stop"
	}
	return "unknown"
}

type timerKind int

const (
	timerKindDisconnect timerKind = iota
	timerKindFirstMove
)

func (tk timerKind) String() string {
	switch tk {
	case timerKindFirstMove:
		return "first_move"
	case timerKindDisconnect:
		return "disconnect"
	}
	return "unknown"
}

type player struct {
	id              string
	color           string
	disconnectTimer *time.Timer
	firstMoveTimer  *time.Timer
}

type timerEvent struct {
	playerID string
	action   timerAction
	kind     timerKind
}

func tick(t *time.Timer) <-chan time.Time {
	if t != nil {
		return t.C
	}
	return nil
}

func (p *player) startFirstMoveTimer(dur time.Duration) {
	if p.firstMoveTimer == nil {
		p.firstMoveTimer = time.NewTimer(dur)
	} else {
		p.firstMoveTimer.Reset(dur)
	}
}

func (p *player) stopFirstMoveTimer() {
	if p.firstMoveTimer != nil {
		p.firstMoveTimer.Stop()
		p.firstMoveTimer = nil
	}
}

func (p *player) startDisconnectTimer(dur time.Duration) {
	if p.disconnectTimer == nil {
		p.disconnectTimer = time.NewTimer(dur)
	} else {
		p.disconnectTimer.Reset(dur)
	}
}

func (p *player) stopDisconnectTimer() {
	if p.disconnectTimer != nil {
		p.disconnectTimer.Stop()
		p.disconnectTimer = nil
	}
}
