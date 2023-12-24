package juicer

import (
	"testing"
)

func TestBitboardClear(t *testing.T) {
	bb := bitboardLightSquares
	bb.clear()

	if bb != bitboardEmpty {
		t.Fatalf("faild to clear bitboard")
	}
}

func TestBitboardFill(t *testing.T) {
	bb := bitboardEmpty
	bb.fill()

	if bb != bitboardFull {
		t.Fatalf("faild to fill bitboard")
	}
}

func TestBitboardSetBit(t *testing.T) {
	bb := bitboardEmpty
	bb.setBit(E4)

	if bb != 0x10000000 {
		t.Fatalf("faild to set bit")
	}
}

func TestBitboardToggleBit(t *testing.T) {
	bb := bitboardEmpty

	bb.toggleBit(E4)
	if bb.getBit(E4) != 1 {
		t.Fatalf("faild to toggle bit")
	}

	bb.toggleBit(E4)
	if bb.getBit(E4) != 0 {
		t.Fatalf("faild to toggle bit")
	}
}

func TestBitboardClearBit(t *testing.T) {
	bb := bitboard(0x10000000)
	bb.clearBit(E4)

	if bb != bitboardEmpty {
		t.Fatalf("faild to clear bit")
	}
}

func TestBitboardGetBit(t *testing.T) {
	bb := bitboard(0x10000000)

	if bb.getBit(E4) != 1 {
		t.Fatalf("faild to get bit")
	}
}

func TestBitboardBitIsSet(t *testing.T) {
	bb := bitboard(0x10000000)

	if !bb.bitIsSet(E4) {
		t.Fatalf("failed bitIsSet check")
	}
}

func TestBitboardMS1B(t *testing.T) {
	bb := bitboard(0x2002012408802004)

	if bb.MS1B() != 61 {
		t.Fatalf("failed MS1B")
	}
}

func TestBitboardLS1B(t *testing.T) {
	bb := bitboard(0x2002012408802004)

	if bb.LS1B() != 2 {
		t.Fatalf("failed LS1B")
	}
}

func TestBitboardPopMS1B(t *testing.T) {
	bb := bitboard(0x2002012408802004)
	bb.PopMS1B()

	if bb != bitboard(0x2012408802004) {
		t.Fatalf("failed PopMS1B")
	}
}

func TestBitboardPopLS1B(t *testing.T) {
	bb := bitboard(0x2002012408802004)
	bb.PopLS1B()

	if bb != bitboard(0x2002012408802000) {
		t.Fatalf("failed PopLS1B")
	}
}

func TestBitboardIsEmptyAndFull(t *testing.T) {
	empty, full := bitboardEmpty, bitboardFull

	if !empty.isEmpty() || empty.isFull() || full.isEmpty() || !full.isFull() {
		t.Fatalf("failed is empty/full")
	}
}

func TestBitboardFlipVertical(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.flipVertical() != bitboard(0x22120A0E1222221E) {
		t.Fatalf("failed flip vertical")
	}
}

func TestBitboardFlipHorizontal(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.flipHorizontal() != bitboard(0x7844444870504844) {
		t.Fatalf("failed flip horizontal")
	}
}

func TestBitboardFlipDiagonalA1H8(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.flipDiagonalA1H8() != bitboard(0x61928C88FF00) {
		t.Fatalf("failed flip diagonal a1h8")
	}
}

func TestBitboardFlipDiagonalA8H1(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.flipDiagonalA8H1() != bitboard(0xFF113149860000) {
		t.Fatalf("failed flip diagonal a8h1")
	}
}

func TestBitboardRotate180(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.rotate180() != bitboard(0x4448507048444478) {
		t.Fatalf("failed rotate 180")
	}
}

func TestBitboardRotate90clockwise(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.rotate90clockwise() != bitboard(0xFF888C92610000) {
		t.Fatalf("failed rotate 90 clockwise")
	}
}

func TestBitboardRotate90counterClockwise(t *testing.T) {
	bb := bitboard(0x1E2222120E0A1222)

	if bb.rotate90counterClockwise() != bitboard(0x86493111FF00) {
		t.Fatalf("failed rotate 90 counter clockwise")
	}
}
