package juicer

import "sync"

var (
	once                   sync.Once
	initializedAttackMasks bool

	whitePawnAttacksMask map[Square]bitboard
	blackPawnAttacksMask map[Square]bitboard
	kingAttacksMask      map[Square]bitboard
	knightsAttacksMask   map[Square]bitboard

	bishopRelevantOccupancyBitsMask map[Square]bitboard
	rookRelevantOccupancyBitsMask   map[Square]bitboard

	bishopRelevantOccupancyBitsPopulationCount map[Square]uint8
	rookRelevantOccupancyBitsPopulationCount   map[Square]uint8
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
