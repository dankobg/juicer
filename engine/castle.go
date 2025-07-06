package engine

import (
	"fmt"
	"strings"
)

type CastleRights uint8

const (
	WhiteKingSideCastle CastleRights = 1 << iota
	WhiteQueenSideCastle
	BlackKingSideCastle
	BlackQueenSideCastle
	CastleRightsNone CastleRights = 0
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

	if cr.whiteHasKingSideCastleRights() {
		sb.WriteString(wkCastleFen)
	}
	if cr.whiteHasQueenSideCastleRights() {
		sb.WriteString(wqCastleFen)
	}
	if cr.blackHasKingSideCastleRights() {
		sb.WriteString(bkCastleFen)
	}
	if cr.blackHasQueenSideCastleRights() {
		sb.WriteString(bqCastleFen)
	}

	return sb.String()
}

func (cr CastleRights) String() string {
	return cr.ToFEN()
}

func (cr CastleRights) whiteHasKingSideCastleRights() bool {
	return (cr & WhiteKingSideCastle) != 0
}

func (cr CastleRights) whiteHasQueenSideCastleRights() bool {
	return (cr & WhiteQueenSideCastle) != 0
}

func (cr CastleRights) whiteHasCastleRights() bool {
	return cr.whiteHasKingSideCastleRights() || cr.whiteHasQueenSideCastleRights()
}

func (cr CastleRights) blackHasKingSideCastleRights() bool {
	return (cr & BlackKingSideCastle) != 0
}

func (cr CastleRights) blackHasQueenSideCastleRights() bool {
	return (cr & BlackQueenSideCastle) != 0
}

func (cr CastleRights) blackHasCastleRights() bool {
	return cr.blackHasKingSideCastleRights() || cr.blackHasQueenSideCastleRights()
}

func (cr *CastleRights) disableWhiteCastleRights() {
	*cr &= ^(WhiteKingSideCastle | WhiteQueenSideCastle)
}

func (cr *CastleRights) disableWhiteKingSideCastleRight() {
	*cr &= ^WhiteKingSideCastle
}

func (cr *CastleRights) disableWhiteQueenSideCastleRight() {
	*cr &= ^WhiteQueenSideCastle
}

func (cr *CastleRights) disableBlackCastleRights() {
	*cr &= ^(BlackKingSideCastle | BlackQueenSideCastle)
}

func (cr *CastleRights) disableBlackKingSideCastleRight() {
	*cr &= ^BlackKingSideCastle
}

func (cr *CastleRights) disableBlackQueenSideCastleRight() {
	*cr &= ^BlackQueenSideCastle
}

func (cr *CastleRights) clear() {
	*cr = CastleRightsNone
}

func NewCastleRightsFromFen(fenCastle string) (CastleRights, error) {
	if reCastleRights.MatchString(fenCastle) {
		var cr CastleRights

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

		return cr, nil
	}

	return CastleRightsNone, fmt.Errorf("invalid castle rights string, doesn't match the pattern")
}
