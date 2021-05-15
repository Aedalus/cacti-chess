package position

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPosition_MakeMove(t *testing.T) {
	t.Run("crash-1", func(t *testing.T) {
		//p, err := position.FromFen("rnbqk2r/ppp2pQp/8/3pP3/4n3/2P5/PPP3PP/R1B1KBNR b KQkq - 0 7")
		p, err := FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
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

		for i, mv := range moves {
			mv, err := p.ParseMove(mv)
			if err != nil {
				panic(err)
			}
			p.MakeMove(mv)
			require.Equal(t, p.hisPly, i+1)
			require.Len(t, p.history, i+1)
		}

		for i := 0; i < len(moves); i++ {
			p.UndoMove()
			require.Equal(t, p.hisPly, len(moves)-i-1)
			require.Len(t, p.history, len(moves)-i-1)
		}
	})
}
