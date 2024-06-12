package server

import (
	juicer "github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/random"
)

type GameType uint8

const (
	GameTypeStandard GameType = iota
	GameTypeChess960
)

type GameMode uint8

const (
	GameModeClassical GameMode = iota
	GameModeRapid
	GameModeBlitz
	GameModeBullet
	GameModeHyperBullet
)

type State uint8

const (
	StateIdle State = iota
	StateStarted
	StateFinished
)

type Player struct {
	ID   string
	Name string
}

type GameState struct {
	GameID       string
	Chess        *juicer.Chess
	White, Black *Player
	GameType     GameType
	GameMode     GameMode
	State        State
}

func NewGameState(white, black *Player, gameType GameType, gameMode GameMode) (*GameState, error) {
	gameID := random.AlphaNumeric(32)

	chess, err := juicer.NewChess(juicer.FENStartingPosition)
	if err != nil {
		return nil, err
	}

	gs := &GameState{
		GameID:   gameID,
		Chess:    chess,
		White:    white,
		Black:    black,
		GameType: gameType,
		GameMode: gameMode,
		State:    StateIdle,
	}

	return gs, nil
}
