package server

import juicer "github.com/dankobg/juicer/engine"

type Player struct {
	ID   string
	Name string
}

type GameState struct {
	Chess        *juicer.Chess
	White, Black *Player
	GameType     string
	GameMode     string
}

func NewGameState(white, black *Player) (*GameState, error) {
	chess, err := juicer.NewChess(juicer.FENStartingPosition)
	if err != nil {
		return nil, err
	}

	gs := &GameState{
		White:    white,
		Black:    black,
		GameType: "standard",
		GameMode: "blitz",
		Chess:    chess,
	}

	return gs, nil
}
