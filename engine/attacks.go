package juicer

import "sync"

var (
	once                   sync.Once
	initializedAttackMasks bool

	whitePawnAttacksMask = make(map[Square]bitboard, 64)
	blackPawnAttacksMask = make(map[Square]bitboard, 64)
	kingAttacksMask      = make(map[Square]bitboard, 64)
	knightsAttacksMask   = make(map[Square]bitboard, 64)

	bishopRelevantOccupancyBitsMask = make(map[Square]bitboard, 64)
	rookRelevantOccupancyBitsMask   = make(map[Square]bitboard, 64)

	bishopRelevantOccupancyBitsPopulationCount = make(map[Square]uint8, 64)
	rookRelevantOccupancyBitsPopulationCount   = make(map[Square]uint8, 64)

	bishopMagics = make(map[Square]bitboard, 64)
	rookMagics   = make(map[Square]bitboard, 64)
)

func initAttacksMaskForNonSlidingPieces() {
	if !initializedAttackMasks {
		once.Do(func() {
			for i := 0; i < boardTotalSquares; i++ {
				sq := Square(i)

				initKingAttacksMask(sq)
				initKnightPawnAttacksMask(sq)
				initPawnAttacksMask(sq)

				initBishopRelevantOccupancyBitsMask(sq)
				initRookRelevantOccupancyBitsMask(sq)

				initBishopRelevantOccupancyBitsPopulationCount(sq)
				initRookRelevantOccupancyBitsPopulationCount(sq)
			}

			initializedAttackMasks = true
		})
	}
}

func initKingAttacksMask(sq Square) {
	kingAttacksMask[sq] = generateKingAttacksMask(sq)
}

func initKnightPawnAttacksMask(sq Square) {
	knightsAttacksMask[sq] = generateKnightAttacksMask(sq)
}

func initPawnAttacksMask(sq Square) {
	whitePawnAttacksMask[sq] = generatePawnAttacksMask(sq, White)
	blackPawnAttacksMask[sq] = generatePawnAttacksMask(sq, Black)
}

func initBishopRelevantOccupancyBitsMask(sq Square) {
	bishopRelevantOccupancyBitsMask[sq] = generateBishopRelevantOccupancyBitsMask(sq)
}

func initRookRelevantOccupancyBitsMask(sq Square) {
	rookRelevantOccupancyBitsMask[sq] = generateRookRelevantOccupancyBitsMask(sq)
}

func initBishopRelevantOccupancyBitsPopulationCount(sq Square) {
	bishopRelevantOccupancyBitsPopulationCount[sq] = bishopRelevantOccupancyBitsMask[sq].populationCount()
}

func initRookRelevantOccupancyBitsPopulationCount(sq Square) {
	rookRelevantOccupancyBitsPopulationCount[sq] = rookRelevantOccupancyBitsMask[sq].populationCount()
}

func generateKingAttacksMask(sq Square) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	clearedAFile := occupancy & bitboardEmptyFilesMask[FileA]
	clearedHFile := occupancy & bitboardEmptyFilesMask[FileH]

	attacks |= clearedAFile << 7 // NoWe
	attacks |= occupancy << 8    // Nort
	attacks |= clearedHFile << 9 // NoEa
	attacks |= clearedHFile << 1 // East

	attacks |= clearedHFile >> 7 // SoEa
	attacks |= occupancy >> 8    // Sout
	attacks |= clearedAFile >> 9 // SoWe
	attacks |= clearedAFile >> 1 // West

	return attacks
}

func generateKnightAttacksMask(sq Square) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	clearedAFile := occupancy & bitboardEmptyFilesMask[FileA]
	clearedHFile := occupancy & bitboardEmptyFilesMask[FileH]
	clearedABFile := clearedAFile & bitboardEmptyFilesMask[FileB]
	clearedGHFile := clearedHFile & bitboardEmptyFilesMask[FileG]

	attacks |= clearedABFile << 6  // West-NoWe
	attacks |= clearedAFile << 15  // Nort-NoWe
	attacks |= clearedHFile << 17  // Nort-NoEa
	attacks |= clearedGHFile << 10 // East-NoEa

	attacks |= clearedGHFile >> 6  // East-SoEa
	attacks |= clearedHFile >> 15  // Sout-SoEa
	attacks |= clearedAFile >> 17  // Sout-SoWe
	attacks |= clearedABFile >> 10 // West-SoWe

	return attacks
}

func generatePawnAttacksMask(sq Square, color Color) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	clearedAFile := occupancy & bitboardEmptyFilesMask[FileA]
	clearedHFile := occupancy & bitboardEmptyFilesMask[FileH]

	if color.IsWhite() {
		attacks |= clearedAFile << 7 // NoWe
		attacks |= clearedHFile << 9 // NoEa
	}

	if color.IsBlack() {
		attacks |= clearedAFile >> 9 // SoWe
		attacks |= clearedHFile >> 7 // SoEa
	}

	return attacks
}

// generateBishopRelevantOccupancyBitsMask generates the bishop relevant occupancy look table
func generateBishopRelevantOccupancyBitsMask(sq Square) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	f, r := int(sq%8), int(sq/8)

	for i := 1; f-i > 0 && r+i < 7; i++ {
		attacks.setBit(Square((r+i)*8 + f - i)) // NoWe
	}

	for i := 1; f+i < 7 && r+i < 7; i++ {
		attacks.setBit(Square((r+i)*8 + f + i)) // NoEa
	}

	for i := 1; f-i > 0 && r-i > 0; i++ {
		attacks.setBit(Square((r-i)*8 + f - i)) // SoWe
	}

	for i := 1; f+i < 7 && r-i > 0; i++ {
		attacks.setBit(Square((r-i)*8 + f + i)) // SoEa
	}

	return attacks
}

// generateRookRelevantOccupancyBitsMask generates the rook relevant occupancy look table
func generateRookRelevantOccupancyBitsMask(sq Square) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	f, r := int(sq%8), int(sq/8)

	for i := 1; r+i < 7; i++ {
		attacks.setBit(Square((r+i)*8 + f)) // North
	}

	for i := 1; f+i < 7; i++ {
		attacks.setBit(Square(r*8 + f + i)) // East
	}

	for i := 1; r-i > 0; i++ {
		attacks.setBit(Square((r-i)*8 + f)) // Sout
	}

	for i := 1; f-i > 0; i++ {
		attacks.setBit(Square(r*8 + f - i)) // West
	}

	return attacks
}

// generateBishopAttacksWithBlockers generates bishop sliding attacks with a blocker bitboard on the fly
// it is only used for finding and initializing magic numbers becayse it is too slow to use in in movegen
func generateBishopAttacksWithBlockers(sq Square, blockers bitboard) bitboard {
	var occupancy, attacks bitboard
	occupancy.setBit(sq)

	f, r := int(sq%8), int(sq/8)
	target := SquareNone

	for i := 1; f-i >= 0 && r+i < 8; i++ {
		target = Square((r+i)*8 + f - i) // NoWe
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; f+i < 8 && r+i < 8; i++ {
		target = Square((r+i)*8 + f + i) // NoEa
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; f-i >= 0 && r-i >= 0; i++ {
		target = Square((r-i)*8 + f - i) // SoWe
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; f+i < 8 && r-i >= 0; i++ {
		target = Square((r-i)*8 + f + i) // SoEa
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	return attacks
}

// generateRookAttacksWithBlockers generates rook sliding attacks with a blocker bitboard on the fly
// it is only used for finding and initializing magic numbers becayse it is too slow to use in in movegen
func generateRookAttacksWithBlockers(sq Square, blockers bitboard) bitboard {
	var piece, attacks bitboard
	piece.setBit(sq)

	f, r := int(sq%8), int(sq/8)
	target := SquareNone

	for i := 1; f+i < 8; i++ {
		target = Square(r*8 + f + i) // East
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; r-i >= 0; i++ {
		target = Square((r-i)*8 + f) // Sout
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; f-i >= 0; i++ {
		target = Square(r*8 + f - i) // West
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	for i := 1; r+i < 8; i++ {
		target = Square((r+i)*8 + f) // Nort
		attacks.setBit(target)

		if blockers&target.occupancyMask() != 0 {
			break
		}
	}

	return attacks
}

// Occupancy generates occupancy bitboards for a given relevant occupancy bitboard
func Occupancy(index, count int, attack bitboard) bitboard {
	var occupancy bitboard

	for i := 0; i < count; i++ {
		sq := Square(attack.LS1B())
		attack.clearBit(sq)

		if index&(1<<i) != 0 {
			occupancy.setBit(sq)
		}
	}

	return occupancy
}
