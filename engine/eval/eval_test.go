package eval

import (
	"cacti-chess/engine/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCaseScore struct {
	name string
	fen  string
	want int
}

func (tc testCaseScore) assert(t *testing.T) {
	t.Helper()

	scr := PositionScorer{}
	p, err := position.FromFen(tc.fen)
	require.Nil(t, err)
	score := scr.Score(p)
	assert.Equal(t, tc.want, score)
}

func TestPositionScorer_Score(t *testing.T) {

	tcs := []testCaseScore{
		{
			"starting position",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			0,
		},
		{
			"black down a rook",
			"rnbqkbn1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			550,
		},
		{
			"wN move",
			"rnbqkbnr/pppppppp/8/8/8/2N5/PPPPPPPP/R1BQKBNR w KQkq - 0 1",
			20,
		},
	}

	for _, tc := range tcs {
		tc.assert(t)
	}
}
