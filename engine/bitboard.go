package juicer

import (
	"fmt"
	"math/bits"
	"strings"
)

// bitboard representation
//   +-------------------------+
// 8 | 56 57 58 59 60 61 62 63 |
// 7 | 48 49 50 51 52 53 54 55 |
// 6 | 40 41 42 43 44 45 46 47 |
// 5 | 32 33 34 35 36 37 38 39 |
// 4 | 24 25 26 27 28 29 30 31 |
// 3 | 16 17 18 19 20 21 22 23 |
// 2 | 8  9  10 11 12 13 14 15 |
// 1 | 0  1  2  3  4  5  6  7  |
// 	 +-------------------------+
// 		 a  b  c  d  e  f  g  h

// bitboard is the board representation using 64-bit unsigned integer.
// the mapping goes from a1 as LSB (bit 0) to h8 as MSB (bit 63)
type bitboard uint64

// bitboardFull represents bitboard where every bit is set to 1
const bitboardFull bitboard = 0xffffffffffffffff

// bitboardEmpty represents bitboard where every bit is set to 0
const bitboardEmpty bitboard = 0x0

// bitboardEmptyFilesMask represents clear mask where only specified file is filled with 0 and rest are 1
var bitboardEmptyFilesMask = map[File]bitboard{
	FileA: 0xFEFEFEFEFEFEFEFE,
	FileB: 0xFDFDFDFDFDFDFDFD,
	FileC: 0xFBFBFBFBFBFBFBFB,
	FileD: 0xF7F7F7F7F7F7F7F7,
	FileE: 0xEFEFEFEFEFEFEFEF,
	FileF: 0xDFDFDFDFDFDFDFDF,
	FileG: 0xBFBFBFBFBFBFBFBF,
	FileH: 0x7F7F7F7F7F7F7F7F,
}

// bitboardUniverseFilesMask represents fill mask where only specified file is filled with 1 and rest are 0
var bitboardUniverseFilesMask = map[File]bitboard{
	FileA: 0x101010101010101,
	FileB: 0x202020202020202,
	FileC: 0x404040404040404,
	FileD: 0x808080808080808,
	FileE: 0x1010101010101010,
	FileF: 0x2020202020202020,
	FileG: 0x4040404040404040,
	FileH: 0x8080808080808080,
}

// bitboardEmptyRanksMask represents clear mask where only specified rank is filled with 0 and rest are 1
var bitboardEmptyRanksMask = map[Rank]bitboard{
	Rank1: 0xFFFFFFFFFFFFFF00,
	Rank2: 0xFFFFFFFFFFFF00FF,
	Rank3: 0xFFFFFFFFFF00FFFF,
	Rank4: 0xFFFFFFFF00FFFFFF,
	Rank5: 0xFFFFFF00FFFFFFFF,
	Rank6: 0xFFFF00FFFFFFFFFF,
	Rank7: 0xFF00FFFFFFFFFFFF,
	Rank8: 0xFFFFFFFFFFFFFF,
}

// bitboardUniverseRanksMask represents fill mask where only specified rank is filled with 1 and rest are 0
var bitboardUniverseRanksMask = map[Rank]bitboard{
	Rank1: 0xFF,
	Rank2: 0xFF00,
	Rank3: 0xFF0000,
	Rank4: 0xFF000000,
	Rank5: 0xFF00000000,
	Rank6: 0xFF0000000000,
	Rank7: 0xFF000000000000,
	Rank8: 0xFF00000000000000,
}

// getBit gets the bit at specified square
func (bb *bitboard) getBit(sq Square) uint8 {
	if !sq.IndexInBoard() {
		return 0
	}

	return uint8((*bb >> sq) & 1)
}

// setBit sets the bit to 1 at specified square
func (bb *bitboard) setBit(squares ...Square) {
	for _, sq := range squares {
		if sq.IndexInBoard() {
			*bb |= 1 << sq
		}
	}
}

// clearBit sets the bit to 0 at specified square
func (bb *bitboard) clearBit(squares ...Square) {
	for _, sq := range squares {
		if sq.IndexInBoard() {
			*bb &= ^(1 << sq)
		}
	}
}

// bitIseSet checks whether the bit is set to 0 at specified square
func (bb *bitboard) bitIsSet(sq Square) bool {
	if !sq.IndexInBoard() {
		return false
	}

	return ((*bb >> sq) & 1) == 1
}

// countBits returns the count of 1 bits in bitboard
func (bb *bitboard) countBits() uint8 {
	return uint8(bits.OnesCount64(uint64(*bb)))
}

func (bb *bitboard) setEmpty() {
	*bb = bitboardEmpty
}

func (bb *bitboard) setFull() {
	*bb = bitboardFull
}

// drawCompact prints the bitboard for debugging as 8x8 grid of 0s and 1s in a compact way
func (bb bitboard) drawCompact() string {
	var sb strings.Builder

	for r := boardSize - 1; r >= 0; r-- {
		sb.WriteString(fmt.Sprintf(" %d ", r+1))

		for f := 0; f < boardSize; f++ {
			idx := r*8 + f
			sb.WriteString(fmt.Sprintf(" %d", bb.getBit(Square(idx))))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("\n    a b c d e f g h")

	return sb.String()
}

// drawPretty prints the bitboard for debugging as 8x8 grid of 0s and 1s in a pretty way
func (bb bitboard) drawPretty() string {
	var sb strings.Builder
	sb.WriteString("   +------------------------+\n")

	for r := boardSize - 1; r >= 0; r-- {
		sb.WriteString(fmt.Sprintf(" %d |", r+1))

		for f := 0; f < 8; f++ {
			idx := r*8 + f
			sb.WriteString(fmt.Sprintf(" %d ", bb.getBit(Square(idx))))
		}

		sb.WriteString("| \n")
	}

	sb.WriteString("   +------------------------+\n")
	sb.WriteString("     a  b  c  d  e  f  g  h")

	return sb.String()
}
