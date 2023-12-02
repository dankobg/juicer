package juicer

import (
	"math/rand"
)

type zobrist struct {
	hash            uint64
	occupanciesKeys [2][6][64]uint64
	castleKeys      map[CastleRights]uint64
	turnKey         uint64
	enpKeys         [64]uint64
}

func newZobrist() *zobrist {
	z := zobrist{
		hash: rand.Uint64(),
	}

	for sq := A1; sq <= H8; sq++ {
		for _, color := range colors {
			for _, pk := range pieceKinds {
				z.occupanciesKeys[color][pk][sq] = rand.Uint64()
			}
		}

		z.enpKeys[sq] = rand.Uint64()
	}

	z.castleKeys = make(map[CastleRights]uint64)

	z.castleKeys[WhiteKingSideCastle] = rand.Uint64()
	z.castleKeys[WhiteQueenSideCastle] = rand.Uint64()
	z.castleKeys[BlackKingSideCastle] = rand.Uint64()
	z.castleKeys[BlackQueenSideCastle] = rand.Uint64()

	z.turnKey = rand.Uint64()

	return &z
}
