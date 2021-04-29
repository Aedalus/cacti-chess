package engine

import (
	"github.com/stretchr/testify/assert"
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
