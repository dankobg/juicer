package juicer

import (
	"testing"
)

func TestGenerateKingAttacksMask(t *testing.T) {
	testCases := map[string]struct {
		sq   Square
		want bitboard
	}{
		"A1": {sq: A1, want: 0x302},
		"B1": {sq: B1, want: 0x705},
		"C1": {sq: C1, want: 0xE0A},
		"D1": {sq: D1, want: 0x1C14},
		"E1": {sq: E1, want: 0x3828},
		"F1": {sq: F1, want: 0x7050},
		"G1": {sq: G1, want: 0xE0A0},
		"H1": {sq: H1, want: 0xC040},
		"A2": {sq: A2, want: 0x30203},
		"B2": {sq: B2, want: 0x70507},
		"C2": {sq: C2, want: 0xE0A0E},
		"D2": {sq: D2, want: 0x1C141C},
		"E2": {sq: E2, want: 0x382838},
		"F2": {sq: F2, want: 0x705070},
		"G2": {sq: G2, want: 0xE0A0E0},
		"H2": {sq: H2, want: 0xC040C0},
		"A3": {sq: A3, want: 0x3020300},
		"B3": {sq: B3, want: 0x7050700},
		"C3": {sq: C3, want: 0xE0A0E00},
		"D3": {sq: D3, want: 0x1C141C00},
		"E3": {sq: E3, want: 0x38283800},
		"F3": {sq: F3, want: 0x70507000},
		"G3": {sq: G3, want: 0xE0A0E000},
		"H3": {sq: H3, want: 0xC040C000},
		"A4": {sq: A4, want: 0x302030000},
		"B4": {sq: B4, want: 0x705070000},
		"C4": {sq: C4, want: 0xE0A0E0000},
		"D4": {sq: D4, want: 0x1C141C0000},
		"E4": {sq: E4, want: 0x3828380000},
		"F4": {sq: F4, want: 0x7050700000},
		"G4": {sq: G4, want: 0xE0A0E00000},
		"H4": {sq: H4, want: 0xC040C00000},
		"A5": {sq: A5, want: 0x30203000000},
		"B5": {sq: B5, want: 0x70507000000},
		"C5": {sq: C5, want: 0xE0A0E000000},
		"D5": {sq: D5, want: 0x1C141C000000},
		"E5": {sq: E5, want: 0x382838000000},
		"F5": {sq: F5, want: 0x705070000000},
		"G5": {sq: G5, want: 0xE0A0E0000000},
		"H5": {sq: H5, want: 0xC040C0000000},
		"A6": {sq: A6, want: 0x3020300000000},
		"B6": {sq: B6, want: 0x7050700000000},
		"C6": {sq: C6, want: 0xE0A0E00000000},
		"D6": {sq: D6, want: 0x1C141C00000000},
		"E6": {sq: E6, want: 0x38283800000000},
		"F6": {sq: F6, want: 0x70507000000000},
		"G6": {sq: G6, want: 0xE0A0E000000000},
		"H6": {sq: H6, want: 0xC040C000000000},
		"A7": {sq: A7, want: 0x302030000000000},
		"B7": {sq: B7, want: 0x705070000000000},
		"C7": {sq: C7, want: 0xE0A0E0000000000},
		"D7": {sq: D7, want: 0x1C141C0000000000},
		"E7": {sq: E7, want: 0x3828380000000000},
		"F7": {sq: F7, want: 0x7050700000000000},
		"G7": {sq: G7, want: 0xE0A0E00000000000},
		"H7": {sq: H7, want: 0xC040C00000000000},
		"A8": {sq: A8, want: 0x203000000000000},
		"B8": {sq: B8, want: 0x507000000000000},
		"C8": {sq: C8, want: 0xA0E000000000000},
		"D8": {sq: D8, want: 0x141C000000000000},
		"E8": {sq: E8, want: 0x2838000000000000},
		"F8": {sq: F8, want: 0x5070000000000000},
		"G8": {sq: G8, want: 0xA0E0000000000000},
		"H8": {sq: H8, want: 0x40C0000000000000},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			attacks := generateKingAttacksMask(tc.sq)

			if attacks != tc.want {
				t.Fatalf("invalid king attacks mask: want %s, got %s", tc.want, attacks)
			}
		})
	}
}

func TestGenerateKnightAttacksMask(t *testing.T) {
	testCases := map[string]struct {
		sq   Square
		want bitboard
	}{
		"A1": {sq: A1, want: 0x20400},
		"B1": {sq: B1, want: 0x50800},
		"C1": {sq: C1, want: 0xa1100},
		"D1": {sq: D1, want: 0x142200},
		"E1": {sq: E1, want: 0x284400},
		"F1": {sq: F1, want: 0x508800},
		"G1": {sq: G1, want: 0xa01000},
		"H1": {sq: H1, want: 0x402000},
		"A2": {sq: A2, want: 0x2040004},
		"B2": {sq: B2, want: 0x5080008},
		"C2": {sq: C2, want: 0xa110011},
		"D2": {sq: D2, want: 0x14220022},
		"E2": {sq: E2, want: 0x28440044},
		"F2": {sq: F2, want: 0x50880088},
		"G2": {sq: G2, want: 0xa0100010},
		"H2": {sq: H2, want: 0x40200020},
		"A3": {sq: A3, want: 0x204000402},
		"B3": {sq: B3, want: 0x508000805},
		"C3": {sq: C3, want: 0xa1100110a},
		"D3": {sq: D3, want: 0x1422002214},
		"E3": {sq: E3, want: 0x2844004428},
		"F3": {sq: F3, want: 0x5088008850},
		"G3": {sq: G3, want: 0xa0100010a0},
		"H3": {sq: H3, want: 0x4020002040},
		"A4": {sq: A4, want: 0x20400040200},
		"B4": {sq: B4, want: 0x50800080500},
		"C4": {sq: C4, want: 0xa1100110a00},
		"D4": {sq: D4, want: 0x142200221400},
		"E4": {sq: E4, want: 0x284400442800},
		"F4": {sq: F4, want: 0x508800885000},
		"G4": {sq: G4, want: 0xa0100010a000},
		"H4": {sq: H4, want: 0x402000204000},
		"A5": {sq: A5, want: 0x2040004020000},
		"B5": {sq: B5, want: 0x5080008050000},
		"C5": {sq: C5, want: 0xa1100110a0000},
		"D5": {sq: D5, want: 0x14220022140000},
		"E5": {sq: E5, want: 0x28440044280000},
		"F5": {sq: F5, want: 0x50880088500000},
		"G5": {sq: G5, want: 0xa0100010a00000},
		"H5": {sq: H5, want: 0x40200020400000},
		"A6": {sq: A6, want: 0x204000402000000},
		"B6": {sq: B6, want: 0x508000805000000},
		"C6": {sq: C6, want: 0xa1100110a000000},
		"D6": {sq: D6, want: 0x1422002214000000},
		"E6": {sq: E6, want: 0x2844004428000000},
		"F6": {sq: F6, want: 0x5088008850000000},
		"G6": {sq: G6, want: 0xa0100010a0000000},
		"H6": {sq: H6, want: 0x4020002040000000},
		"A7": {sq: A7, want: 0x400040200000000},
		"B7": {sq: B7, want: 0x800080500000000},
		"C7": {sq: C7, want: 0x1100110a00000000},
		"D7": {sq: D7, want: 0x2200221400000000},
		"E7": {sq: E7, want: 0x4400442800000000},
		"F7": {sq: F7, want: 0x8800885000000000},
		"G7": {sq: G7, want: 0x100010a000000000},
		"H7": {sq: H7, want: 0x2000204000000000},
		"A8": {sq: A8, want: 0x4020000000000},
		"B8": {sq: B8, want: 0x8050000000000},
		"C8": {sq: C8, want: 0x110a0000000000},
		"D8": {sq: D8, want: 0x22140000000000},
		"E8": {sq: E8, want: 0x44280000000000},
		"F8": {sq: F8, want: 0x88500000000000},
		"G8": {sq: G8, want: 0x10a00000000000},
		"H8": {sq: H8, want: 0x20400000000000},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			attacks := generateKnightAttacksMask(tc.sq)

			if attacks != tc.want {
				t.Fatalf("invalid knight attacks mask: want %s, got %s", tc.want, attacks)
			}
		})
	}
}

func TestGeneratePawnAttacksMask(t *testing.T) {
	testCases := map[string]struct {
		sq        Square
		wantWhite bitboard
		wantBlack bitboard
	}{
		"A1": {sq: A1, wantWhite: 0x200, wantBlack: 0x0},
		"B1": {sq: B1, wantWhite: 0x500, wantBlack: 0x0},
		"C1": {sq: C1, wantWhite: 0xa00, wantBlack: 0x0},
		"D1": {sq: D1, wantWhite: 0x1400, wantBlack: 0x0},
		"E1": {sq: E1, wantWhite: 0x2800, wantBlack: 0x0},
		"F1": {sq: F1, wantWhite: 0x5000, wantBlack: 0x0},
		"G1": {sq: G1, wantWhite: 0xa000, wantBlack: 0x0},
		"H1": {sq: H1, wantWhite: 0x4000, wantBlack: 0x0},
		"A2": {sq: A2, wantWhite: 0x20000, wantBlack: 0x2},
		"B2": {sq: B2, wantWhite: 0x50000, wantBlack: 0x5},
		"C2": {sq: C2, wantWhite: 0xa0000, wantBlack: 0xa},
		"D2": {sq: D2, wantWhite: 0x140000, wantBlack: 0x14},
		"E2": {sq: E2, wantWhite: 0x280000, wantBlack: 0x28},
		"F2": {sq: F2, wantWhite: 0x500000, wantBlack: 0x50},
		"G2": {sq: G2, wantWhite: 0xa00000, wantBlack: 0xa0},
		"H2": {sq: H2, wantWhite: 0x400000, wantBlack: 0x40},
		"A3": {sq: A3, wantWhite: 0x2000000, wantBlack: 0x200},
		"B3": {sq: B3, wantWhite: 0x5000000, wantBlack: 0x500},
		"C3": {sq: C3, wantWhite: 0xa000000, wantBlack: 0xa00},
		"D3": {sq: D3, wantWhite: 0x14000000, wantBlack: 0x1400},
		"E3": {sq: E3, wantWhite: 0x28000000, wantBlack: 0x2800},
		"F3": {sq: F3, wantWhite: 0x50000000, wantBlack: 0x5000},
		"G3": {sq: G3, wantWhite: 0xa0000000, wantBlack: 0xa000},
		"H3": {sq: H3, wantWhite: 0x40000000, wantBlack: 0x4000},
		"A4": {sq: A4, wantWhite: 0x200000000, wantBlack: 0x20000},
		"B4": {sq: B4, wantWhite: 0x500000000, wantBlack: 0x50000},
		"C4": {sq: C4, wantWhite: 0xa00000000, wantBlack: 0xa0000},
		"D4": {sq: D4, wantWhite: 0x1400000000, wantBlack: 0x140000},
		"E4": {sq: E4, wantWhite: 0x2800000000, wantBlack: 0x280000},
		"F4": {sq: F4, wantWhite: 0x5000000000, wantBlack: 0x500000},
		"G4": {sq: G4, wantWhite: 0xa000000000, wantBlack: 0xa00000},
		"H4": {sq: H4, wantWhite: 0x4000000000, wantBlack: 0x400000},
		"A5": {sq: A5, wantWhite: 0x20000000000, wantBlack: 0x2000000},
		"B5": {sq: B5, wantWhite: 0x50000000000, wantBlack: 0x5000000},
		"C5": {sq: C5, wantWhite: 0xa0000000000, wantBlack: 0xa000000},
		"D5": {sq: D5, wantWhite: 0x140000000000, wantBlack: 0x14000000},
		"E5": {sq: E5, wantWhite: 0x280000000000, wantBlack: 0x28000000},
		"F5": {sq: F5, wantWhite: 0x500000000000, wantBlack: 0x50000000},
		"G5": {sq: G5, wantWhite: 0xa00000000000, wantBlack: 0xa0000000},
		"H5": {sq: H5, wantWhite: 0x400000000000, wantBlack: 0x40000000},
		"A6": {sq: A6, wantWhite: 0x2000000000000, wantBlack: 0x200000000},
		"B6": {sq: B6, wantWhite: 0x5000000000000, wantBlack: 0x500000000},
		"C6": {sq: C6, wantWhite: 0xa000000000000, wantBlack: 0xa00000000},
		"D6": {sq: D6, wantWhite: 0x14000000000000, wantBlack: 0x1400000000},
		"E6": {sq: E6, wantWhite: 0x28000000000000, wantBlack: 0x2800000000},
		"F6": {sq: F6, wantWhite: 0x50000000000000, wantBlack: 0x5000000000},
		"G6": {sq: G6, wantWhite: 0xa0000000000000, wantBlack: 0xa000000000},
		"H6": {sq: H6, wantWhite: 0x40000000000000, wantBlack: 0x4000000000},
		"A7": {sq: A7, wantWhite: 0x200000000000000, wantBlack: 0x20000000000},
		"B7": {sq: B7, wantWhite: 0x500000000000000, wantBlack: 0x50000000000},
		"C7": {sq: C7, wantWhite: 0xa00000000000000, wantBlack: 0xa0000000000},
		"D7": {sq: D7, wantWhite: 0x1400000000000000, wantBlack: 0x140000000000},
		"E7": {sq: E7, wantWhite: 0x2800000000000000, wantBlack: 0x280000000000},
		"F7": {sq: F7, wantWhite: 0x5000000000000000, wantBlack: 0x500000000000},
		"G7": {sq: G7, wantWhite: 0xa000000000000000, wantBlack: 0xa00000000000},
		"H7": {sq: H7, wantWhite: 0x4000000000000000, wantBlack: 0x400000000000},
		"A8": {sq: A8, wantWhite: 0x0, wantBlack: 0x2000000000000},
		"B8": {sq: B8, wantWhite: 0x0, wantBlack: 0x5000000000000},
		"C8": {sq: C8, wantWhite: 0x0, wantBlack: 0xa000000000000},
		"D8": {sq: D8, wantWhite: 0x0, wantBlack: 0x14000000000000},
		"E8": {sq: E8, wantWhite: 0x0, wantBlack: 0x28000000000000},
		"F8": {sq: F8, wantWhite: 0x0, wantBlack: 0x50000000000000},
		"G8": {sq: G8, wantWhite: 0x0, wantBlack: 0xa0000000000000},
		"H8": {sq: H8, wantWhite: 0x0, wantBlack: 0x40000000000000},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			attacksWhite := generatePawnAttacksMask(tc.sq, White)
			attacksBlack := generatePawnAttacksMask(tc.sq, Black)

			if attacksWhite != tc.wantWhite {
				t.Fatalf("invalid white pawn attacks mask: want %s, got %s", tc.wantWhite, attacksWhite)
			}
			if attacksBlack != tc.wantBlack {
				t.Fatalf("invalid black pawn attacks mask: want %s, got %s", tc.wantBlack, attacksBlack)
			}
		})
	}
}

func TestGenerateBishopRelevantOccupancyBitsMask(t *testing.T) {
	testCases := map[string]struct {
		sq   Square
		want bitboard
	}{
		"A1": {sq: A1, want: 0x40201008040200},
		"B1": {sq: B1, want: 0x402010080400},
		"C1": {sq: C1, want: 0x4020100a00},
		"D1": {sq: D1, want: 0x40221400},
		"E1": {sq: E1, want: 0x2442800},
		"F1": {sq: F1, want: 0x204085000},
		"G1": {sq: G1, want: 0x20408102000},
		"H1": {sq: H1, want: 0x2040810204000},
		"A2": {sq: A2, want: 0x20100804020000},
		"B2": {sq: B2, want: 0x40201008040000},
		"C2": {sq: C2, want: 0x4020100a0000},
		"D2": {sq: D2, want: 0x4022140000},
		"E2": {sq: E2, want: 0x244280000},
		"F2": {sq: F2, want: 0x20408500000},
		"G2": {sq: G2, want: 0x2040810200000},
		"H2": {sq: H2, want: 0x4081020400000},
		"A3": {sq: A3, want: 0x10080402000200},
		"B3": {sq: B3, want: 0x20100804000400},
		"C3": {sq: C3, want: 0x4020100a000a00},
		"D3": {sq: D3, want: 0x402214001400},
		"E3": {sq: E3, want: 0x24428002800},
		"F3": {sq: F3, want: 0x2040850005000},
		"G3": {sq: G3, want: 0x4081020002000},
		"H3": {sq: H3, want: 0x8102040004000},
		"A4": {sq: A4, want: 0x8040200020400},
		"B4": {sq: B4, want: 0x10080400040800},
		"C4": {sq: C4, want: 0x20100a000a1000},
		"D4": {sq: D4, want: 0x40221400142200},
		"E4": {sq: E4, want: 0x2442800284400},
		"F4": {sq: F4, want: 0x4085000500800},
		"G4": {sq: G4, want: 0x8102000201000},
		"H4": {sq: H4, want: 0x10204000402000},
		"A5": {sq: A5, want: 0x4020002040800},
		"B5": {sq: B5, want: 0x8040004081000},
		"C5": {sq: C5, want: 0x100a000a102000},
		"D5": {sq: D5, want: 0x22140014224000},
		"E5": {sq: E5, want: 0x44280028440200},
		"F5": {sq: F5, want: 0x8500050080400},
		"G5": {sq: G5, want: 0x10200020100800},
		"H5": {sq: H5, want: 0x20400040201000},
		"A6": {sq: A6, want: 0x2000204081000},
		"B6": {sq: B6, want: 0x4000408102000},
		"C6": {sq: C6, want: 0xa000a10204000},
		"D6": {sq: D6, want: 0x14001422400000},
		"E6": {sq: E6, want: 0x28002844020000},
		"F6": {sq: F6, want: 0x50005008040200},
		"G6": {sq: G6, want: 0x20002010080400},
		"H6": {sq: H6, want: 0x40004020100800},
		"A7": {sq: A7, want: 0x20408102000},
		"B7": {sq: B7, want: 0x40810204000},
		"C7": {sq: C7, want: 0xa1020400000},
		"D7": {sq: D7, want: 0x142240000000},
		"E7": {sq: E7, want: 0x284402000000},
		"F7": {sq: F7, want: 0x500804020000},
		"G7": {sq: G7, want: 0x201008040200},
		"H7": {sq: H7, want: 0x402010080400},
		"A8": {sq: A8, want: 0x2040810204000},
		"B8": {sq: B8, want: 0x4081020400000},
		"C8": {sq: C8, want: 0xa102040000000},
		"D8": {sq: D8, want: 0x14224000000000},
		"E8": {sq: E8, want: 0x28440200000000},
		"F8": {sq: F8, want: 0x50080402000000},
		"G8": {sq: G8, want: 0x20100804020000},
		"H8": {sq: H8, want: 0x40201008040200},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			attacks := generateBishopRelevantOccupancyBitsMask(tc.sq)

			if attacks != tc.want {
				t.Fatalf("invalid bishop relevant occupancy bits mask: want %s, got %s", tc.want, attacks)
			}
		})
	}
}

func TestGenerateRookRelevantOccupancyBitsMask(t *testing.T) {
	testCases := map[string]struct {
		sq   Square
		want bitboard
	}{
		"A1": {sq: A1, want: 0x101010101017e},
		"B1": {sq: B1, want: 0x202020202027c},
		"C1": {sq: C1, want: 0x404040404047a},
		"D1": {sq: D1, want: 0x8080808080876},
		"E1": {sq: E1, want: 0x1010101010106e},
		"F1": {sq: F1, want: 0x2020202020205e},
		"G1": {sq: G1, want: 0x4040404040403e},
		"H1": {sq: H1, want: 0x8080808080807e},
		"A2": {sq: A2, want: 0x1010101017e00},
		"B2": {sq: B2, want: 0x2020202027c00},
		"C2": {sq: C2, want: 0x4040404047a00},
		"D2": {sq: D2, want: 0x8080808087600},
		"E2": {sq: E2, want: 0x10101010106e00},
		"F2": {sq: F2, want: 0x20202020205e00},
		"G2": {sq: G2, want: 0x40404040403e00},
		"H2": {sq: H2, want: 0x80808080807e00},
		"A3": {sq: A3, want: 0x10101017e0100},
		"B3": {sq: B3, want: 0x20202027c0200},
		"C3": {sq: C3, want: 0x40404047a0400},
		"D3": {sq: D3, want: 0x8080808760800},
		"E3": {sq: E3, want: 0x101010106e1000},
		"F3": {sq: F3, want: 0x202020205e2000},
		"G3": {sq: G3, want: 0x404040403e4000},
		"H3": {sq: H3, want: 0x808080807e8000},
		"A4": {sq: A4, want: 0x101017e010100},
		"B4": {sq: B4, want: 0x202027c020200},
		"C4": {sq: C4, want: 0x404047a040400},
		"D4": {sq: D4, want: 0x8080876080800},
		"E4": {sq: E4, want: 0x1010106e101000},
		"F4": {sq: F4, want: 0x2020205e202000},
		"G4": {sq: G4, want: 0x4040403e404000},
		"H4": {sq: H4, want: 0x8080807e808000},
		"A5": {sq: A5, want: 0x1017e01010100},
		"B5": {sq: B5, want: 0x2027c02020200},
		"C5": {sq: C5, want: 0x4047a04040400},
		"D5": {sq: D5, want: 0x8087608080800},
		"E5": {sq: E5, want: 0x10106e10101000},
		"F5": {sq: F5, want: 0x20205e20202000},
		"G5": {sq: G5, want: 0x40403e40404000},
		"H5": {sq: H5, want: 0x80807e80808000},
		"A6": {sq: A6, want: 0x17e0101010100},
		"B6": {sq: B6, want: 0x27c0202020200},
		"C6": {sq: C6, want: 0x47a0404040400},
		"D6": {sq: D6, want: 0x8760808080800},
		"E6": {sq: E6, want: 0x106e1010101000},
		"F6": {sq: F6, want: 0x205e2020202000},
		"G6": {sq: G6, want: 0x403e4040404000},
		"H6": {sq: H6, want: 0x807e8080808000},
		"A7": {sq: A7, want: 0x7e010101010100},
		"B7": {sq: B7, want: 0x7c020202020200},
		"C7": {sq: C7, want: 0x7a040404040400},
		"D7": {sq: D7, want: 0x76080808080800},
		"E7": {sq: E7, want: 0x6e101010101000},
		"F7": {sq: F7, want: 0x5e202020202000},
		"G7": {sq: G7, want: 0x3e404040404000},
		"H7": {sq: H7, want: 0x7e808080808000},
		"A8": {sq: A8, want: 0x7e01010101010100},
		"B8": {sq: B8, want: 0x7c02020202020200},
		"C8": {sq: C8, want: 0x7a04040404040400},
		"D8": {sq: D8, want: 0x7608080808080800},
		"E8": {sq: E8, want: 0x6e10101010101000},
		"F8": {sq: F8, want: 0x5e20202020202000},
		"G8": {sq: G8, want: 0x3e40404040404000},
		"H8": {sq: H8, want: 0x7e80808080808000},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			attacks := generateRookRelevantOccupancyBitsMask(tc.sq)

			if attacks != tc.want {
				t.Fatalf("invalid rook relevant occupancy bits mask: want %s, got %s", tc.want, attacks)
			}
		})
	}
}
