package gameserver

import (
	"fmt"
	"juicer/engine"
	"time"
)

type Clock struct {
	white, black time.Timer

	duration  time.Duration
	remaining time.Duration
	timer     *time.Timer

	whiteMoved chan struct{}
	blacMoved  chan struct{}
}

func newClock(dur time.Duration) *Clock {
	return &Clock{
		duration:  dur,
		remaining: dur,
		timer:     time.NewTimer(dur),
	}
}

func (cl *Clock) Start() {
	<-cl.timer.C
	fmt.Println("Time's up!")
}

func (cl *Clock) Pause() {}

func (cl *Clock) Stop() {}

func (cl *Clock) TimeLeft() time.Duration {
	return time.Duration(300)
}
func (cl *Clock) WhiteTimeLeft() time.Duration {
	fmt.Println(cl.white)
	return time.Duration(300)
}
func (cl *Clock) BlackTimeLeft() time.Duration {
	fmt.Println(cl.black)
	return time.Duration(300)
}

func (cl *Clock) LastMoved() string {
	for {
		select {
		case <-cl.whiteMoved:
			return "w"
		case <-cl.blacMoved:
			return "b"
		}
	}
}

type GameState struct {
	Chess                  *engine.Chess
	White, Black           *engine.Player
	whiteClock, blackClock *Clock
	GameType               string
	GameMode               string
	ElloRating             int16
	GameStarted            bool
	GameFinished           bool
}

func NewGameState(p1, p2 *engine.Player) (*GameState, error) {
	c, err := engine.NewChess("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the game: %w", err)
	}

	wclock := newClock(5 * time.Minute)
	bclock := newClock(5 * time.Minute)

	gs := &GameState{
		Chess:      c,
		White:      p1,
		Black:      p2,
		whiteClock: wclock,
		blackClock: bclock,
	}
	return gs, nil
}
