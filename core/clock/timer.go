package clock

import (
	"time"
)

type TimerState int

const (
	TimerIdle TimerState = iota
	TimerRunning
	TimerPaused
)

type Timer struct {
	initialTime time.Duration
	time        time.Duration
	timer       *time.Timer
	startedAt   time.Time
	State       TimerState
}

func (ts TimerState) String() string {
	switch ts {
	case TimerIdle:
		return "idle"
	case TimerRunning:
		return "running"
	case TimerPaused:
		return "paused"
	}
	return ""
}

func NewTimer(dur time.Duration) *Timer {
	return &Timer{
		initialTime: dur,
		time:        dur,
		State:       TimerIdle,
	}
}

func (t *Timer) Start() {
	if t.State == TimerRunning {
		return
	}
	t.timer = time.NewTimer(t.time)
	t.startedAt = time.Now()
	t.State = TimerRunning
}

func (t *Timer) Pause() {
	if t.State == TimerRunning && t.timer != nil {
		t.time -= time.Since(t.startedAt)
		t.timer.Stop()
		t.State = TimerPaused
		t.startedAt = time.Now()
	}
}

func (t *Timer) Reset() {
	t.time = t.initialTime
	t.timer.Stop()
	t.timer = nil
	t.State = TimerIdle
}

func (t *Timer) Add(dur time.Duration) {
	if t.timer != nil {
		t.timer.Stop()
		t.time = t.Remaining() + dur
		t.timer = time.NewTimer(t.time)
		t.startedAt = time.Now()
	}
}

func (t *Timer) Remaining() time.Duration {
	if t.State == TimerRunning {
		return max(t.time-time.Since(t.startedAt), 0)
	} else {
		return t.time
	}
}

func (t *Timer) Tick() <-chan time.Time {
	return Tick(t.timer)
}

func Tick(t *time.Timer) <-chan time.Time {
	if t != nil {
		return t.C
	}
	return nil
}
