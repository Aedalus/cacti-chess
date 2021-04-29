package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertExpectedMovelistLen(t *testing.T, expected int, fenStr string) {
	t.Helper()
	p, err := ParseFen(fenStr)
	assert.Nil(t, err)

	list := &movelist{}
	generateAllMoves(p, list)
	//fmt.Println(list)
	assert.Equal(t, expected, list.count)
}

func Test_MoveGen_WhitePawns(t *testing.T) {
	fen := "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
	assertExpectedMovelistLen(t, 26, fen)
}
