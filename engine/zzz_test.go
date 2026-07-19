package engine

import (
	"fmt"
	"testing"
)

func TestJuicer(t *testing.T) {
	InitPrecalculatedTables()

	c, _ := NewChess(FENStartingPosition)

	c.MakeMoveUCI("e2e4")
	c.MakeMoveUCI("e7e5")
	c.MakeMoveUCI("g1f3")
	c.MakeMoveUCI("b8c6")
	c.MakeMoveUCI("f1b5")
	c.MakeMoveUCI("a7a6")
	c.MakeMoveUCI("b5a4")
	c.MakeMoveUCI("b7b5")
	c.MakeMoveUCI("a4b3")
	c.MakeMoveUCI("d7d5")
	c.MakeMoveUCI("e4d5")
	c.MakeMoveUCI("h7h6")
	c.MakeMoveUCI("d5c6")
	c.MakeMoveUCI("d8d7")
	c.MakeMoveUCI("c6d7")
	c.MakeMoveUCI("e8e7")
	c.MakeMoveUCI("d7c8n")
	c.MakeMoveUCI("e7d8")
	c.MakeMoveUCI("d1e2")
	c.MakeMoveUCI("d8c8")

	fmt.Println(c.Position.Turn)
	fmt.Println(c.Position.PrintBoard())
}
