package search

import (
	"cacti-chess/engine/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSearchInfo_SearchPosition(t *testing.T) {
	t.Run("starting alphabeta - depth 0", func(t *testing.T) {
		p, err := position.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 1, false)
		assert.Equal(t, float64(30), val)
		line := s.pvTable.GetBestLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"d2d4"}, names)
	})

	t.Run("checkmate test - 0", func(t *testing.T) {
		p, err := position.FromFen("k7/7Q/K7/8/8/8/8/8 w - - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 2, false)
		require.Equal(t, float64(29000), val)
		line := s.pvTable.GetBestLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"h7b7"}, names)
	})

	t.Run("checkmate test - 1", func(t *testing.T) {
		p, err := position.FromFen("k7/7Q/K7/8/8/8/8/8 b - - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 3, false)
		require.Equal(t, float64(-29000), val)
		line := s.pvTable.GetBestLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"a8b8", "h7b7"}, names)
	})

	t.Run("checkmate test - 2", func(t *testing.T) {
		p, err := position.FromFen("r5rk/5p1p/5R2/4B3/8/8/7P/7K w - - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 6, false)
		assert.Equal(t, float64(29000), val)
		line := s.pvTable.GetBestLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"f6a6", "f7f6", "e5f6", "g8g7", "a6a8"}, names)
	})

	t.Run("mate in 2 - depth 0", func(t *testing.T) {
		p, err := position.FromFen("2bqkbn1/2pppp2/np2N3/r3P1p1/p2N2B1/5Q2/PPPPKPP1/RNB2r2 w KQkq - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 4, false)
		assert.Equal(t, float64(29000), val)
		line := s.pvTable.GetBestLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"f3f7", "e8f7", "g4h5"}, names)
	})

	t.Run("crash-1", func(t *testing.T) {
		//p, err := position.FromFen("rnbqk2r/ppp2pQp/8/3pP3/4n3/2P5/PPP3PP/R1B1KBNR b KQkq - 0 7")
		p, err := position.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		moves := []string{
			"e2e4", "g8f6",
			"b1c3", "e7e5",
			"f2f4", "f8b4",
			"f4e5", "b4c3",
			"d2c3", "f6e4",
			"d1g4", "d7d5",
			"g4g7",
		}

		for _, mv := range moves {
			mv, err := p.ParseMove(mv)
			if err != nil {
				panic(err)
			}
			p.MakeMove(mv)

			s := New()
			s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 3, false)
			line := s.pvTable.GetBestLine(p)
			assert.Len(t, line, 3)

		}
	})
}
