package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	fenPartsLength          = 6
	fenNoneSymbol           = "-"
	fenSeparator            = " "
	fenPositionSeparator    = "/"
	startingWhiteKingSquare = E1
	startingBlackKingSquare = E8
)

var (
	reIsDigit        = regexp.MustCompile(`^[0-9]$`)
	reEnpSquare      = regexp.MustCompile(`^(-|[abcdefgh][36])$`)
	reCastleRights   = regexp.MustCompile(`(?m)^(-|\bK?Q?k?q?)$`)
	reTurnColor      = regexp.MustCompile(`^(w|b)$`)
	reFenPieceSymbol = regexp.MustCompile(`^[prnbqkPRNBQK]$`)
)

type fenToken struct {
	halfMoveClock uint8
	fullMoveClock uint16
	enpSquare     Square
	castleRights  CastleRights
	turnColor     Color
	position      string
}

// validateFenMetadataParts validates the the parts after the position string token which are
// half and full move clocks, current turn color, en-passant square and castle rights
func validateFenMetadataParts(fen string, opts validateFenOps) (fenToken, error) {
	var emptyRet fenToken

	// tokens length must be 6 after splitting the fen by a single space delimiter
	tokens := strings.Split(fen, fenSeparator)
	if len(tokens) != fenPartsLength {
		return emptyRet, fmt.Errorf("invalid FEN: length must be exactly 6 after splitting by a single space delimiter")
	}

	var (
		fullMoveClockToken = tokens[5]
		halfMoveClockToken = tokens[4]
		enpSquareToken     = tokens[3]
		castleRightsToken  = tokens[2]
		turnColorToken     = tokens[1]
		positionToken      = tokens[0]
	)

	// turn color must be either `w` | `b`
	if !reTurnColor.MatchString(turnColorToken) {
		return emptyRet, fmt.Errorf("invalid FEN: invalid active turn color")
	}

	turn := White
	if turnColorToken == Black.String() {
		turn = Black
	}

	// full move clock must be a number >= 1
	fullMoveClock, err := strconv.ParseUint(fullMoveClockToken, 10, 16)
	if err != nil || fullMoveClock == 0 {
		return emptyRet, fmt.Errorf("invalid FEN: full move clock must be a number >= 1")
	}

	// half move clock must be a number >= 0
	halfMoveClock, err := strconv.ParseUint(halfMoveClockToken, 10, 8)
	if err != nil {
		return emptyRet, fmt.Errorf("invalid FEN: half move clock must be a number >= 0")
	}

	var n uint64
	if turn.IsBlack() {
		n = 1
	}

	// half move clock must be within the limit
	if !(halfMoveClock <= ((fullMoveClock-1)*2)+n) {
		return emptyRet, fmt.Errorf("invalid FEN: half move clock must be whithin the valid limit")
	}

	// in case of an en-passant square, the half move clock must be equal to 0
	if enpSquareToken != fenNoneSymbol && halfMoveClock != 0 {
		return emptyRet, fmt.Errorf("invalid FEN: half move clock must be 0 if en-passant square exists")
	}

	// en-passant square must be a valid square or `-` if empty
	if !reEnpSquare.MatchString(enpSquareToken) {
		return emptyRet, fmt.Errorf("invalid FEN: en-passant target square is invalid")
	}

	enpSquare := SquareNone
	if enpSquareToken != fenNoneSymbol {
		enpSquare = coordToSquare[enpSquareToken]
	}

	if enpSquare != SquareNone {
		if (turn.IsWhite() && enpSquare.Rank() == Rank3) || (turn.IsBlack() && enpSquare.Rank() == Rank6) {
			return emptyRet, fmt.Errorf("invalid FEN: en-passant target square coordinate is invalid")
		}
	}

	// castle rights string must be of valid fen castle string format
	if !reCastleRights.MatchString(castleRightsToken) {
		return emptyRet, fmt.Errorf("invalid FEN: invalid castling rights string")
	}

	var cr CastleRights
	if castleRightsToken != fenNoneSymbol {
		if strings.Contains(castleRightsToken, "K") {
			cr |= WhiteKingSideCastle
		}
		if strings.Contains(castleRightsToken, "Q") {
			cr |= WhiteQueenSideCastle
		}
		if strings.Contains(castleRightsToken, "k") {
			cr |= BlackKingSideCastle
		}
		if strings.Contains(castleRightsToken, "q") {
			cr |= BlackQueenSideCastle
		}
	}

	return fenToken{
		halfMoveClock: uint8(halfMoveClock),
		fullMoveClock: uint16(fullMoveClock),
		enpSquare:     enpSquare,
		castleRights:  cr,
		turnColor:     turn,
		position:      positionToken,
	}, nil
}

// validatePositionPart validates squares and pieces
func validatePositionPart(ft fenToken, opts validateFenOps) (map[Square]Piece, error) {
	ranks := strings.Split(ft.position, fenPositionSeparator)
	if len(ranks) != boardSize {
		return nil, fmt.Errorf("invalid FEN: it does not contain 8 ranks delimited by %q character", fenPositionSeparator)
	}

	piecesCount := make(map[Piece]uint8, 0)
	squares := make(map[Square]Piece, boardTotalSquares)

	for r := 0; r < len(ranks); r++ {
		var sumSquaresInRank uint8
		var previousWasNumber bool

		for f := 0; f < len(ranks[r]); f++ {
			if reIsDigit.MatchString(string(ranks[r][f])) {
				if previousWasNumber {
					return nil, fmt.Errorf("invalid FEN: position string is invalid, it has consecutive numbers")
				}

				n, err := strconv.ParseUint(string(ranks[r][f]), 10, 8)
				if err != nil {
					return nil, fmt.Errorf("invalid FEN: failed to parse row number")
				}

				sumSquaresInRank += uint8(n)
				previousWasNumber = true
			} else {
				piece, err := NewPieceFromFenSymbol(string(ranks[r][f]))
				if err != nil {
					return nil, fmt.Errorf("invalid FEN: position string contains invalid piece symbol")
				}

				sq := Square((7-r)*8 + int(sumSquaresInRank))
				squares[sq] = piece

				piecesCount[piece]++
				sumSquaresInRank++
				previousWasNumber = false
			}
		}

		if sumSquaresInRank != boardSize {
			return nil, fmt.Errorf("invalid FEN: position string is invalid, too many squares in rank")
		}
	}

	if piecesCount[WhiteKing] == 0 {
		return nil, fmt.Errorf("invalid FEN: position is missing white king")
	}
	if piecesCount[BlackKing] == 0 {
		return nil, fmt.Errorf("invalid FEN: position is missing black king")
	}
	if c := piecesCount[WhiteKing]; c > 1 {
		return nil, fmt.Errorf("invalid FEN: position is having too many white kings (%d)", c)
	}
	if c := piecesCount[BlackKing]; c > 1 {
		return nil, fmt.Errorf("invalid FEN: position is having too many black kings (%d)", c)
	}

	for _, char := range ranks[0] {
		if string(char) == WhitePawn.String() {
			return nil, fmt.Errorf("invalid FEN: white pawn is on 8th rank")
		}
	}
	for _, char := range ranks[7] {
		if string(char) == BlackPawn.String() {
			return nil, fmt.Errorf("invalid FEN: black pawn is on 1st rank")
		}
	}

	return squares, nil
}

type validateFenOps struct {
	strict bool
}

type positionMeta struct {
	fenToken
	squares map[Square]Piece
}

// validateFEN validates the fen string
func validateFEN(fen string, opts validateFenOps) (*positionMeta, error) {
	tkn, err := validateFenMetadataParts(fen, opts)
	if err != nil {
		return nil, err
	}

	squares, err := validatePositionPart(tkn, opts)
	if err != nil {
		return nil, err
	}

	meta := positionMeta{
		fenToken: tkn,
		squares:  squares,
	}

	return &meta, nil
}
