package engine

const (
	FENStartingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	FENEmptyPosition    = "8/8/8/8/8/8/8/8 w KQkq - 0 1"
)

type Chess struct {
	position *Position
}
