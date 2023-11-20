package juicer

import (
	"fmt"
	"testing"
)

func TestJuicer(t *testing.T) {
	InitAllAttackMasksTables()

	occ := bitboardEmpty
	occ.setBit(D4)
	occ.setBit(F5)

	fmt.Println(getQueenAttacks(E4, occ).draw(nil))
}
