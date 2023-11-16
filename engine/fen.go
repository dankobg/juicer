package juicer

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
	reIsDigit        = regexp.MustCompile("^[0-9]$")
	reEnpSquare      = regexp.MustCompile("^(-|[abcdefgh][36])$")
	reCastleRights   = regexp.MustCompile("[^kKqQ-]")
	reTurnColor      = regexp.MustCompile("^(w|b)$")
	reFenPieceSymbol = regexp.MustCompile("^[prnbqkPRNBQK]$")
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
	noopRet := fenToken{}

	// tokens length must be 6 after splitting the fen by a single space delimiter
	tokens := strings.Split(fen, fenSeparator)
	if len(tokens) != fenPartsLength {
		return noopRet, fmt.Errorf("invalid FEN: length must be exactly 6 after splitting by a single space delimiter")
	}

	var (
		halfMoveClockToken = tokens[5]
		fullMoveClockToken = tokens[4]
		enpSquareToken     = tokens[3]
		castleRightsToken  = tokens[2]
		turnColorToken     = tokens[1]
		positionToken      = tokens[0]
	)

	// turn color must be either `w` | `b`
	if !reTurnColor.MatchString(turnColorToken) {
		return noopRet, fmt.Errorf("invalid FEN: invalid active turn color")
	}

	// full move clock must be a number >= 1
	fullMoveClock, err := strconv.ParseUint(fullMoveClockToken, 10, 8)
	if err != nil || fullMoveClock == 0 {
		return noopRet, fmt.Errorf("invalid FEN: full move clock must be a number >= 1")
	}

	// half move clock must be a number >= 0
	halfMoveClock, err := strconv.ParseUint(halfMoveClockToken, 10, 64)
	if err != nil {
		return noopRet, fmt.Errorf("invalid FEN: half move clock must be a number >= 0")
	}

	var n uint64
	if turnColorToken == Black.String() {
		n = 1
	}

	// half move clock must be within the limit
	if !(halfMoveClock <= ((fullMoveClock-1)*2)+n) {
		return noopRet, fmt.Errorf("invalid FEN: half move clock must be whithin the valid limit")
	}

	// in case of an en-passant square, the half move clock must be equal to 0
	if enpSquareToken != fenNoneSymbol && halfMoveClock != 0 {
		return noopRet, fmt.Errorf("invalid FEN: half move clock must be 0 if en-passant square exists")
	}

	// en-passant square must be a valid square or `-` if empty
	if !reEnpSquare.MatchString(enpSquareToken) {
		return noopRet, fmt.Errorf("invalid FEN: en-passant target square is invalid")
	}

	// castle rights string must be of valid fen castle string format
	if reCastleRights.MatchString(castleRightsToken) {
		return noopRet, fmt.Errorf("invalid FEN: invalid castling rights string")
	}

	turn := White
	if turnColorToken == Black.String() {
		turn = Black
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

	enpSquare := SquareNone
	if enpSquareToken != fenNoneSymbol {
		enpSquare = coordToSquare[enpSquareToken]
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
	// position string contains 8 ranks
	ranks := strings.Split(ft.position, fenPositionSeparator)
	if len(ranks) != boardSize {
		return nil, fmt.Errorf("invalid FEN: it does not contain 8 ranks delimited by %q character", fenPositionSeparator)
	}

	squares := make(map[Square]Piece, boardTotalSquares)
	var whiteKingsCount, blackKingsCount int
	whiteKingSquare, blackKingSquare := SquareNone, SquareNone

	for r := 0; r < len(ranks); r++ {
		var sumSquaresInRank int64
		var previousWasNumber bool

		for f := 0; f < len(ranks[f]); f++ {
			if reIsDigit.MatchString(string(ranks[r][f])) {
				if previousWasNumber {
					return nil, fmt.Errorf("invalid FEN: position string is invalid, it has consecutive numbers")
				}

				n, err := strconv.ParseInt(string(ranks[r][f]), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid FEN: failed to parse row number")
				}

				sumSquaresInRank += n
				previousWasNumber = true
			} else {
				symbol := string(ranks[r][f])
				if !reFenPieceSymbol.MatchString(symbol) {
					return nil, fmt.Errorf("invalid FEN: position string contains invalid piece symbol")
				}

				piece, err := NewPieceFromFenSymbol(symbol)
				if err != nil {
					return nil, fmt.Errorf("invalid FEN: position string contains invalid piece symbol")
				}

				sqIdx := r*8 + f
				sq := Square(sqIdx)
				if sq.IndexInBoard() {
					squares[sq] = piece
				}

				if piece.IsKing() {
					if piece.Color().IsWhite() {
						whiteKingsCount++
						whiteKingSquare = sq
					}
					if piece.Color().IsBlack() {
						blackKingsCount++
						blackKingSquare = sq
					}
				}

				sumSquaresInRank += 1
				previousWasNumber = false
			}
		}

		if sumSquaresInRank != boardSize {
			return nil, fmt.Errorf("invalid FEN: position string is invalid, too many squares in rank")
		}
	}

	if ft.enpSquare != SquareNone {
		if (ft.enpSquare.Rank() == Rank3 && ft.turnColor == White) || (ft.enpSquare.Rank() == Rank6 && ft.turnColor == Black) {
			return nil, fmt.Errorf("invalid FEN: illegal en-passant target square")
		}
	}

	if ft.turnColor == White && whiteKingSquare != startingWhiteKingSquare {
		return nil, fmt.Errorf("invalid FEN: white king is not on a starting position which conflicts with the castle string")
		// @TODO: check if k/q side rook is not on starting square and if castle rights says you can castle on that k/q side
	}

	if ft.turnColor == Black && blackKingSquare != startingBlackKingSquare {
		return nil, fmt.Errorf("invalid FEN: black king is not on a starting position which conflicts with the castle string")
		// @TODO: check if k/q side rook is not on starting square and if castle rights says you can castle on that k/q side
	}

	if whiteKingsCount == 0 {
		return nil, fmt.Errorf("invalid FEN: position is missing white king")
	}
	if blackKingsCount == 0 {
		return nil, fmt.Errorf("invalid FEN: position is missing black king")
	}
	if whiteKingsCount > 1 {
		return nil, fmt.Errorf("invalid FEN: position is having too many white kings (%d)", whiteKingsCount)
	}
	if blackKingsCount > 1 {
		return nil, fmt.Errorf("invalid FEN: position is having too many black kings (%d)", blackKingsCount)
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

func loadPositionFromFEN(fen string) (*Position, error) {
	meta, err := validateFEN(fen, validateFenOps{})
	if err != nil {
		return nil, err
	}

	board := NewEmptyBoard()

	for sq, piece := range meta.squares {
		if piece == WhiteKing {
			board.whiteKingOccupancy.setBit(sq)
		} else if piece == WhiteQueen {
			board.whiteKingOccupancy.setBit(sq)
		} else if piece == WhiteRook {
			board.whiteRooksOccupancy.setBit(sq)
		} else if piece == WhiteBishop {
			board.whiteBishopsOccupancy.setBit(sq)
		} else if piece == WhiteKnight {
			board.whiteKnightsOccupancy.setBit(sq)
		} else if piece == WhitePawn {
			board.whitePawnsOccupancy.setBit(sq)
		} else if piece == BlackKing {
			board.blackKingOccupancy.setBit(sq)
		} else if piece == BlackQueen {
			board.blackQueensOccupancy.setBit(sq)
		} else if piece == BlackRook {
			board.blackRooksOccupancy.setBit(sq)
		} else if piece == BlackBishop {
			board.blackBishopsOccupancy.setBit(sq)
		} else if piece == BlackKnight {
			board.blackKnightsOccupancy.setBit(sq)
		} else if piece == BlackPawn {
			board.blackPawnsOccupancy.setBit(sq)
		}
	}

	pos := Position{
		board:         &board,
		turn:          meta.turnColor,
		enpSquare:     meta.enpSquare,
		castleRights:  meta.castleRights,
		halfMoveClock: meta.halfMoveClock,
		fullMoveClock: meta.fullMoveClock,
		// check:                false,
		// checkmate:            false,
		// stalemate:            false,
		// draw:                 false,
		// threeFold:            false,
		// insufficientMaterial: false,
		// terminated:           false,
		// outcome:              "",
		// comments:             []string{},
		// headers:              []string{},
		// capturedPieces:       []Piece{},
		// alivePieces:          []Piece{},
	}

	return &pos, nil
}
