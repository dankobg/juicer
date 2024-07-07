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

func (gm gameMode) String() string {
	switch gm {
	case gameModeClassical:
		return "classical"
	case gameModeRapid:
		return "rapid"
	case gameModeBlitz:
		return "blitz"
	case gameModeBullet:
		return "bullet"
	case gameModeHyperBullet:
		return "hyperbullet"
	default:
		return "unknown"
	}
}

type result uint8

const (
	resultWhiteWon result = iota
	resultBlackWon
	resultDraw
	resultAborted
)

func (r result) String() string {
	switch r {
	case resultWhiteWon:
		return "white won"
	case resultBlackWon:
		return "black won"
	case resultDraw:
		return "draw"
	case resultAborted:
		return "aborted"
	default:
		return "unknown"
	}
}

type resultStatus uint8

const (
	resultStatusUnknown resultStatus = iota
	resultStatusCheckmate
	resultStatusInsufficientMaterial
	resultStatusThreefoldRepetition
	resultStatusFivefoldRepetition
	resultStatusFiftyMoveRule
	resultStatusSeventyFiveMoveRule
	resultStatusStalemate
	resultStatusResignation
	resultStatusDrawAgreed
	resultStatusFlagged
	resultStatusAdjudication
	resultStatusTimedOut
	resultStatusAborted
)

func (rs resultStatus) String() string {
	switch rs {
	case resultStatusUnknown:
		return "unknown"
	case resultStatusCheckmate:
		return "checkmate"
	case resultStatusInsufficientMaterial:
		return "insufficient material"
	case resultStatusThreefoldRepetition:
		return "threefold repetition"
	case resultStatusFivefoldRepetition:
		return "fivefold repetition"
	case resultStatusFiftyMoveRule:
		return "fifty move rule"
	case resultStatusSeventyFiveMoveRule:
		return "seventy five move rule"
	case resultStatusStalemate:
		return "stalemate"
	case resultStatusResignation:
		return "resignation"
	case resultStatusDrawAgreed:
		return "draw agreed"
	case resultStatusFlagged:
		return "flagged"
	case resultStatusAdjudication:
		return "adjudication"
	case resultStatusTimedOut:
		return "timed out"
	case resultStatusAborted:
		return "aborted"
	default:
		return "unknown"
	}
}

type MatchState uint8

const (
	matchStateIdle MatchState = iota
	matchStateWaitingStart
	matchStateStarted
	matchStateFinished
	matchStateAborted
)

type gameState struct {
	GameID             string
	Chess              *juicer.Chess
	WhiteID            string
	BlackID            string
	GameType           gameType
	GameMode           gameMode
	Result             result
	ResultStatus       resultStatus
	MatchState         MatchState
	StartTime          time.Time
	LastMove           time.Time
	RemainingTimeWhite time.Duration
	RemainingTimeBlack time.Duration
}

func NewGameState(whiteID, blackID string, gameType gameType, gameMode gameMode) (*gameState, error) {
	gameID := random.AlphaNumeric(32)

	chess, err := juicer.NewChess(juicer.FENStartingPosition)
	if err != nil {
		return nil, err
	}

	gs := &gameState{
		GameID:             gameID,
		Chess:              chess,
		WhiteID:            whiteID,
		BlackID:            blackID,
		GameType:           gameType,
		GameMode:           gameMode,
		MatchState:         matchStateWaitingStart,
		ResultStatus:       resultStatusUnknown,
		StartTime:          time.Now(),
		RemainingTimeWhite: time.Minute * 5,
		RemainingTimeBlack: time.Minute * 5,
	}

	return gs, nil
}

func (gs *gameState) updatePlayerClockAfterMove() {
	now := time.Now()
	var elapsed time.Duration
	if !gs.LastMove.IsZero() {
		elapsed = now.Sub(gs.LastMove)
	}

	if gs.Chess.Position.Turn.IsWhite() {
		gs.RemainingTimeBlack -= elapsed
	} else if gs.Chess.Position.Turn.IsBlack() {
		gs.RemainingTimeWhite -= elapsed
	}

	gs.LastMove = now
}

func (gs *gameState) updatePlayerClockOnRejoin() {
	now := time.Now()
	var elapsed time.Duration
	if !gs.LastMove.IsZero() {
		elapsed = now.Sub(gs.LastMove)
	}

	if gs.Chess.Position.Turn.IsWhite() {
		gs.RemainingTimeWhite -= elapsed
	} else if gs.Chess.Position.Turn.IsBlack() {
		gs.RemainingTimeBlack -= elapsed
	}

	gs.LastMove = now
}
