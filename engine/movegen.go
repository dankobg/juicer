package juicer

import "sync"

var (
	initialized bool
	once        sync.Once

	whitePawnAttacksMask map[Square]bitboard
	blackPawnAttacksMask map[Square]bitboard
	kingAttacksMask      map[Square]bitboard
	knightsAttacksMask   map[Square]bitboard
)

func initAttacksMaskForNonSlidingPieces() {
	if !initialized {
		once.Do(func() {
			for i := 0; i < boardTotalSquares; i++ {
				sq := Square(i)

				initKingAttacksMask(sq)
				initKnightPawnAttacksMask(sq)
				initPawnAttacksMask(sq)
			}

			initialized = true
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

func generateKingAttacksMask(sq Square) bitboard {
	occupancy, attacks := bitboardEmpty, bitboardEmpty
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
	occupancy, attacks := bitboardEmpty, bitboardEmpty
	occupancy.setBit(sq)

	clearedAFile := occupancy & bitboardEmptyFilesMask[FileA]
	clearedHFile := occupancy & bitboardEmptyFilesMask[FileH]
	clearedABFile := clearedAFile & bitboardEmptyFilesMask[FileB]
	clearedGHFile := clearedHFile & bitboardEmptyFilesMask[FileG]

	attacks |= clearedHFile << 17  // Nort-NoEa
	attacks |= clearedGHFile << 10 // East-NoEa
	attacks |= clearedABFile << 6  // West-NoWe
	attacks |= clearedAFile << 15  // Nort-NoWe

	attacks |= clearedAFile >> 17  // Sout-SoWe
	attacks |= clearedABFile >> 10 // West-SoWe
	attacks |= clearedGHFile >> 6  // East-SoEa
	attacks |= clearedHFile >> 15  // Sout-SoEa

	return attacks
}

func generatePawnAttacksMask(sq Square, color Color) bitboard {
	occupancy, attacks := bitboardEmpty, bitboardEmpty
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
