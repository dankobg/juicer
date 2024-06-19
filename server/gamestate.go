package server

import (
	"time"

	juicer "github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/random"
)

type gameType uint8

const (
	gameTypeStandard gameType = iota
	gameTypeChess960
)

type gameMode uint8

const (
	gameModeClassical gameMode = iota
	gameModeRapid
	gameModeBlitz
	gameModeBullet
	gameModeHyperBullet
)

type result uint8

const (
	resultWhiteWon result = iota
	resultBlackWon
	resultDraw
)

type ResultStatus uint8

const (
	resultStatusUnknown ResultStatus = iota
	resultStatusCheckmate
	resultStatusInsufficientMaterial
	resultStatusThreeFoldRepetition
	resultStatusFiveFoldRepetition
	resultStatusFiftyMoveRule
	resultStatusSeventyFiveMoveRule
	resultStatusStalemate
	resultStatusResignation
	resultStatusDrawAgreed
	resultStatusFlagged
	resultStatusAdjudication
	resultStatusTimedOut
)

type MatchState uint8

const (
	matchStateIdle MatchState = iota
	matchStateWaitingStart
	matchStateStarted
	matchStateFinished
	matchStateAborted
)

type player struct {
	ID   string
	Name string
}

type gameState struct {
	GameID       string
	Chess        *juicer.Chess
	White, Black *player
	GameType     gameType
	GameMode     gameMode
	Result       result
	ResultStatus ResultStatus
	MatchState   MatchState
	StartTime    time.Time
}

func NewGameState(white, black *player, gameType gameType, gameMode gameMode) (*gameState, error) {
	gameID := random.AlphaNumeric(32)

	chess, err := juicer.NewChess(juicer.FENStartingPosition)
	if err != nil {
		return nil, err
	}

	gs := &gameState{
		GameID:       gameID,
		Chess:        chess,
		White:        white,
		Black:        black,
		GameType:     gameType,
		GameMode:     gameMode,
		MatchState:   matchStateWaitingStart,
		ResultStatus: resultStatusUnknown,
	}

	return gs, nil
}
