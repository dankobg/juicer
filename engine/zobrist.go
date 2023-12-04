package juicer

import (
	"math/rand"
	"sync"
)

var (
	onceZobrist        sync.Once
	zobristInitialized bool
	defaultZobrist     zobrist
)

type zobrist struct {
	seed            uint64
	occupanciesKeys [2][6][64]uint64
	castleKeys      map[CastleRights]uint64
	turnKey         uint64
	enpKeys         [64]uint64
}

func initZobrist() {
	if !zobristInitialized {
		onceZobrist.Do(func() {
			defaultZobrist = zobrist{
				seed: rand.Uint64(),
			}

			for sq := A1; sq <= H8; sq++ {
				for _, color := range colors {
					for _, pk := range pieceKinds {
						defaultZobrist.occupanciesKeys[color][pk][sq] = rand.Uint64()
					}
				}

				defaultZobrist.enpKeys[sq] = rand.Uint64()
			}

			defaultZobrist.castleKeys = make(map[CastleRights]uint64)

			defaultZobrist.castleKeys[WhiteKingSideCastle] = rand.Uint64()
			defaultZobrist.castleKeys[WhiteQueenSideCastle] = rand.Uint64()
			defaultZobrist.castleKeys[BlackKingSideCastle] = rand.Uint64()
			defaultZobrist.castleKeys[BlackQueenSideCastle] = rand.Uint64()

			defaultZobrist.turnKey = rand.Uint64()

			zobristInitialized = true
		})
	}
}
