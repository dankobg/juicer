package engine

import (
	"testing"
)

func TestFiles(t *testing.T) {
	testCases := map[string]struct {
		file  File
		idx   int
		label string
	}{
		"FileA is a file": {file: FileA, label: "a", idx: 0},
		"FileB is b file": {file: FileB, label: "b", idx: 1},
		"FileC is c file": {file: FileC, label: "c", idx: 2},
		"FileD is d file": {file: FileD, label: "d", idx: 3},
		"FileE is e file": {file: FileE, label: "e", idx: 4},
		"FileF is f file": {file: FileF, label: "f", idx: 5},
		"FileG is g file": {file: FileG, label: "g", idx: 6},
		"FileH is h file": {file: FileH, label: "h", idx: 7},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.file.String() != tc.label {
				t.Fatalf("invalid file label: want %s, got %s", tc.label, tc.file)
			}
			if int(tc.file) != tc.idx {
				t.Fatalf("invalid file index: want %d, got %d", tc.idx, int(tc.file))
			}
		})
	}
}

func TestRanks(t *testing.T) {
	testCases := map[string]struct {
		rank  Rank
		idx   int
		label string
	}{
		"Rank1 is first rank":   {rank: Rank1, label: "1", idx: 0},
		"Rank2 is second rank":  {rank: Rank2, label: "2", idx: 1},
		"Rank3 is third rank":   {rank: Rank3, label: "3", idx: 2},
		"Rank4 is fourth rank":  {rank: Rank4, label: "4", idx: 3},
		"Rank5 is fifth rank":   {rank: Rank5, label: "5", idx: 4},
		"Rank6 is sixth rank":   {rank: Rank6, label: "6", idx: 5},
		"Rank7 is seventh rank": {rank: Rank7, label: "7", idx: 6},
		"Rank8 is eighth rank":  {rank: Rank8, label: "8", idx: 7},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.rank.String() != tc.label {
				t.Fatalf("invalid rank label: want %s, got %s", tc.label, tc.rank)
			}
			if int(tc.rank) != tc.idx {
				t.Fatalf("invalid rank index: want %d, got %d", tc.idx, int(tc.rank))
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    File
		wantErr bool
	}{
		"creates a file":           {in: "a", wantErr: false, want: FileA},
		"creates b file":           {in: "b", wantErr: false, want: FileB},
		"creates c file":           {in: "c", wantErr: false, want: FileC},
		"creates d file":           {in: "d", wantErr: false, want: FileD},
		"creates e file":           {in: "e", wantErr: false, want: FileE},
		"creates f file":           {in: "f", wantErr: false, want: FileF},
		"creates g file":           {in: "g", wantErr: false, want: FileG},
		"creates h file":           {in: "h", wantErr: false, want: FileH},
		"fails with invalid empty": {in: "", wantErr: true},
		"fails with invalid j":     {in: "j", wantErr: true},
		"fails with invalid len":   {in: "aa", wantErr: true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			file, err := NewFile(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid file, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && file != tc.want {
				t.Fatalf("invalid file, want %s, got %s", tc.want, file)
			}
		})
	}
}

func TestNewRank(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    Rank
		wantErr bool
	}{
		"creates first rank":       {in: "1", wantErr: false, want: Rank1},
		"creates second rank":      {in: "2", wantErr: false, want: Rank2},
		"creates third rank":       {in: "3", wantErr: false, want: Rank3},
		"creates fourth rank":      {in: "4", wantErr: false, want: Rank4},
		"creates fifth rank":       {in: "5", wantErr: false, want: Rank5},
		"creates sixth rank":       {in: "6", wantErr: false, want: Rank6},
		"creates seventh rank":     {in: "7", wantErr: false, want: Rank7},
		"creates eighth rank":      {in: "8", wantErr: false, want: Rank8},
		"fails with invalid empty": {in: "", wantErr: true},
		"fails with invalid 9":     {in: "9", wantErr: true},
		"fails with invalid len":   {in: "69", wantErr: true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rank, err := NewRank(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid rank, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && rank != tc.want {
				t.Fatalf("invalid rank, want %s, got %s", tc.want, rank)
			}
		})
	}
}

func TestSquares(t *testing.T) {
	testCases := map[string]struct {
		sq    Square
		idx   int
		label string
		color Color
		file  File
		rank  Rank
	}{
		"Square A1": {sq: A1, label: "a1", idx: 0, color: Black, file: FileA, rank: Rank1},
		"Square B1": {sq: B1, label: "b1", idx: 1, color: White, file: FileB, rank: Rank1},
		"Square C1": {sq: C1, label: "c1", idx: 2, color: Black, file: FileC, rank: Rank1},
		"Square D1": {sq: D1, label: "d1", idx: 3, color: White, file: FileD, rank: Rank1},
		"Square E1": {sq: E1, label: "e1", idx: 4, color: Black, file: FileE, rank: Rank1},
		"Square F1": {sq: F1, label: "f1", idx: 5, color: White, file: FileF, rank: Rank1},
		"Square G1": {sq: G1, label: "g1", idx: 6, color: Black, file: FileG, rank: Rank1},
		"Square H1": {sq: H1, label: "h1", idx: 7, color: White, file: FileH, rank: Rank1},
		"Square A2": {sq: A2, label: "a2", idx: 8, color: White, file: FileA, rank: Rank2},
		"Square B2": {sq: B2, label: "b2", idx: 9, color: Black, file: FileB, rank: Rank2},
		"Square C2": {sq: C2, label: "c2", idx: 10, color: White, file: FileC, rank: Rank2},
		"Square D2": {sq: D2, label: "d2", idx: 11, color: Black, file: FileD, rank: Rank2},
		"Square E2": {sq: E2, label: "e2", idx: 12, color: White, file: FileE, rank: Rank2},
		"Square F2": {sq: F2, label: "f2", idx: 13, color: Black, file: FileF, rank: Rank2},
		"Square G2": {sq: G2, label: "g2", idx: 14, color: White, file: FileG, rank: Rank2},
		"Square H2": {sq: H2, label: "h2", idx: 15, color: Black, file: FileH, rank: Rank2},
		"Square A3": {sq: A3, label: "a3", idx: 16, color: Black, file: FileA, rank: Rank3},
		"Square B3": {sq: B3, label: "b3", idx: 17, color: White, file: FileB, rank: Rank3},
		"Square C3": {sq: C3, label: "c3", idx: 18, color: Black, file: FileC, rank: Rank3},
		"Square D3": {sq: D3, label: "d3", idx: 19, color: White, file: FileD, rank: Rank3},
		"Square E3": {sq: E3, label: "e3", idx: 20, color: Black, file: FileE, rank: Rank3},
		"Square F3": {sq: F3, label: "f3", idx: 21, color: White, file: FileF, rank: Rank3},
		"Square G3": {sq: G3, label: "g3", idx: 22, color: Black, file: FileG, rank: Rank3},
		"Square H3": {sq: H3, label: "h3", idx: 23, color: White, file: FileH, rank: Rank3},
		"Square A4": {sq: A4, label: "a4", idx: 24, color: White, file: FileA, rank: Rank4},
		"Square B4": {sq: B4, label: "b4", idx: 25, color: Black, file: FileB, rank: Rank4},
		"Square C4": {sq: C4, label: "c4", idx: 26, color: White, file: FileC, rank: Rank4},
		"Square D4": {sq: D4, label: "d4", idx: 27, color: Black, file: FileD, rank: Rank4},
		"Square E4": {sq: E4, label: "e4", idx: 28, color: White, file: FileE, rank: Rank4},
		"Square F4": {sq: F4, label: "f4", idx: 29, color: Black, file: FileF, rank: Rank4},
		"Square G4": {sq: G4, label: "g4", idx: 30, color: White, file: FileG, rank: Rank4},
		"Square H4": {sq: H4, label: "h4", idx: 31, color: Black, file: FileH, rank: Rank4},
		"Square A5": {sq: A5, label: "a5", idx: 32, color: Black, file: FileA, rank: Rank5},
		"Square B5": {sq: B5, label: "b5", idx: 33, color: White, file: FileB, rank: Rank5},
		"Square C5": {sq: C5, label: "c5", idx: 34, color: Black, file: FileC, rank: Rank5},
		"Square D5": {sq: D5, label: "d5", idx: 35, color: White, file: FileD, rank: Rank5},
		"Square E5": {sq: E5, label: "e5", idx: 36, color: Black, file: FileE, rank: Rank5},
		"Square F5": {sq: F5, label: "f5", idx: 37, color: White, file: FileF, rank: Rank5},
		"Square G5": {sq: G5, label: "g5", idx: 38, color: Black, file: FileG, rank: Rank5},
		"Square H5": {sq: H5, label: "h5", idx: 39, color: White, file: FileH, rank: Rank5},
		"Square A6": {sq: A6, label: "a6", idx: 40, color: White, file: FileA, rank: Rank6},
		"Square B6": {sq: B6, label: "b6", idx: 41, color: Black, file: FileB, rank: Rank6},
		"Square C6": {sq: C6, label: "c6", idx: 42, color: White, file: FileC, rank: Rank6},
		"Square D6": {sq: D6, label: "d6", idx: 43, color: Black, file: FileD, rank: Rank6},
		"Square E6": {sq: E6, label: "e6", idx: 44, color: White, file: FileE, rank: Rank6},
		"Square F6": {sq: F6, label: "f6", idx: 45, color: Black, file: FileF, rank: Rank6},
		"Square G6": {sq: G6, label: "g6", idx: 46, color: White, file: FileG, rank: Rank6},
		"Square H6": {sq: H6, label: "h6", idx: 47, color: Black, file: FileH, rank: Rank6},
		"Square A7": {sq: A7, label: "a7", idx: 48, color: Black, file: FileA, rank: Rank7},
		"Square B7": {sq: B7, label: "b7", idx: 49, color: White, file: FileB, rank: Rank7},
		"Square C7": {sq: C7, label: "c7", idx: 50, color: Black, file: FileC, rank: Rank7},
		"Square D7": {sq: D7, label: "d7", idx: 51, color: White, file: FileD, rank: Rank7},
		"Square E7": {sq: E7, label: "e7", idx: 52, color: Black, file: FileE, rank: Rank7},
		"Square F7": {sq: F7, label: "f7", idx: 53, color: White, file: FileF, rank: Rank7},
		"Square G7": {sq: G7, label: "g7", idx: 54, color: Black, file: FileG, rank: Rank7},
		"Square H7": {sq: H7, label: "h7", idx: 55, color: White, file: FileH, rank: Rank7},
		"Square A8": {sq: A8, label: "a8", idx: 56, color: White, file: FileA, rank: Rank8},
		"Square B8": {sq: B8, label: "b8", idx: 57, color: Black, file: FileB, rank: Rank8},
		"Square C8": {sq: C8, label: "c8", idx: 58, color: White, file: FileC, rank: Rank8},
		"Square D8": {sq: D8, label: "d8", idx: 59, color: Black, file: FileD, rank: Rank8},
		"Square E8": {sq: E8, label: "e8", idx: 60, color: White, file: FileE, rank: Rank8},
		"Square F8": {sq: F8, label: "f8", idx: 61, color: Black, file: FileF, rank: Rank8},
		"Square G8": {sq: G8, label: "g8", idx: 62, color: White, file: FileG, rank: Rank8},
		"Square H8": {sq: H8, label: "h8", idx: 63, color: Black, file: FileH, rank: Rank8},
	}
	t.Run("valid SquareNone", func(t *testing.T) {
		if SquareNone >= A1 && SquareNone <= H8 {
			t.Fatalf("invalid SquareNone")
		}
	})

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.sq.String() != tc.label {
				t.Fatalf("invalid square label: want %s, got %s", tc.label, tc.sq)
			}
			if int(tc.sq) != tc.idx {
				t.Fatalf("invalid square index: want %d, got %d", tc.idx, int(tc.sq))
			}
			if tc.file != tc.sq.File() {
				t.Fatalf("invalid square file: want %d, got %d", tc.file, tc.sq.File())
			}
			if tc.rank != tc.sq.Rank() {
				t.Fatalf("invalid square rank: want %d, got %d", tc.rank, tc.sq.Rank())
			}
			if tc.color != tc.sq.Color() {
				t.Fatalf("invalid square color: want %d, got %d", tc.color, tc.sq.Color())
			}
		})
	}
}

func TestNewSquare(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    Square
		wantErr bool
	}{
		"creates A1 square":        {in: "a1", want: A1, wantErr: false},
		"creates B1 square":        {in: "b1", want: B1, wantErr: false},
		"creates C1 square":        {in: "c1", want: C1, wantErr: false},
		"creates D1 square":        {in: "d1", want: D1, wantErr: false},
		"creates E1 square":        {in: "e1", want: E1, wantErr: false},
		"creates F1 square":        {in: "f1", want: F1, wantErr: false},
		"creates G1 square":        {in: "g1", want: G1, wantErr: false},
		"creates H1 square":        {in: "h1", want: H1, wantErr: false},
		"creates A2 square":        {in: "a2", want: A2, wantErr: false},
		"creates B2 square":        {in: "b2", want: B2, wantErr: false},
		"creates C2 square":        {in: "c2", want: C2, wantErr: false},
		"creates D2 square":        {in: "d2", want: D2, wantErr: false},
		"creates E2 square":        {in: "e2", want: E2, wantErr: false},
		"creates F2 square":        {in: "f2", want: F2, wantErr: false},
		"creates G2 square":        {in: "g2", want: G2, wantErr: false},
		"creates H2 square":        {in: "h2", want: H2, wantErr: false},
		"creates A3 square":        {in: "a3", want: A3, wantErr: false},
		"creates B3 square":        {in: "b3", want: B3, wantErr: false},
		"creates C3 square":        {in: "c3", want: C3, wantErr: false},
		"creates D3 square":        {in: "d3", want: D3, wantErr: false},
		"creates E3 square":        {in: "e3", want: E3, wantErr: false},
		"creates F3 square":        {in: "f3", want: F3, wantErr: false},
		"creates G3 square":        {in: "g3", want: G3, wantErr: false},
		"creates H3 square":        {in: "h3", want: H3, wantErr: false},
		"creates A4 square":        {in: "a4", want: A4, wantErr: false},
		"creates B4 square":        {in: "b4", want: B4, wantErr: false},
		"creates C4 square":        {in: "c4", want: C4, wantErr: false},
		"creates D4 square":        {in: "d4", want: D4, wantErr: false},
		"creates E4 square":        {in: "e4", want: E4, wantErr: false},
		"creates F4 square":        {in: "f4", want: F4, wantErr: false},
		"creates G4 square":        {in: "g4", want: G4, wantErr: false},
		"creates H4 square":        {in: "h4", want: H4, wantErr: false},
		"creates A5 square":        {in: "a5", want: A5, wantErr: false},
		"creates B5 square":        {in: "b5", want: B5, wantErr: false},
		"creates C5 square":        {in: "c5", want: C5, wantErr: false},
		"creates D5 square":        {in: "d5", want: D5, wantErr: false},
		"creates E5 square":        {in: "e5", want: E5, wantErr: false},
		"creates F5 square":        {in: "f5", want: F5, wantErr: false},
		"creates G5 square":        {in: "g5", want: G5, wantErr: false},
		"creates H5 square":        {in: "h5", want: H5, wantErr: false},
		"creates A6 square":        {in: "a6", want: A6, wantErr: false},
		"creates B6 square":        {in: "b6", want: B6, wantErr: false},
		"creates C6 square":        {in: "c6", want: C6, wantErr: false},
		"creates D6 square":        {in: "d6", want: D6, wantErr: false},
		"creates E6 square":        {in: "e6", want: E6, wantErr: false},
		"creates F6 square":        {in: "f6", want: F6, wantErr: false},
		"creates G6 square":        {in: "g6", want: G6, wantErr: false},
		"creates H6 square":        {in: "h6", want: H6, wantErr: false},
		"creates A7 square":        {in: "a7", want: A7, wantErr: false},
		"creates B7 square":        {in: "b7", want: B7, wantErr: false},
		"creates C7 square":        {in: "c7", want: C7, wantErr: false},
		"creates D7 square":        {in: "d7", want: D7, wantErr: false},
		"creates E7 square":        {in: "e7", want: E7, wantErr: false},
		"creates F7 square":        {in: "f7", want: F7, wantErr: false},
		"creates G7 square":        {in: "g7", want: G7, wantErr: false},
		"creates H7 square":        {in: "h7", want: H7, wantErr: false},
		"creates A8 square":        {in: "a8", want: A8, wantErr: false},
		"creates B8 square":        {in: "b8", want: B8, wantErr: false},
		"creates C8 square":        {in: "c8", want: C8, wantErr: false},
		"creates D8 square":        {in: "d8", want: D8, wantErr: false},
		"creates E8 square":        {in: "e8", want: E8, wantErr: false},
		"creates F8 square":        {in: "f8", want: F8, wantErr: false},
		"creates G8 square":        {in: "g8", want: G8, wantErr: false},
		"creates H8 square":        {in: "h8", want: H8, wantErr: false},
		"fails with invalid empty": {in: "", wantErr: true},
		"fails with invalid file":  {in: "j5", wantErr: true},
		"fails with invalid rank":  {in: "g9", wantErr: true},
		"fails with len short":     {in: "h", wantErr: true},
		"fails with len long":      {in: "h51", wantErr: true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sq, err := NewSquareFromCoord(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid square, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && sq != tc.want {
				t.Fatalf("invalid square: want %s, got %s", tc.want, sq)
			}
		})
	}
}
