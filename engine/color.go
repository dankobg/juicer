package engine

// Color represents the color (used for pieces, squares, turns etc.)
type Color int8

const (
	ColorNone Color = iota - 1
	White
	Black
)

// String returns a FEN compatible color character
func (c Color) String() string {
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	}

	return ""
}

// Name is a color's display friendly name
func (c Color) Name() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	}

	return "No Color"
}

// IsWhite returns true if color is white
func (c Color) IsWhite() bool {
	return c == White
}

// IsBlack returns true if color is black
func (c Color) IsBlack() bool {
	return c == Black
}

// Opposite returns the opposite color
func (c Color) Opposite() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	}

	return ColorNone
}
