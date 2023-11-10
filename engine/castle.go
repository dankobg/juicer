package engine

type CastleRights int

const (
	WhiteKingSideCastle CastleRights = 1 << iota
	WhiteQueenSideCastle
	BlackKingSideCastle
	BlackQueenSideCastle
)

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
