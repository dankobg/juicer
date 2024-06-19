package engine

import "testing"

func TestColor(t *testing.T) {
	testCases := map[string]struct {
		color            Color
		fen              string
		opposite         Color
		isWhite, isBlack bool
	}{
		"White is w": {color: White, fen: "w", opposite: Black, isWhite: true, isBlack: false},
		"Black is b": {color: Black, fen: "b", opposite: White, isWhite: false, isBlack: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.color.String() != tc.fen {
				t.Fatalf("invalid color: want %s, got %s", tc.fen, tc.color)
			}
			if tc.color.Opposite() != tc.opposite {
				t.Fatalf("invalid opposite color: want %d, got %d", tc.opposite, tc.color.Opposite())
			}
			if (tc.isWhite != tc.color.IsWhite()) || (tc.isBlack != tc.color.IsBlack()) {
				t.Fatalf("invalid variant: want %v, %v, got %v, %v", tc.isWhite, tc.isBlack, tc.color.IsWhite(), tc.color.IsBlack())
			}
		})
	}
}

func TestNewColorFromFenStr(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    Color
		wantErr bool
	}{
		"creates white color":      {in: "w", wantErr: false, want: White},
		"creates black color":      {in: "b", wantErr: false, want: Black},
		"fails with invalid empty": {in: "", wantErr: true},
		"fails with invalid white": {in: "white", wantErr: true},
		"fails with invalid black": {in: "black", wantErr: true},
		"fails with invalid len w": {in: "ww", wantErr: true},
		"fails with invalid len b": {in: "bb", wantErr: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			color, err := NewColorFromFenStr(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid color, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && color != tc.want {
				t.Fatalf("invalid color: want %s, got %s", tc.want, color)
			}
		})
	}
}
