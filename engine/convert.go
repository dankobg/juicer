package engine

import (
	"strconv"
	"strings"
)

var fileToColumn = map[File]Column{"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7}
var columnToFile = map[Column]File{0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h"}
var rankToRow = map[Rank]Row{1: 7, 2: 6, 3: 5, 4: 4, 5: 3, 6: 2, 7: 1, 8: 0}
var rowToRank = map[Row]Rank{0: 8, 1: 7, 2: 6, 3: 5, 4: 4, 5: 3, 6: 2, 7: 1}

func swapColor(color Color) Color {
	if color == Black {
		return White
	}
	return Black
}

func convertFileToColumn(file File) Column {
	c := fileToColumn[file]
	return c
}

func convertColumnToFile(col Column) File {
	f := columnToFile[col]
	return f
}

func convertRankToRow(rank Rank) Row {
	r := rankToRow[rank]
	return r
}

func convertRowToRank(row Row) Rank {
	r := rowToRank[row]
	return r
}

func convertCoordinateToFileAndRank(coord Coordinate) (File, Rank) {
	chars := strings.Split(coord.String(), "")

	file := File(chars[0])

	rankNum, _ := strconv.ParseUint(chars[1], 10, 8)
	rank := Rank(uint8(rankNum))

	return file, rank
}

func convertCoordinateToRowAndColumn(coord Coordinate) (Row, Column) {
	chars := strings.Split(coord.String(), "")

	file := File(chars[0])

	rankNum, _ := strconv.ParseUint(chars[1], 10, 8)
	rank := Rank(uint8(rankNum))

	row := convertRankToRow(rank)
	col := convertFileToColumn(file)

	return row, col
}

func convertFileAndRankToCoordinate(file File, rank Rank) Coordinate {
	return Coordinate(file.String() + rank.String())
}

func convertRowAndColumnToCoordinate(row Row, col Column) Coordinate {
	file := convertColumnToFile(col)
	rank := convertRowToRank(row)

	return convertFileAndRankToCoordinate(file, rank)
}
