package server

import juicer "github.com/dankobg/juicer/engine"

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
	Chess        *juicer.Chess
	White, Black *Player
	GameType     GameType
	GameMode     GameMode
	State        State
}

func NewGameState(white, black *Player, gameType GameType, gameMode GameMode) (*GameState, error) {
	chess, err := juicer.NewChess(juicer.FENStartingPosition)
	if err != nil {
		return nil, err
	}

	gs := &GameState{
		Chess:    chess,
		White:    white,
		Black:    black,
		GameType: gameType,
		GameMode: gameMode,
		State:    StateIdle,
	}

	return gs, nil
}
