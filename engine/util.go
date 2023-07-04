package engine

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Metadata struct {
	board          [][]*Square
	turn           Color
	castlingRights CastlingRights
	enpSquare      *Square
	halfMoves      uint8
	fullMoves      uint8
}

const (
	fenLength     = 6
	fenNoneSymbol = "-"

	whiteKing = "K"
	blackKing = "k"
)

var (
	reIsDigit      = regexp.MustCompile("^[0-9]$")
	reEnpSquare    = regexp.MustCompile("^(-|[abcdefgh][36])$")
	reCastleRights = regexp.MustCompile("[^kKqQ-]")
	reTurnColor    = regexp.MustCompile("^(w|b)$")
	reFenSymbol    = regexp.MustCompile("^[prnbqkPRNBQK]$")
)

func isDigit(s string) bool {
	return reIsDigit.Match([]byte(s))
}

func validateFEN(FEN string) error {
	// 1. criterion: fen length must be 6 after splitting the fen by a single space delimiter
	tokens := strings.Split(FEN, " ")
	if len(tokens) != fenLength {
		return errors.New("invalid FEN: length must be exactly 6 after splitting by a single space delimiter")
	}

	// 2. criterion: full clock number is a number >= 1
	fullMoveClock, err := strconv.ParseUint(tokens[5], 10, 8)
	if err != nil || fullMoveClock == 0 {
		return errors.New("invalid FEN: full move clock must be a number >= 1")
	}

	// 3. criterion: half move clock is a number >= 0
	_, err = strconv.ParseUint(tokens[4], 10, 8)
	if err != nil {
		return errors.New("invalid FEN: half move clock must be a number >= 0")
	}

	// 4. criterion: en-passant target square is a vaid square coordinate or `-` if empty
	if !reEnpSquare.Match([]byte(tokens[3])) {
		return errors.New("invalid FEN: en-passant target square is invalid")
	}

	// 5. criterion: castling rights is a valid castle-string
	if reCastleRights.Match([]byte(tokens[2])) {
		return errors.New("invalid FEN: invalid castling rights string")
	}

	// 6. criterion: active turn color must be `w` (white) or `b` (black)
	if !reTurnColor.Match([]byte(tokens[1])) {
		return errors.New("invalid FEN: invalid active turn color")
	}

	// 7. criterion: first field contains 8 rows
	rows := strings.Split(tokens[0], "/")
	if len(rows) != boardSize {
		return errors.New("invalid FEN: it does not contain 8 `/` delimited rows")
	}

	// 8. criterion: every row is valid
	for r := 0; r < len(rows); r++ {
		var sumFields int64
		previousWasNumber := false

		for c := 0; c < len(rows[r]); c++ {
			if isDigit(string(rows[r][c])) {
				if previousWasNumber {
					return errors.New("invalid FEN: piece data is invalid (consecutive number)")
				}

				n, err := strconv.ParseInt(string(rows[r][c]), 10, 64)
				if err != nil {
					return errors.New("invalid FEN: failed to parse row number")
				}
				sumFields += n
				previousWasNumber = true
			} else {
				if !reFenSymbol.Match([]byte(string(rows[r][c]))) {
					return errors.New("invalid FEN: piece data is invalid, invalid piece")
				}

				sumFields += 1
				previousWasNumber = false
			}
		}

		if sumFields != boardSize {
			return errors.New("invalid FEN: piece data is invalid, too many squares in rank")
		}

		if string(tokens[3][0]) != fenNoneSymbol {
			if (string(tokens[3][1]) == "3" && tokens[1] == White.String()) || (string(tokens[3][1]) == "6" && tokens[1] == Black.String()) {
				return errors.New("invalid FEN: illegal en-passant target square")
			}
		}

		var whiteKings, blackKings int

		if !strings.Contains(string(tokens[0]), whiteKing) {
			return errors.New("invalid FEN: missing white king")
		}
		if !strings.Contains(string(tokens[0]), blackKing) {
			return errors.New("invalid FEN: missing black king")
		}

		for _, char := range tokens[0] {
			if string(char) == whiteKing {
				whiteKings++
			}
			if string(char) == blackKing {
				blackKings++
			}
		}

		if whiteKings > 1 {
			return errors.New("invalid FEN: too many white kings")
		}
		if blackKings > 1 {
			return errors.New("invalid FEN: too many black kings")
		}
	}

	return nil
}

func convertMetadataToFEN(meta Metadata) string {
	var fen string

	squares := meta.board

	for r := 0; r < len(squares); r++ {
		var empty int

		for c := 0; c < len(squares[r]); c++ {
			if squares[r][c].IsEmpty() {
				empty += 1
			}

			if squares[r][c].HasPiece() {
				if empty > 0 {
					fen += fmt.Sprintf("%v", empty)
					empty = 0
				}

				fen += squares[r][c].piece.ToFENSymbol().String()
			}
		}

		if empty > 0 {
			fen += fmt.Sprint(empty)
		}

		if r < len(squares)-1 {
			fen += "/"
		}
	}

	var castleToken string

	if meta.castlingRights == 0 {
		castleToken = "-"
	} else {
		if (meta.castlingRights & WhiteKingSideCastle) > 0 {
			castleToken += "K"
		}
		if (meta.castlingRights & WhiteQueenSideCastle) > 0 {
			castleToken += "Q"
		}
		if (meta.castlingRights & BlackKingSideCastle) > 0 {
			castleToken += "k"
		}
		if (meta.castlingRights & BlackQueenSideCastle) > 0 {
			castleToken += "q"
		}
	}

	enpSquareToken := "-"
	if meta.enpSquare != nil {
		enpSquareToken += meta.enpSquare.Coordinate().String()
	}

	finalFenTokenPart := fmt.Sprintf(" %s %s %s %v %v", meta.turn, castleToken, enpSquareToken, meta.halfMoves, meta.fullMoves)
	fen += finalFenTokenPart

	return fen
}

func convertFENPositionsTokenToBoard(fen string) [][]*Square {
	squares := initBoardSquares()

	for r, fenRow := range strings.Split(fen, "/") {
		var c int

		for _, char := range fenRow {
			if char >= '1' && char <= '8' {
				c += int(char - '0')
			} else {
				squares[r][c].piece = NewPieceFromFENSymbol(PieceFENSymbol(char), true)
				c++
			}
		}
	}

	return squares
}

func ParseMetadataFromFEN(fen string) (*Metadata, error) {
	if err := validateFEN(fen); err != nil {
		return nil, err
	}

	tokens := strings.Split(fen, " ")

	squares := convertFENPositionsTokenToBoard(tokens[0])

	var enpSquare *Square
	if tokens[3] != fenNoneSymbol {
		row, col := convertCoordinateToRowAndColumn(Coordinate(tokens[3]))
		enpSquare = NewSquare(row, col, nil)
	}

	var turn Color
	if tokens[1] == White.String() {
		turn = White
	}
	if tokens[1] == Black.String() {
		turn = Black
	}

	halfMoves, err := strconv.ParseUint(tokens[4], 10, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to parse half moves: %w", err)
	}

	fullMoves, err := strconv.ParseUint(tokens[5], 10, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to parse full moves: %w", err)
	}

	var castlingRights CastlingRights
	castle := tokens[2]
	if castle != fenNoneSymbol {
		if strings.Contains(castle, "K") {
			castlingRights |= WhiteKingSideCastle
		}
		if strings.Contains(castle, "Q") {
			castlingRights |= WhiteQueenSideCastle
		}
		if strings.Contains(castle, "k") {
			castlingRights |= BlackKingSideCastle
		}
		if strings.Contains(castle, "q") {
			castlingRights |= BlackQueenSideCastle
		}
	}

	meta := &Metadata{
		board:          squares,
		turn:           turn,
		castlingRights: castlingRights,
		enpSquare:      enpSquare,
		halfMoves:      uint8(halfMoves),
		fullMoves:      uint8(fullMoves),
	}

	return meta, nil
}

func getUnambiguousMoveNotation(move Move, moves []Move) string {
	var ambiguous bool
	var sameRank bool
	var sameFile bool

	for _, m := range moves {
		if move.pieceMoved.SymbolEquals(m.PieceMoved()) && !move.fromSquare.CoordEquals(m.fromSquare) && !move.toSquare.CoordEquals(m.toSquare) {
			ambiguous = true

			if move.fromSquare.Rank() == m.fromSquare.Rank() {
				sameRank = true
			}

			if move.toSquare.File() == m.toSquare.File() {
				sameFile = true
			}
		}
	}

	if ambiguous {
		coord := move.fromSquare.Coordinate()

		if sameRank && sameFile {
			return coord.String()
		} else if sameFile {
			return coord.Rank().String()
		} else {
			return coord.File().String()
		}
	}

	return ""
}

func make2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)

	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}

	return matrix
}

func initBoardSquares() [][]*Square {
	squares := make2D[*Square](boardSize, boardSize)

	for r := 0; r < len(squares); r++ {
		for c := 0; c < len(squares[r]); c++ {
			sq := NewSquare(Row(r), Column(c), nil)
			squares[r][c] = sq
		}
	}

	return squares
}

func drawBoardASCII(squares [][]*Square) string {
	if len(squares) == 0 {
		return ""
	}

	s := "   +------------------------+\n"

	for r := 0; r < len(squares); r++ {
		for c := 0; c < len(squares[r]); c++ {
			if c%8 == 0 {
				s += " " + convertRowToRank(Row(r)).String() + " |"
			}

			if squares[r][c].HasPiece() {
				s += " " + squares[r][c].piece.ToFENSymbol().String() + " "
			} else {
				s += " - "
			}

			if (c+1)%8 == 0 {
				s += "| \n"
			}
		}
	}

	s += "   +------------------------+\n"
	s += "     a  b  c  d  e  f  g  h"

	return s
}

func calculateSquareColor(row, col uint8) Color {
	if (row+col)%2 == 0 {
		return White
	}
	return Black
}
