package engine

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type Row uint8
type Column uint8

type File string
type Rank uint8

type Coordinate string

func idxInRange[T constraints.Integer](min, max T, nums ...T) bool {
	ok := true

	for _, idx := range nums {
		if !(idx >= min && idx <= max) {
			ok = false
			break
		}
	}

	return ok
}

func validRowCol[T constraints.Integer](nums ...T) bool {
	return idxInRange(0, 7, nums...)
}

func validRank[T constraints.Integer](nums ...T) bool {
	return idxInRange(1, 8, nums...)
}

func validFile(file string) bool {
	return file == "a" || file == "b" || file == "c" || file == "d" || file == "e" || file == "f" || file == "g" || file == "h"
}

func (r Row) String() string {
	return strconv.FormatUint(uint64(r), 10)
}

func (r Row) Valid() bool {
	return validRowCol(uint8(r))
}

func (c Column) String() string {
	return strconv.FormatUint(uint64(c), 10)
}

func (c Column) Valid() bool {
	return validRowCol(uint8(c))
}

func (f File) String() string {
	return string(f)
}

func (f File) Valid() bool {
	return validFile(string(f))
}

func (r Rank) String() string {
	return strconv.FormatUint(uint64(r), 10)
}

func (r Rank) Valid() bool {
	return validRank(uint8(r))
}

func (c Coordinate) String() string {
	return string(c)
}

func (c Coordinate) Valid() bool {
	coord := string(c)

	if len(coord) != 2 {
		return false
	}

	file := File(coord[0])
	rank := Rank(coord[1])

	return file.Valid() && rank.Valid()
}

func (c Coordinate) File() File {
	file, _ := convertCoordinateToFileAndRank(c)
	return file
}

func (c Coordinate) Rank() Rank {
	_, rank := convertCoordinateToFileAndRank(c)
	return rank
}

func (c Coordinate) RowCol() (Row, Column) {
	row, col := convertCoordinateToRowAndColumn(c)
	return row, col
}

func (c Coordinate) Row() Row {
	row, _ := convertCoordinateToRowAndColumn(c)
	return row
}

func (c Coordinate) Column() Column {
	_, col := convertCoordinateToRowAndColumn(c)
	return col
}

type Square struct {
	row    Row
	column Column
	color  Color
	piece  *Piece
}

func NewSquare(row Row, column Column, piece *Piece) *Square {
	color := calculateSquareColor(uint8(row), uint8(column))

	return &Square{
		row:    row,
		column: column,
		piece:  piece,
		color:  color,
	}
}

func (s *Square) String() string {
	return s.Coordinate().String()
}

func (s *Square) Row() Row {
	return s.row
}

func (s *Square) Column() Column {
	return s.column
}

func (s *Square) Color() Color {
	return s.color
}

func (s *Square) Piece() *Piece {
	return s.piece
}

func (s *Square) File() File {
	return convertColumnToFile(s.column)
}

func (s *Square) Rank() Rank {
	return convertRowToRank(s.row)
}

func (s *Square) Coordinate() Coordinate {
	return convertRowAndColumnToCoordinate(s.row, s.column)
}

func (s *Square) RowCol() (Row, Column) {
	return s.Coordinate().RowCol()
}

func (s *Square) IsDark() bool {
	return s.color.IsBlack()
}

func (s *Square) IsLight() bool {
	return s.color.IsWhite()
}

func (s *Square) HasPiece() bool {
	return s.piece != nil
}

func (s *Square) IsEmpty() bool {
	return !s.HasPiece()
}

func (s *Square) HasFriendlyPiece(currentTurn Color) bool {
	return s.HasPiece() && s.piece.IsFriendly(currentTurn)
}

func (s *Square) HasEnemyPiece(currentTurn Color) bool {
	return s.HasPiece() && s.piece.IsEnemy(currentTurn)
}

func (s *Square) isEmptyOrHasEnemyPiece(currentTurn Color) bool {
	return s.IsEmpty() || s.HasEnemyPiece(currentTurn)
}

func (s *Square) ToFENSymbol() string {
	if s.piece != nil {
		return s.piece.ToFENSymbol().String()
	}
	return ""
}

func (s *Square) CoordEquals(other Square) bool {
	return s.Coordinate().String() == other.Coordinate().String()
}
