package juicer

import "testing"

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
				t.Fatalf("want %s, got %s", tc.want, attacks)
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
				t.Fatalf("want %s, got %s", tc.want, attacks)
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
				t.Fatalf("wantWhite %s, got %s", tc.wantWhite, attacksWhite)
			}

			if attacksBlack != tc.wantBlack {
				t.Fatalf("wantBlack %s, got %s", tc.wantBlack, attacksBlack)
			}
		})
	}
}
