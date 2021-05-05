package position

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCaseMoveGen struct {
	name     string
	fen      string
	expected int
	print    bool
}

func (tc testCaseMoveGen) assert(t *testing.T) {
	t.Helper()
	p, err := FromFen(tc.fen)
	require.Nil(t, err)
	require.NotNil(t, p)

	if p == nil {
		panic("p should not be nil")
	}

	mlist := p.generateAllMoves()
	assert.Equal(t, tc.expected, mlist.count)
}

func Test_MoveGen_MovelistCount(t *testing.T) {

	cases := []testCaseMoveGen{
		{
			"white pawns",
			"rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1",
			42,
			false,
		},
		{

			"black pawns",
			"rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/PNBQKB1R b KQkq e3 0 1",
			42,
			false,
		},
		{
			"knights + kings",
			"5k2/1n6/4n3/6N1/8/3N4/8/5K2 b - - 0 1",
			16,
			false,
		},
		{
			"rooks",
			"6k1/8/5r2/8/1nR5/5N2/8/6K1 b - - 0 1",
			23,
			false,
		},
		{
			"queens",
			"6k1/8/4nq2/8/1nQ5/5N2/1N6/6K1 b - - 0 1",
			36,
			false,
		},
		{
			"bishops",
			"6k1/1b6/4n3/8/1n4B1/1B3N2/1N6/2b3K1 b - - 0 1",
			32,
			false,
		},
		{
			"castle",
			"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			26,
			false,
		},
		{
			"castle2",
			"3rk2r/8/8/8/8/8/6p1/R3K2R w KQk - 0 1",
			24,
			false,
		},
	}

	for _, c := range cases {
		c.assert(t)
	}
}
