package gameserver

import "testing"

func TestClientExample(t *testing.T) {
	sum := 15

	a, b := 5, 10

	if a+b != sum {
		t.Errorf("expected sum: %v but got: %+v\n", sum, a+b)
	}
}
