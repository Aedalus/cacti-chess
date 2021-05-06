package position

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const maxUInt64 uint64 = 18446744073709551615

func Test_movekey(t *testing.T) {
	t.Run("from-get", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < BOARD_SQ_NUMBER; i++ {
			k := movekey(i)
			assert.Equal(t, i, k.getFrom())

			// make sure other keys don't mess it up
			k = k | 1<<7
			assert.Equal(t, i, k.getFrom())
		}
	})

	t.Run("from-set", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < BOARD_SQ_NUMBER; i++ {
			k := movekey(0).setFrom(i)
			assert.Equal(t, movekey(i), k)
			assert.Equal(t, i, k.getFrom())
			assert.Equal(t, uint64(0), uint64(k) & ^moveKeyFromBitmask)

			j := movekey(maxUInt64).setFrom(i)
			assert.Equal(t, i, j.getFrom())
			assert.True(t, uint64(j)&^moveKeyFromBitmask == ^moveKeyFromBitmask)
		}
	})

	t.Run("to-get", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < 65; i++ {
			val := uint64(i << 7)
			k := movekey(val)
			assert.Equal(t, i, k.getTo())

			// make sure other keys don't mess it up
			k = movekey(uint64(k) | 1)
			assert.Equal(t, i, k.getTo())
		}
	})

	t.Run("to-set", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < BOARD_SQ_NUMBER; i++ {
			k := movekey(0).setTo(i)
			assert.Equal(t, movekey(i)<<7, k)
			assert.Equal(t, i, k.getTo())
			assert.Equal(t, uint64(0), uint64(k) & ^moveKeyToBitmask)

			j := movekey(maxUInt64).setTo(i)
			assert.Equal(t, i, j.getTo())
			assert.True(t, uint64(j)&^moveKeyToBitmask == ^moveKeyToBitmask)
		}

		// test that it preserves flags
		k := movekey(0)
		k = movekey(uint64(k) | uint64(1))
		assert.Equal(t, movekey(1), k)

		k = k.setTo(2)
		assert.Equal(t, movekey(0x101), k)

		k = k.setTo(4)
		assert.Equal(t, movekey(0x201), k)
	})

	t.Run("enPas", func(t *testing.T) {
		k := movekey(0)

		assert.Equal(t, movekey(0), k)
		assert.False(t, k.isEnPas())

		k = k.setEnPas()
		assert.Equal(t, movekey(1<<18), k)
		assert.True(t, k.isEnPas())
		assert.Equal(t, uint64(0), uint64(k) & ^moveKeyEnPasBitmask)

		k = k.clearEnPas()
		assert.Equal(t, movekey(0), k)
		assert.False(t, k.isEnPas())

		j := movekey(maxUInt64)
		assert.True(t, j.isEnPas())

		j = j.clearEnPas()
		assert.Equal(t, movekey(maxUInt64 & ^uint64(1<<18)), j)
	})

	t.Run("pawnStart", func(t *testing.T) {
		k := movekey(0)

		assert.Equal(t, uint64(0), uint64(k))
		assert.False(t, k.isPawnStart())

		k = k.setPawnStart()
		assert.Equal(t, movekey(1<<19), k)
		assert.True(t, k.isPawnStart())
		assert.Equal(t, uint64(0), uint64(k) & ^moveKeyPawnStartBitmask)

		k = k.clearPawnStart()
		assert.Equal(t, movekey(0), k)
		assert.False(t, k.isPawnStart())

		j := movekey(maxUInt64)
		assert.True(t, j.isPawnStart())
		j = j.clearPawnStart()
		assert.Equal(t, movekey(maxUInt64 & ^uint64(1<<19)), j)
	})

	t.Run("castle", func(t *testing.T) {
		k := movekey(0)

		assert.Equal(t, uint64(0), uint64(k))
		assert.False(t, k.isCastle())

		k = k.setCastle()
		assert.Equal(t, movekey(1<<24), k)
		assert.True(t, k.isCastle())

		k = k.clearCastle()
		assert.Equal(t, movekey(0), k)
		assert.False(t, k.isCastle())

		j := movekey(maxUInt64)
		assert.True(t, j.isCastle())
		j = j.clearCastle()
		assert.Equal(t, movekey(maxUInt64 & ^uint64(1<<24)), j)
	})

	t.Run("capture", func(t *testing.T) {
		k := movekey(0)
		j := movekey(maxUInt64)

		for i := 0; i < PIECE_COUNT; i++ {
			k = k.setCaptured(piece(i))
			assert.Equal(t, piece(i), k.getCaptured())
			k = k.setCaptured(EMPTY)
			assert.Equal(t, EMPTY, k.getCaptured())
			assert.Equal(t, uint64(0), uint64(k) & ^moveKeyCapturedBitmask)

			j = j.setCaptured(piece(i))
			assert.Equal(t, piece(i), j.getCaptured())
			j = j.setCaptured(EMPTY)
			assert.Equal(t, EMPTY, j.getCaptured())
			assert.True(t, uint64(j) & ^moveKeyCapturedBitmask == ^moveKeyCapturedBitmask)
		}
	})

	t.Run("promote", func(t *testing.T) {
		k := movekey(0)
		j := movekey(maxUInt64)

		for i := 0; i < PIECE_COUNT; i++ {
			k = k.setPromoted(piece(i))
			assert.Equal(t, piece(i), k.getPromoted())
			k = k.setPromoted(EMPTY)
			assert.Equal(t, EMPTY, k.getPromoted())
			assert.Equal(t, uint64(0), uint64(k) & ^moveKeyPromotedPieceBitmask)

			j = j.setPromoted(piece(i))
			assert.Equal(t, piece(i), j.getPromoted())
			j = j.setPromoted(EMPTY)
			assert.Equal(t, EMPTY, j.getPromoted())
			assert.True(t, uint64(j) & ^moveKeyPromotedPieceBitmask == ^moveKeyPromotedPieceBitmask)
		}
	})

	t.Run("enPas doesn't override capture", func(t *testing.T) {
		k := movekey(0).setEnPas()
		cap := k.getCaptured()
		assert.Equal(t, EMPTY, cap)
	})
}

func Test_printMove(t *testing.T) {
	m := movekey(0).setFrom(A2).setTo(A4)

	want := "from: a2\n"
	want += "to: a4\n"
	want += "enPas: false\n"
	want += "castle: false\n"
	want += "pawnStart: false\n"
	want += "captured: EMPTY\n"
	want += "promoted: EMPTY\n"
	assert.Equal(t, want, m.String())
}

func Test_ParseMove(t *testing.T) {
	t.Run("will return no move for not found", func(t *testing.T) {
		p, err := FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		move, err := p.ParseMove("a1a8")
		require.Nil(t, err)
		assert.Equal(t, move, movekey(0))
	})

	t.Run("it will find basic moves", func(t *testing.T) {
		p, err := FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		move, err := p.ParseMove("a2a4")
		require.Nil(t, err)
		assert.NotNil(t, move)

		assert.Equal(t, A2, move.getFrom())
		assert.Equal(t, A4, move.getTo())
	})

	t.Run("it will find castle moves", func(t *testing.T) {
		p, err := FromFen("rnbqkbn1/pppppppP/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		for _, tc := range []struct {
			pce     string
			promPce piece
		}{
			{"n", wN},
			{"q", wQ},
			{"r", wR},
			{"b", wB},
		} {
			str := "h7h8" + tc.pce
			move, err := p.ParseMove(str)
			assert.Nil(t, err)
			assert.Equal(t, H7, move.getFrom())
			assert.Equal(t, H8, move.getTo())
			assert.Equal(t, tc.promPce, move.getPromoted())
		}
		move, err := p.ParseMove("a2a4")
		require.Nil(t, err)
		assert.NotNil(t, move)

		assert.Equal(t, A2, move.getFrom())
		assert.Equal(t, A4, move.getTo())
	})
}
