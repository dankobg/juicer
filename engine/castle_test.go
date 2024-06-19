package engine

import (
	"testing"
)

func TestCastleRights(t *testing.T) {
	testCases := map[string]struct {
		cr  CastleRights
		fen string
	}{
		"Castle none":              {cr: CastleRightsNone, fen: "-"},
		"Castle wk":                {cr: WhiteKingSideCastle, fen: "K"},
		"Castle wq":                {cr: WhiteQueenSideCastle, fen: "Q"},
		"Castle bk":                {cr: BlackKingSideCastle, fen: "k"},
		"Castle bq":                {cr: BlackQueenSideCastle, fen: "q"},
		"Castle wk + wq":           {cr: WhiteKingSideCastle | WhiteQueenSideCastle, fen: "KQ"},
		"Castle wq + bk":           {cr: WhiteQueenSideCastle | BlackKingSideCastle, fen: "Qk"},
		"Castle wk + bk":           {cr: WhiteKingSideCastle | BlackKingSideCastle, fen: "Kk"},
		"Castle wk + bq":           {cr: WhiteKingSideCastle | BlackQueenSideCastle, fen: "Kq"},
		"Castle wq + bq":           {cr: WhiteQueenSideCastle | BlackQueenSideCastle, fen: "Qq"},
		"Castle bk + bq":           {cr: BlackKingSideCastle | BlackQueenSideCastle, fen: "kq"},
		"Castle wk + wq + bk":      {cr: WhiteKingSideCastle | WhiteQueenSideCastle | BlackKingSideCastle, fen: "KQk"},
		"Castle wk + wq + bq":      {cr: WhiteKingSideCastle | WhiteQueenSideCastle | BlackQueenSideCastle, fen: "KQq"},
		"Castle wk + bk + bq":      {cr: WhiteKingSideCastle | BlackKingSideCastle | BlackQueenSideCastle, fen: "Kkq"},
		"Castle wq + bk + bq":      {cr: WhiteQueenSideCastle | BlackKingSideCastle | BlackQueenSideCastle, fen: "Qkq"},
		"Castle wk + wq + bk + bq": {cr: WhiteKingSideCastle | WhiteQueenSideCastle | BlackKingSideCastle | BlackQueenSideCastle, fen: "KQkq"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.cr.ToFEN() != tc.fen {
				t.Fatalf("invalid castle rights fen: want %s, got %s", tc.fen, tc.cr.ToFEN())
			}
		})
	}
}

func TestNewCastleRightsFromFen(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    CastleRights
		wantErr bool
	}{
		"Creates - cr":                {in: "-", wantErr: false},
		"Creates K cr":                {in: "K", wantErr: false},
		"Creates Q cr":                {in: "Q", wantErr: false},
		"Creates k cr":                {in: "k", wantErr: false},
		"Creates q cr":                {in: "q", wantErr: false},
		"Creates KQ cr":               {in: "KQ", wantErr: false},
		"Creates Qk cr":               {in: "Qk", wantErr: false},
		"Creates Kk cr":               {in: "Kk", wantErr: false},
		"Creates Kq cr":               {in: "Kq", wantErr: false},
		"Creates Qq cr":               {in: "Qq", wantErr: false},
		"Creates kq cr":               {in: "kq", wantErr: false},
		"Creates KQk cr":              {in: "KQk", wantErr: false},
		"Creates KQq cr":              {in: "KQq", wantErr: false},
		"Creates Kkq cr":              {in: "Kkq", wantErr: false},
		"Creates Qkq cr":              {in: "Qkq", wantErr: false},
		"Creates KQkq cr":             {in: "KQkq", wantErr: false},
		"fails with empty":            {in: "", wantErr: true},
		"fails with duplicate -":      {in: "--", wantErr: true},
		"fails with duplicate K":      {in: "KKQkq", wantErr: true},
		"fails with duplicate Q":      {in: "KQQkq", wantErr: true},
		"fails with duplicate k":      {in: "KQkkq", wantErr: true},
		"fails with duplicate q":      {in: "KQkqq", wantErr: true},
		"fails with wrong order QK":   {in: "QK", wantErr: true},
		"fails with wrong order kQ":   {in: "kQ", wantErr: true},
		"fails with wrong order kK":   {in: "kK", wantErr: true},
		"fails with wrong order qK":   {in: "qK", wantErr: true},
		"fails with wrong order qQ":   {in: "qQ", wantErr: true},
		"fails with wrong order qk":   {in: "qk", wantErr: true},
		"fails with wrong order QKk":  {in: "QKk", wantErr: true},
		"fails with wrong order QKq":  {in: "QKq", wantErr: true},
		"fails with wrong order kqK":  {in: "kqK", wantErr: true},
		"fails with wrong order kqQ":  {in: "kqQ", wantErr: true},
		"fails with wrong order kqKQ": {in: "kqKQ", wantErr: true},
		"fails with wrong order kqQK": {in: "kqQK", wantErr: true},
		"fails with wrong order kKq":  {in: "kKq", wantErr: true},
		"fails with wrong order kQq":  {in: "kQq", wantErr: true},
		"fails with wrong order kKQq": {in: "kKQq", wantErr: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cr, err := NewCastleRightsFromFen(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid castle rights, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && cr != tc.want {
				t.Fatalf("invalid castle rights, want %s, got %s", tc.want, cr)
			}
		})
	}
}
