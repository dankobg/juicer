package juicer

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	boardSize         = 8
	boardTotalSquares = 64

	fileChars = "abcdefgh"
	rankChars = "12345678"
)

var (
	reValidFile  = regexp.MustCompile("^[a-h]$")
	reValidRank  = regexp.MustCompile("^[1-8]$")
	reValidCoord = regexp.MustCompile("^[a-h][1-8]$")

	coords = []string{
		"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
		"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
		"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
		"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
		"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
		"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
		"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
		"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	}

	coordToSquare = map[string]Square{
		"a1": A1, "a2": A2, "a3": A3, "a4": A4, "a5": A5, "a6": A6, "a7": A7, "a8": A8,
		"b1": B1, "b2": B2, "b3": B3, "b4": B4, "b5": B5, "b6": B6, "b7": B7, "b8": B8,
		"c1": C1, "c2": C2, "c3": C3, "c4": C4, "c5": C5, "c6": C6, "c7": C7, "c8": C8,
		"d1": D1, "d2": D2, "d3": D3, "d4": D4, "d5": D5, "d6": D6, "d7": D7, "d8": D8,
		"e1": E1, "e2": E2, "e3": E3, "e4": E4, "e5": E5, "e6": E6, "e7": E7, "e8": E8,
		"f1": F1, "f2": F2, "f3": F3, "f4": F4, "f5": F5, "f6": F6, "f7": F7, "f8": F8,
		"g1": G1, "g2": G2, "g3": G3, "g4": G4, "g5": G5, "g6": G6, "g7": G7, "g8": G8,
		"h1": H1, "h2": H2, "h3": H3, "h4": H4, "h5": H5, "h6": H6, "h7": H7, "h8": H8,
	}
)

// Square is one of the 64 squares on the board
type Square int8

func (sq Square) IndexInBoard() bool {
	return sq >= 0 && sq <= 63
}

// File returns the square's file
func (sq Square) File() File {
	return File(int(sq) % boardSize)
}

// Rank returns the square's rank
func (sq Square) Rank() Rank {
	return Rank(int(sq) / boardSize)
}

// Coordinate returns the square's coordinate (file, rank pair) e.g. `e4`
func (sq Square) Coordinate() string {
	return sq.File().String() + sq.Rank().String()
}

func (sq Square) String() string {
	return sq.Coordinate()
}

// NewSquare creates a new Square from a File and a Rank
func NewSquare(f File, r Rank) Square {
	return Square(int8(r)*boardSize + int8(f))
}

// NewSquareFromCoord creates a new Square from a coordinate
func NewSquareFromCoord(coord string) (Square, error) {
	if len(coord) != 2 {
		return SquareNone, fmt.Errorf("invalid coordinate length: %s", coord)
	}

	fileChar := string(coord[0])
	rankChar := string(coord[1])

	if !reValidCoord.MatchString(coord) {
		return SquareNone, fmt.Errorf("invalid coordinate: %s", coord)
	}

	f := strings.Index(fileChars, fileChar)
	r := strings.Index(rankChars, rankChar)

	return NewSquare(File(f), Rank(r)), nil
}

func (sq Square) Color() Color {
	if ((sq / 8) % 2) == (sq % 2) {
		return Black
	}
	return White
}

const (
	SquareNone Square = iota - 1
	A1
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// File is the square's file from a-h
type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

func NewFile(file string) (File, error) {
	if len(file) != 1 {
		return File(-1), fmt.Errorf("invaid file character length: %s", file)
	}

	if !reValidFile.MatchString(file) {
		return File(-1), fmt.Errorf("invalid file character: %s", file)
	}

	f := strings.Index(fileChars, file)
	return File(f), nil
}

func (f File) String() string {
	return fileChars[f : f+1]
}

// Rank is the square's rank from 1-8
type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

func NewRank(rank string) (Rank, error) {
	if len(rank) != 1 {
		return Rank(-1), fmt.Errorf("invaid rank character length: %s", rank)
	}

	if !reValidRank.MatchString(rank) {
		return Rank(-1), fmt.Errorf("invalid rank character: %s", rank)
	}

	r := strings.Index(rankChars, rank)
	return Rank(r), nil
}

func (r Rank) String() string {
	return rankChars[r : r+1]
}

func (sq Square) occupancyMask() bitboard {
	return 1 << sq
}
