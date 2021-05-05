package position

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMappingBoards(t *testing.T) {
	// spot check some of the mapping
	t.Run("the 120 engine maps correctly", func(t *testing.T) {
		// borders should all be 01
		assert.Equal(t, -1, sq120to64[0])
		assert.Equal(t, -1, sq120to64[9])
		assert.Equal(t, -1, sq120to64[40])
		assert.Equal(t, -1, sq120to64[49])
		assert.Equal(t, -1, sq120to64[110])
		assert.Equal(t, -1, sq120to64[119])

		// first rank
		assert.Equal(t, 0, sq120to64[21])
		assert.Equal(t, 1, sq120to64[22])
		assert.Equal(t, 2, sq120to64[23])

		// middle rank
		assert.Equal(t, 8, sq120to64[31])
		assert.Equal(t, 15, sq120to64[38])

		// top rank
		assert.Equal(t, 48, sq120to64[81])
		assert.Equal(t, 56, sq120to64[91])
		assert.Equal(t, 63, sq120to64[98])
	})

	t.Run("the 64 engine maps correctly", func(t *testing.T) {
		// first rank
		assert.Equal(t, 21, sq64to120[0])
		assert.Equal(t, 28, sq64to120[7])

		// middle rank
		assert.Equal(t, 61, sq64to120[32])
		assert.Equal(t, 68, sq64to120[39])

		// last rank
		assert.Equal(t, 91, sq64to120[56])
		assert.Equal(t, 98, sq64to120[63])
	})
}

func TestBitboard64(t *testing.T) {
	t.Run("it can set/clear/toggle/has", func(t *testing.T) {
		for i := 0; i < 64; i++ {
			b := &bitboard64{val: 0}
			assert.Equal(t, false, b.has(i))
		}

		b := &bitboard64{}
		b.set(0)
		b.set(2)
		b.set(4)

		assert.True(t, b.has(0))
		assert.False(t, b.has(1))
		assert.True(t, b.has(2))
		assert.False(t, b.has(3))
		assert.True(t, b.has(4))

		b.clear(0)
		b.clear(2)
		b.clear(4)

		assert.False(t, b.has(0))
		assert.False(t, b.has(1))
		assert.False(t, b.has(2))
		assert.False(t, b.has(3))
		assert.False(t, b.has(4))

	})

	t.Run("test has", func(t *testing.T) {
		b := &bitboard64{}
		b.set(0)
		b.set(1)
		assert.True(t, b.has(0))
		assert.True(t, b.has(1))
		assert.False(t, b.has(2))
	})
}
