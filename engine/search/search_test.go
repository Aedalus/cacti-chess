package search

import (
	"cacti-chess/engine/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSearchInfo_SearchPosition(t *testing.T) {
	//t.Run("starting alphabeta - depth 0", func(t *testing.T) {
	//	p, err := position.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	//	require.Nil(t, err)
	//
	//	s := New()
	//	val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 1, false)
	//	assert.Equal(t, float64(0), val)
	//	line := s.pvTable.GetPVLine(p)
	//	names := []string{}
	//	for _, mv := range line {
	//		names = append(names, mv.ShortString())
	//	}
	//	assert.Equal(t, []string{"d2d4"}, names)
	//})

	//t.Run("mate in 1 - depth 0", func(t *testing.T) {
	//	p, err := position.FromFen("k7/7Q/K7/8/8/8/8/8 w - - 0 1")
	//	require.Nil(t, err)
	//
	//	s := New()
	//	val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 2, false)
	//	require.Equal(t, float64(29000), val)
	//	line := s.pvTable.GetPVLine(p)
	//	names := []string{}
	//	for _, mv := range line {
	//		names = append(names, mv.ShortString())
	//	}
	//	assert.Equal(t, []string{"f6a6", "g8g7", "a6a8"}, names)
	//})

	t.Run("mate in 1 - depth 0", func(t *testing.T) {
		p, err := position.FromFen("k7/7Q/K7/8/8/8/8/8 b - - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 3, false)
		require.Equal(t, float64(-29000), val)
		line := s.pvTable.GetPVLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"a8b8", "h7b7"}, names)
	})

	t.Run("mate in 1 - depth 0", func(t *testing.T) {
		p, err := position.FromFen("r5rk/5p1p/5R2/4B3/8/8/7P/7K w - - 0 1")
		require.Nil(t, err)

		s := New()
		val := s.AlphaBeta(p, math.Inf(-1), math.Inf(1), 6, false)
		assert.Equal(t, float64(29000), val)
		line := s.pvTable.GetPVLine(p)
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
		line := s.pvTable.GetPVLine(p)
		names := []string{}
		for _, mv := range line {
			names = append(names, mv.ShortString())
		}
		assert.Equal(t, []string{"f3f7", "e8f7", "g4h5"}, names)
	})
}
