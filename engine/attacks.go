package juicer

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"sync"
)

var (
	once                    sync.Once
	initializedAttackTables bool

	pawnAttacksMask    [2][64]bitboard
	kingAttacksMask    [64]bitboard
	knightsAttacksMask [64]bitboard

	bishopRelevantOccupancyBitsMask [64]bitboard
	rookRelevantOccupancyBitsMask   [64]bitboard

	BishopRelevantOccupancyBitsPopulationCount [64]uint8
	RookRelevantOccupancyBitsPopulationCount   [64]uint8

	bishopAttacksMask [64][512]bitboard
	rookAttacksMask   [64][4096]bitboard

	bishopMagics = [64]bitboard{
		0x20010400808600, 0xa008010410820000, 0x1004440082038008, 0x904040098084800,
		0x600c052000520541, 0x4002010420402022, 0x11040104400480, 0x200104104202080,
		0x1200210204080080, 0x6c18600204e20682, 0x2202004200e0, 0x100044404810840,
		0x400220211108110, 0x20002011009000c, 0xa00200a2084210, 0x202008098011000,
		0xc40002004019206, 0x116042040804c500, 0x419002080a80200a, 0x4000844000800,
		0x404b080a04800, 0x4608080482012002, 0x44040500a0880841, 0x2002100909050d00,
		0x8404004030a400, 0x90709004040080, 0x11444043040d0204, 0x8080100202020,
		0x801001181004000, 0x4140822002021000, 0x102089092009006, 0x540a042100540203,
		0x50100409482820, 0x8010880900041004, 0x230100500414, 0x200800050810,
		0x8294064010040100, 0x9010100220044404, 0x154202022004008e, 0x9420220008401,
		0x71080840110401, 0x2000a40420400201, 0x802619048001004, 0x209280a058000500,
		0x2004044810100a00, 0xa0208d000804300, 0x638a80d000684, 0x1910401000080,
		0x800420210400200, 0x4404410090100, 0x8020808400880000, 0x400081042120c21,
		0x4009001022120001, 0x4902220802082000, 0x410841000820290, 0x820020401002440,
		0x800420041084000, 0x10818c05a000, 0x301804213d000, 0x800040018208801,
		0x1b80000004104405, 0x2500214084184884, 0x1000628801050400, 0x8040229e24002080,
	}

	rookMagics = [64]bitboard{
		0x18010a040018000, 0x40002000401001, 0x290010a841e00100, 0x29001000050900a0,
		0x4080030400800800, 0x1200040200100801, 0x2200208200040851, 0x220000820425004c,
		0x104800740008020, 0x420400020005000, 0x844801000200480, 0x4004808008001000,
		0x4009000410080100, 0x3000400020900, 0x4804000810020104, 0x74800641800900,
		0x862818014400020, 0x40048020004480, 0x11a1010040200012, 0x20828010000800,
		0x848808004020800, 0x4522808004000200, 0x10100020004, 0x400206000092411c,
		0x818004444000a000, 0x180a000c0005002, 0xb104100200100, 0x24022202000a4010,
		0x100040080080080, 0x2010200080490, 0x180390400221098, 0x410008200010044,
		0x310400089800020, 0x8c0804009002902, 0x1004402001001504, 0x105021001000920,
		0x40080800801, 0xa02001002000804, 0x108284204005041, 0x8004082002411,
		0x2802281c0028001, 0x9044000910020, 0x200010008080, 0x40201001010008,
		0x8000080004008080, 0x3010400420080110, 0x414210040008, 0x10348400460001,
		0x80002000401040, 0x460200088400080, 0x8201822000100280, 0x600100008008280,
		0xc0800800040080, 0x24040080020080, 0x22c11a0108100c00, 0x204008114104200,
		0x8800800010290041, 0x401500228206, 0x8002a00011090041, 0x42008100101,
		0x283000800100205, 0x2008810010402, 0x490102200880104, 0x800010920940042,
	}

	F1G1 = bitboardUniverseFilesMask[FileF]&bitboardUniverseRanksMask[Rank1] | bitboardUniverseFilesMask[FileG]&bitboardUniverseRanksMask[Rank1]
	C1D1 = bitboardUniverseFilesMask[FileC]&bitboardUniverseRanksMask[Rank1] | bitboardUniverseFilesMask[FileD]&bitboardUniverseRanksMask[Rank1]
	B1D1 = bitboardUniverseFilesMask[FileB]&bitboardUniverseRanksMask[Rank1] | bitboardUniverseFilesMask[FileD]&bitboardUniverseRanksMask[Rank1]

	F8G8 = bitboardUniverseFilesMask[FileF]&bitboardUniverseRanksMask[Rank8] | bitboardUniverseFilesMask[FileG]&bitboardUniverseRanksMask[Rank8]
	C8D8 = bitboardUniverseFilesMask[FileC]&bitboardUniverseRanksMask[Rank8] | bitboardUniverseFilesMask[FileD]&bitboardUniverseRanksMask[Rank8]
	B8D8 = bitboardUniverseFilesMask[FileB]&bitboardUniverseRanksMask[Rank8] | bitboardUniverseFilesMask[FileD]&bitboardUniverseRanksMask[Rank8]
)

func initAllAttackMasksTables() {
	if !initializedAttackTables {
		once.Do(func() {
			initAttackMasksForNonSlidingPieces()
			initBishopAndRookPopCounts()
			initAttackMasksForSlidingPieces()

			initializedAttackTables = true
		})
	}
}

func initAttackMasksForNonSlidingPieces() {
	for sq := A1; sq <= H8; sq++ {
		initKingAttacksMask(sq)
		initKnightPawnAttacksMask(sq)
		initPawnAttacksMask(sq)

		initBishopRelevantOccupancyBitsMask(sq)
		initRookRelevantOccupancyBitsMask(sq)
	}
}

func initBishopAndRookPopCounts() {
	for sq := A1; sq <= H8; sq++ {
		initBishopRelevantOccupancyBitsPopulationCount(sq)
		initRookRelevantOccupancyBitsPopulationCount(sq)
	}
}

func initAttackMasksForSlidingPieces() {
	for sq := A1; sq <= H8; sq++ {
		initBishopAndRookAttacksMask(sq)
	}
}

func initKingAttacksMask(sq Square) {
	kingAttacksMask[sq] = generateKingAttacksMask(sq)
}

func initKnightPawnAttacksMask(sq Square) {
	knightsAttacksMask[sq] = generateKnightAttacksMask(sq)
}

func initPawnAttacksMask(sq Square) {
	pawnAttacksMask[White][sq] = generatePawnAttacksMask(sq, White)
	pawnAttacksMask[Black][sq] = generatePawnAttacksMask(sq, Black)
}

func initBishopRelevantOccupancyBitsMask(sq Square) {
	bishopRelevantOccupancyBitsMask[sq] = generateBishopRelevantOccupancyBitsMask(sq)
}

func initRookRelevantOccupancyBitsMask(sq Square) {
	rookRelevantOccupancyBitsMask[sq] = generateRookRelevantOccupancyBitsMask(sq)
}

func initBishopRelevantOccupancyBitsPopulationCount(sq Square) {
	BishopRelevantOccupancyBitsPopulationCount[sq] = bishopRelevantOccupancyBitsMask[sq].populationCount()
}

func initRookRelevantOccupancyBitsPopulationCount(sq Square) {
	RookRelevantOccupancyBitsPopulationCount[sq] = rookRelevantOccupancyBitsMask[sq].populationCount()
}

func initBishopAndRookAttacksMask(sq Square) {
	attackBishop := bishopRelevantOccupancyBitsMask[sq]
	attackRook := rookRelevantOccupancyBitsMask[sq]

	bishopPopcount := BishopRelevantOccupancyBitsPopulationCount[sq]
	rookPopcount := RookRelevantOccupancyBitsPopulationCount[sq]
	bishopOccIdx := 1 << bishopPopcount
	rookOccIdx := 1 << rookPopcount

	for i := 0; i < bishopOccIdx; i++ {
		occ := SetOccupancy(i, int(bishopPopcount), attackBishop)
		mIdx := (occ * bishopMagics[sq]) >> (64 - bishopPopcount)
		arr := bishopAttacksMask[sq]
		arr[mIdx] = generateBishopAttacksWithBlockers(sq, occ)
		bishopAttacksMask[sq] = arr
	}

	for i := 0; i < rookOccIdx; i++ {
		occ := SetOccupancy(i, int(rookPopcount), attackRook)
		mIdx := (occ * rookMagics[sq]) >> (64 - rookPopcount)
		arr := rookAttacksMask[sq]
		arr[mIdx] = generateRookAttacksWithBlockers(sq, occ)
		rookAttacksMask[sq] = arr
	}
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

// generateBishopRelevantOccupancyBitsMask generates the bishop relevant occupancy mask that doesn't include edge squares
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

// generateRookRelevantOccupancyBitsMask generates the rook relevant occupancy mask that doesn't include edge squares
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

// SetOccupancy generate occupancy bitboards for a given relevant occupancy bitboard
func SetOccupancy(index, count int, attack bitboard) bitboard {
	occ := bitboardEmpty

	for i := 0; i < count; i++ {
		sq := attack.LS1B()
		attack.clearBit(Square(sq))

		if index&(1<<i) != 0 {
			occ.setBit(Square(sq))
		}
	}

	return occ
}

// generateRandomMagicNumberCandidate generates random bitboard with few non-zero bits for magic number candidates
func generateRandomMagicNumberCandidate() bitboard {
	rnd := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	return bitboard(rnd.Uint64() & rnd.Uint64() & rnd.Uint64())
}

// findMagicNumbers finds the magic numbers with a brute force approach for a square
func findMagicNumbers(sq Square, isBishop bool) bitboard {

	var occ, attacks, usedAttacks [4096]bitboard
	attack := generateRookRelevantOccupancyBitsMask(sq)
	bitCount := attack.populationCount()

	if isBishop {
		attack = generateBishopRelevantOccupancyBitsMask(sq)
		bitCount = attack.populationCount()
	}

	permutations := 1 << bitCount

	for i := 0; i < permutations; i++ {
		occ[i] = SetOccupancy(i, int(bitCount), attack)

		if isBishop {
			attacks[sq] = generateBishopAttacksWithBlockers(sq, occ[sq])
		} else {
			attacks[sq] = generateRookAttacksWithBlockers(sq, occ[sq])
		}
	}

	for randC := 0; randC < 1<<48; randC++ {
		magicNum := generateRandomMagicNumberCandidate()

		if ((attack * magicNum) & bitboardUniverseRanksMask[Rank8]).populationCount() < 6 {
			continue
		}

		usedAttacks = [4096]bitboard{}
		var i int
		var fail bool

		for i = 0; !fail && i < permutations; i++ {
			magicIdx := int((occ[i] * magicNum) >> (64 - bitCount))
			if usedAttacks[magicIdx] == 0 {
				usedAttacks[magicIdx] = attacks[i]
			} else if usedAttacks[magicIdx] != attacks[i] {
				fail = true
			}
		}

		if !fail {
			return magicNum
		}
	}

	panic(fmt.Errorf("failed to find magic numbers"))
}

// getBishopAttacks returns the bishop attack mask with blocker occupancy
func getBishopAttacks(sq Square, occupancy bitboard) bitboard {
	occupancy &= bishopRelevantOccupancyBitsMask[sq]
	occupancy *= bishopMagics[sq]
	occupancy >>= 64 - BishopRelevantOccupancyBitsPopulationCount[sq]
	return bishopAttacksMask[sq][occupancy]
}

// getRookAttacks returns the rook attack mask with blocker occupancy
func getRookAttacks(sq Square, occupancy bitboard) bitboard {
	occupancy &= rookRelevantOccupancyBitsMask[sq]
	occupancy *= rookMagics[sq]
	occupancy >>= 64 - RookRelevantOccupancyBitsPopulationCount[sq]
	return rookAttacksMask[sq][occupancy]
}

// getQueenAttacks returns the queen attack mask with blocker occupancy as combination of rook and bishop attacks
func getQueenAttacks(sq Square, occupancy bitboard) bitboard {
	return getBishopAttacks(sq, occupancy) | getRookAttacks(sq, occupancy)
}
