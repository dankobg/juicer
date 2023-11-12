package juicer

import (
	"testing"
)

func TestSquareFiles(t *testing.T) {
	testCases := map[string]struct {
		file File
		want string
	}{
		"FileA is a file": {file: FileA, want: "a"},
		"FileB is b file": {file: FileB, want: "b"},
		"FileC is c file": {file: FileC, want: "c"},
		"FileD is d file": {file: FileD, want: "d"},
		"FileE is e file": {file: FileE, want: "e"},
		"FileF is f file": {file: FileF, want: "f"},
		"FileG is g file": {file: FileG, want: "g"},
		"FileH is h file": {file: FileH, want: "h"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.file.String() != tc.want {
				t.Fatalf("want %s, got %s", tc.want, tc.file)
			}
		})
	}
}

func TestSquareRanks(t *testing.T) {
	testCases := map[string]struct {
		rank Rank
		want string
	}{
		"Rank1 is first rank":   {rank: Rank1, want: "1"},
		"Rank2 is second rank":  {rank: Rank2, want: "2"},
		"Rank3 is third rank":   {rank: Rank3, want: "3"},
		"Rank4 is fourth rank":  {rank: Rank4, want: "4"},
		"Rank5 is fifth rank":   {rank: Rank5, want: "5"},
		"Rank6 is sixth rank":   {rank: Rank6, want: "6"},
		"Rank7 is seventh rank": {rank: Rank7, want: "7"},
		"Rank8 is eighth rank":  {rank: Rank8, want: "8"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.rank.String() != tc.want {
				t.Fatalf("want %s, got %s", tc.want, tc.rank)
			}
		})
	}
}

func TestSquareNewFile(t *testing.T) {
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
				t.Fatalf("error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}

			if !tc.wantErr && file != tc.want {
				t.Fatalf("want %s, got %s", tc.want, file)
			}
		})
	}
}

func TestSquareNewRank(t *testing.T) {
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
				t.Fatalf("error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}

			if !tc.wantErr && rank != tc.want {
				t.Fatalf("want %s, got %s", tc.want, rank)
			}
		})
	}
}
