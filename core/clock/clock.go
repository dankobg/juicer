package clock

import (
	"time"
)

type ClockState uint8

const (
	ClockIdle ClockState = iota + 1
	ClockRunning
	ClockPaused
)

func (st ClockState) String() string {
	switch st {
	case ClockIdle:
		return "idle"
	case ClockRunning:
		return "running"
	case ClockPaused:
		return "paused"
	}
	return ""
}

type Color uint8

const (
	White Color = iota + 1
	Black
)

func (clr Color) String() string {
	switch clr {
	case White:
		return "white"
	case Black:
		return "black"
	}
	return ""
}

type Clock struct {
	InitialClock     time.Duration
	InitialIncrement time.Duration
	Increment        time.Duration
	White            *Timer
	Black            *Timer
	CurrentTurn      Color
	State            ClockState
}

func NewClock(clock, increment time.Duration) *Clock {
	return &Clock{
		InitialClock:     clock,
		InitialIncrement: increment,
		Increment:        increment,
		White:            NewTimer(clock),
		Black:            NewTimer(clock),
		CurrentTurn:      White,
		State:            ClockIdle,
	}
}

func (c *Clock) currentTimer() *Timer {
	if c.CurrentTurn == White {
		return c.White
	}
	return c.Black
}

func (c *Clock) toggleCurrentTurn() {
	if c.CurrentTurn == White {
		c.SetCurrentTurn(Black)
	} else {
		c.SetCurrentTurn(White)
	}
}

func (c *Clock) SetCurrentTurn(color Color) {
	c.CurrentTurn = color
}

func (c *Clock) Start() {
	if c.State == ClockRunning {
		return
	}
	c.State = ClockRunning
	c.currentTimer().Start()
}

func (c *Clock) Pause() {
	if c.State == ClockRunning {
		c.State = ClockPaused
		c.currentTimer().Pause()
	}
}

func (c *Clock) Reset() {
	c.State = ClockIdle
	c.Increment = c.InitialIncrement
	c.SetCurrentTurn(White)
	c.White.Reset()
	c.Black.Reset()
}

func (c *Clock) Restart() {
	c.Reset()
	c.Start()
}

func (c *Clock) Toggle() {
	if c.State == ClockIdle {
		c.Start()
	} else if c.State == ClockRunning {
		c.currentTimer().Pause()
		if c.Increment > 0 {
			c.currentTimer().Add(c.Increment)
		}
		c.toggleCurrentTurn()
		c.currentTimer().Start()
	}
}
