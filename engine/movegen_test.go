package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertExpectedMovelistLen(t *testing.T, expected int, fenStr string) {
	t.Helper()
	p, err := ParseFen(fenStr)
	assert.Nil(t, err)

	list := &movelist{}
	generateAllMoves(p, list)
	fmt.Println(list)
	assert.Equal(t, expected, list.count)
}

func Test_MoveGen_MovelistCount(t *testing.T) {

	cases := []struct {
		name     string
		fen      string
		expected int
	}{
		{
			"white pawns",
			"rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1",
			32,
		}, {

			"black pawns",
			"rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/PNBQKB1R b KQkq e3 0 1",
			32,
		}, {
			"knights + kings",
			"5k2/1n6/4n3/6N1/8/3N4/8/5K2 b - - 0 1",
			16,
		},
	}

	for _, c := range cases {
		assertExpectedMovelistLen(t, c.expected, c.fen)
	}
}
