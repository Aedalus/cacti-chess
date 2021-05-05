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
			k := &movekey{uint64(i)}
			assert.Equal(t, i, k.getFrom())

			// make sure other keys don't mess it up
			k.val = k.val | 1<<7
			assert.Equal(t, i, k.getFrom())
		}
	})

	t.Run("from-set", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < BOARD_SQ_NUMBER; i++ {
			k := &movekey{}
			k.setFrom(i)
			assert.Equal(t, uint64(i), k.val)
			assert.Equal(t, i, k.getFrom())
			assert.Equal(t, uint64(0), k.val & ^moveKeyFromBitmask)

			j := &movekey{maxUInt64}
			j.setFrom(i)
			assert.Equal(t, i, j.getFrom())
			assert.True(t, j.val&^moveKeyFromBitmask == ^moveKeyFromBitmask)
		}
	})

	t.Run("to-get", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < 65; i++ {
			val := uint64(i << 7)
			k := &movekey{val}
			//fmt.Println(val, val>>7)
			assert.Equal(t, i, k.getTo())

			// make sure other keys don't mess it up
			k.val = k.val | 1
			assert.Equal(t, i, k.getTo())
		}
	})

	t.Run("to-set", func(t *testing.T) {
		// test that its right if no flag has been set
		for i := 0; i < BOARD_SQ_NUMBER; i++ {
			k := &movekey{}
			k.setTo(i)
			assert.Equal(t, uint64(i)<<7, k.val)
			assert.Equal(t, i, k.getTo())
			assert.Equal(t, uint64(0), k.val & ^moveKeyToBitmask)

			j := &movekey{maxUInt64}
			j.setTo(i)
			assert.Equal(t, i, j.getTo())
			assert.True(t, j.val&^moveKeyToBitmask == ^moveKeyToBitmask)
		}

		// test that it preserves flags
		k := &movekey{}
		k.val = k.val | uint64(1)
		assert.Equal(t, uint64(1), k.val)

		k.setTo(2)
		assert.Equal(t, uint64(0x101), k.val)

		k.setTo(4)
		assert.Equal(t, uint64(0x201), k.val)
	})

	t.Run("enPas", func(t *testing.T) {
		k := &movekey{}

		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isEnPas())

		k.setEnPas()
		assert.Equal(t, uint64(1<<18), k.val)
		assert.True(t, k.isEnPas())
		assert.Equal(t, uint64(0), k.val & ^moveKeyEnPasBitmask)

		k.clearEnPas()
		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isEnPas())

		j := &movekey{maxUInt64}
		assert.True(t, j.isEnPas())

		j.clearEnPas()
		assert.Equal(t, maxUInt64 & ^uint64(1<<18), j.val)
	})

	t.Run("pawnStart", func(t *testing.T) {
		k := &movekey{}

		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isPawnStart())

		k.setPawnStart()
		assert.Equal(t, uint64(1<<19), k.val)
		assert.True(t, k.isPawnStart())
		assert.Equal(t, uint64(0), k.val & ^moveKeyPawnStartBitmask)

		k.clearPawnStart()
		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isPawnStart())

		j := &movekey{maxUInt64}
		assert.True(t, j.isPawnStart())
		j.clearPawnStart()
		assert.Equal(t, maxUInt64 & ^uint64(1<<19), j.val)
	})

	t.Run("castle", func(t *testing.T) {
		k := &movekey{}

		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isCastle())

		k.setCastle()
		assert.Equal(t, uint64(1<<24), k.val)
		assert.True(t, k.isCastle())

		k.clearCastle()
		assert.Equal(t, uint64(0), k.val)
		assert.False(t, k.isCastle())

		j := &movekey{maxUInt64}
		assert.True(t, j.isCastle())
		j.clearCastle()
		assert.Equal(t, maxUInt64 & ^uint64(1<<24), j.val)
	})

	t.Run("capture", func(t *testing.T) {
		k := &movekey{}
		j := &movekey{maxUInt64}

		for i := 0; i < PIECE_COUNT; i++ {
			k.setCaptured(piece(i))
			assert.Equal(t, piece(i), k.getCaptured())
			k.setCaptured(EMPTY)
			assert.Equal(t, EMPTY, k.getCaptured())
			assert.Equal(t, uint64(0), k.val & ^moveKeyCapturedBitmask)

			j.setCaptured(piece(i))
			assert.Equal(t, piece(i), j.getCaptured())
			j.setCaptured(EMPTY)
			assert.Equal(t, EMPTY, j.getCaptured())
			assert.True(t, j.val & ^moveKeyCapturedBitmask == ^moveKeyCapturedBitmask)
		}
	})

	t.Run("promote", func(t *testing.T) {
		k := &movekey{}
		j := &movekey{maxUInt64}

		for i := 0; i < PIECE_COUNT; i++ {
			k.setPromoted(piece(i))
			assert.Equal(t, piece(i), k.getPromoted())
			k.setPromoted(EMPTY)
			assert.Equal(t, EMPTY, k.getPromoted())
			assert.Equal(t, uint64(0), k.val & ^moveKeyPromotedPieceBitmask)

			j.setPromoted(piece(i))
			assert.Equal(t, piece(i), j.getPromoted())
			j.setPromoted(EMPTY)
			assert.Equal(t, EMPTY, j.getPromoted())
			assert.True(t, j.val & ^moveKeyPromotedPieceBitmask == ^moveKeyPromotedPieceBitmask)
		}
	})

	t.Run("enPas doesn't override capture", func(t *testing.T) {
		k := &movekey{}
		k.setEnPas()
		cap := k.getCaptured()
		assert.Equal(t, EMPTY, cap)
	})
}

func Test_printMove(t *testing.T) {
	m := &movekey{}
	m.setFrom(A2)
	m.setTo(A4)

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

		move, err := p.parseMove("a1a8")
		require.Nil(t, err)
		assert.Equal(t, move, &movekey{0})
	})

	t.Run("it will find basic moves", func(t *testing.T) {
		p, err := FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.Nil(t, err)

		move, err := p.parseMove("a2a4")
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
			move, err := p.parseMove(str)
			assert.Nil(t, err)
			assert.Equal(t, H7, move.getFrom())
			assert.Equal(t, H8, move.getTo())
			assert.Equal(t, tc.promPce, move.getPromoted())
		}
		move, err := p.parseMove("a2a4")
		require.Nil(t, err)
		assert.NotNil(t, move)

		assert.Equal(t, A2, move.getFrom())
		assert.Equal(t, A4, move.getTo())
	})
}
