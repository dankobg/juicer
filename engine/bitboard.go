package engine

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
// it uses little endian rank-file (LERF) mapping which has A1 as LSB (bit 0) and H8 as MSB (bit 63)
// it goes like this: (a1, b1 ... h1, a2, b2 ... h2, a8, b8 ... h8)
type bitboard uint64

const (
	bitboardFull         bitboard = 0xFFFFFFFFFFFFFFFF
	bitboardEmpty        bitboard = 0x0
	bitboardLightSquares bitboard = 0x55AA55AA55AA55AA
	bitboardDarkSquares  bitboard = 0xAA55AA55AA55AA55
	bitboardCorners      bitboard = 0x8100000000000081
	bitboardCenter       bitboard = 0x0000001818000000
)

var (
	// bitboardEmptyFilesMask represents clear mask where only specified file is filled with 0 and rest are 1
	bitboardEmptyFilesMask = map[File]bitboard{
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
	bitboardUniverseFilesMask = map[File]bitboard{
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
	bitboardEmptyRanksMask = map[Rank]bitboard{
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
	bitboardUniverseRanksMask = map[Rank]bitboard{
		Rank1: 0xFF,
		Rank2: 0xFF00,
		Rank3: 0xFF0000,
		Rank4: 0xFF000000,
		Rank5: 0xFF00000000,
		Rank6: 0xFF0000000000,
		Rank7: 0xFF000000000000,
		Rank8: 0xFF00000000000000,
	}
)

// clear sets the bitboard to empty
func (bb *bitboard) clear() {
	*bb = bitboardEmpty
}

// fill sets the bitboard to filled
func (bb *bitboard) fill() {
	*bb = bitboardFull
}

// setBit sets the bit to 1 at specified square
func (bb *bitboard) setBit(sq Square) {
	if sq.IndexInBoard() {
		*bb |= 1 << sq
	}
}

// toggleBit toggles the bit, it sets 0->1 and 1->0
func (bb *bitboard) toggleBit(sq Square) {
	if sq.IndexInBoard() {
		*bb ^= 1 << sq
	}
}

// clearBit sets the bit to 0 at specified square
func (bb *bitboard) clearBit(sq Square) {
	if sq.IndexInBoard() {
		*bb &= ^(1 << sq)
	}
}

// getBit gets the bit at specified square
func (bb bitboard) getBit(sq Square) uint8 {
	if !sq.IndexInBoard() {
		return 0
	}

	return uint8((bb >> sq) & 1)
}

// bitIsSet checks if bit value is 1 at specified square
func (bb bitboard) bitIsSet(sq Square) bool {
	if !sq.IndexInBoard() {
		return false
	}

	return ((bb >> sq) & 1) == 1
}

// populationCount returns the count of 1 bits in bitboard
func (bb bitboard) populationCount() uint8 {
	return uint8(bits.OnesCount64(uint64(bb)))
}

// MS1B gets the index of most significant 1 bit
func (bb bitboard) MS1B() int {
	return 63 - bits.LeadingZeros64(uint64(bb))
}

// LS1B gets the index of least significant 1 bit
func (bb bitboard) LS1B() int {
	return bits.TrailingZeros64(uint64(bb))
}

// PopLS1B pops (clears) the LS1B and returns its index
func (bb *bitboard) PopLS1B() int {
	ls1b := bb.LS1B()
	bb.clearBit(Square(ls1b))

	return ls1b
}

// PopMS1B pops (clears) the MS1B and returns its index
func (bb *bitboard) PopMS1B() int {
	ms1b := bb.MS1B()
	bb.clearBit(Square(ms1b))

	return ms1b
}

// isEmpty checks whether bitboard is empty (all 0s)
func (bb bitboard) isEmpty() bool {
	return bb == bitboardEmpty
}

// isFull checks whether bitboard is full (all 1s)
func (bb bitboard) isFull() bool {
	return bb == bitboardFull
}

// flipVertical flips the bitboard vertically
func (bb bitboard) flipVertical() bitboard {
	return bitboard(bits.ReverseBytes64(uint64(bb)))
}

// flipHorizontal flips the bitboard horizontally
func (bb bitboard) flipHorizontal() bitboard {
	k1 := bitboard(0x5555555555555555)
	k2 := bitboard(0x3333333333333333)
	k4 := bitboard(0x0F0F0F0F0F0F0F0F)

	bb = ((bb >> 1) & k1) + 2*(bb&k1)
	bb = ((bb >> 2) & k2) + 4*(bb&k2)
	bb = ((bb >> 4) & k4) + 16*(bb&k4)

	return bb
}

// flipDiagonalA1H8 flips the bitboard diagonally from a1 to h1
func (bb bitboard) flipDiagonalA1H8() bitboard {
	var t bitboard
	k1 := bitboard(0x5500550055005500)
	k2 := bitboard(0x3333000033330000)
	k4 := bitboard(0x0F0F0F0F00000000)

	t = k4 & (bb ^ (bb << 28))
	bb ^= t ^ (t >> 28)
	t = k2 & (bb ^ (bb << 14))
	bb ^= t ^ (t >> 14)
	t = k1 & (bb ^ (bb << 7))
	bb ^= t ^ (t >> 7)

	return bb
}

// flipDiagonalA8H1 flips the bitboard diagonally from a8 to h8
func (bb bitboard) flipDiagonalA8H1() bitboard {
	var t bitboard
	k1 := bitboard(0xAA00AA00AA00AA00)
	k2 := bitboard(0xCCCC0000CCCC0000)
	k4 := bitboard(0xF0F0F0F00F0F0F0F)

	t = bb ^ (bb << 36)
	bb ^= k4 & (t ^ (bb >> 36))
	t = k2 & (bb ^ (bb << 18))
	bb ^= t ^ (t >> 18)
	t = k1 & (bb ^ (bb << 9))
	bb ^= t ^ (t >> 9)

	return bb
}

// rotate180 rotates the bitboard 180 degrees
func (bb bitboard) rotate180() bitboard {
	bb = bb.flipVertical()
	bb = bb.flipHorizontal()

	return bb
}

// rotate90clockwise rotates the bitboard 90 degrees
func (bb bitboard) rotate90clockwise() bitboard {
	bb = bb.flipVertical()
	bb = bb.flipDiagonalA8H1()

	return bb
}

// rotate90counterClockwise rotates the bitboard 270 degrees (90 counter-clockwise)
func (bb bitboard) rotate90counterClockwise() bitboard {
	bb = bb.flipVertical()
	bb = bb.flipDiagonalA1H8()

	return bb
}

// draw prints the board in 8x8 grid in ascii style
// it prints 1/0 whether the piece bit is set at specified square
func (bb bitboard) draw(options *DrawOptions) string {
	return printBoard(options, func(sq Square) string {
		return fmt.Sprint(bb.getBit(sq))
	})
}

// String implements the stringer and returns the uint64 as string
func (bb bitboard) String() string {
	return fmt.Sprint(uint64(bb))
}

type DrawOptions struct {
	Compact bool
	Side    Color
}

// printBoard prints the board in 8x8 grid with ascii style
func printBoard(options *DrawOptions, printerFunc func(sq Square) string) string {
	opts := DrawOptions{Side: White}
	if options != nil {
		opts = *options
	}

	var sb strings.Builder

	if !opts.Compact {
		sb.WriteString("   +------------------------+\n")
	}

	for r := boardSize - 1; r >= 0; r-- {
		rankStartChar := ""
		if !opts.Compact {
			rankStartChar = "|"
		}

		rankIdx, rankLabel := r, r+1
		if opts.Side == Black {
			rankIdx, rankLabel = 7-r, 8-r
		}

		sb.WriteString(fmt.Sprintf(" %d %s", rankLabel, rankStartChar))

		for f := range 8 {
			fileSpacingChar := ""
			if !opts.Compact {
				fileSpacingChar = " "
			}

			fileIdx := f
			if opts.Side == Black {
				fileIdx = 7 - f
			}

			sq := Square(rankIdx*8 + fileIdx)
			sb.WriteString(fmt.Sprintf(" %s%s", printerFunc(sq), fileSpacingChar))
		}

		rankEndChar := ""
		if !opts.Compact {
			rankEndChar = "|"
		}

		sb.WriteString(fmt.Sprintf("%s \n", rankEndChar))
	}

	if opts.Compact {
		fileLabels := "\n    a b c d e f g h"
		if opts.Side == Black {
			fileLabels = "\n    h g f e d c b a"
		}

		sb.WriteString(fileLabels)
	} else {
		sb.WriteString("   +------------------------+\n")

		fileLabels := "     a  b  c  d  e  f  g  h"
		if opts.Side == Black {
			fileLabels = "     h  g  f  e  d  c  b  a"
		}

		sb.WriteString(fileLabels)
	}

	return sb.String()
}
