package juicer

import (
	"fmt"
	"strings"
)

type CastleRights uint8

const (
	CastleRightsNone    CastleRights = 0
	WhiteKingSideCastle CastleRights = 1 << iota
	WhiteQueenSideCastle
	BlackKingSideCastle
	BlackQueenSideCastle
)

const (
	wkCastleFen = "K"
	wqCastleFen = "Q"
	bkCastleFen = "k"
	bqCastleFen = "q"
)

func (cr CastleRights) ToFEN() string {
	if cr == CastleRightsNone {
		return fenNoneSymbol
	}

	var sb strings.Builder

	if cr.whiteCanCastleKingSide() {
		sb.WriteString(wkCastleFen)
	}
	if cr.whiteCanCastleQueenSide() {
		sb.WriteString(wqCastleFen)
	}
	if cr.blackCanCastleKingSide() {
		sb.WriteString(bkCastleFen)
	}
	if cr.blackCanCastleQueenSide() {
		sb.WriteString(bqCastleFen)
	}

	return sb.String()
}

func (cr CastleRights) String() string {
	return cr.ToFEN()
}

func (cr CastleRights) whiteCanCastleKingSide() bool {
	return (cr & WhiteKingSideCastle) != 0
}

func (cr CastleRights) whiteCanCastleQueenSide() bool {
	return (cr & WhiteQueenSideCastle) != 0
}

func (cr CastleRights) whiteCanCastle() bool {
	return cr.whiteCanCastleKingSide() || cr.whiteCanCastleQueenSide()
}

func (cr CastleRights) blackCanCastleKingSide() bool {
	return (cr & BlackKingSideCastle) != 0
}

func (cr CastleRights) blackCanCastleQueenSide() bool {
	return (cr & BlackQueenSideCastle) != 0
}

func (cr CastleRights) blackCanCastle() bool {
	return cr.blackCanCastleKingSide() || cr.blackCanCastleQueenSide()
}

func (cr *CastleRights) preventWhiteFromCastling() {
	*cr &= ^(WhiteKingSideCastle | WhiteQueenSideCastle)
}

func (cr *CastleRights) preventWhiteFromCastlingKingSide() {
	*cr &= ^WhiteKingSideCastle
}

func (cr *CastleRights) preventWhiteFromCastlingQueenSide() {
	*cr &= ^WhiteQueenSideCastle
}

func (cr *CastleRights) preventBlackFromCastling() {
	*cr &= ^(BlackKingSideCastle | BlackQueenSideCastle)
}

func (cr *CastleRights) preventBlackFromCastlingKingSide() {
	*cr &= ^BlackKingSideCastle
}

func (cr *CastleRights) preventBlackFromCastlingQueenSide() {
	*cr &= ^BlackQueenSideCastle
}

func (cr *CastleRights) clear() {
	*cr = CastleRightsNone
}

func NewCastleRightsFromFen(fenCastle string) (CastleRights, error) {
	if reCastleRights.MatchString(fenCastle) {
		var cr CastleRights

		if fenCastle != fenNoneSymbol {
			if strings.Contains(fenCastle, wkCastleFen) {
				cr |= WhiteKingSideCastle
			}
			if strings.Contains(fenCastle, wqCastleFen) {
				cr |= WhiteQueenSideCastle
			}
			if strings.Contains(fenCastle, bkCastleFen) {
				cr |= BlackKingSideCastle
			}
			if strings.Contains(fenCastle, bqCastleFen) {
				cr |= BlackQueenSideCastle
			}
		}
	}

	return CastleRightsNone, fmt.Errorf("invalid castle rights string, doesn't match the pattern")
}
