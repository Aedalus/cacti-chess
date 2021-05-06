package position

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_FenParsePiecesStr(t *testing.T) {
	t.Run("it can parse starting position", func(t *testing.T) {
		fenStr := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
		got, err := parsePiecesStr(fenStr)

		assert.Nil(t, err)

		want := [64]Piece{
			PwR, PwN, PwB, PwQ, PwK, PwB, PwN, PwR,
			PwP, PwP, PwP, PwP, PwP, PwP, PwP, PwP,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			PbP, PbP, PbP, PbP, PbP, PbP, PbP, PbP,
			PbR, PbN, PbB, PbQ, PbK, PbB, PbN, PbR,
		}
		assert.Equal(t, want, got)
	})
	t.Run("it can parse gap numbers", func(t *testing.T) {
		fenStr := "1ppppppp/2pppppp/3ppppp/4pppp/5ppp/6pp/7p/8"
		got, err := parsePiecesStr(fenStr)

		assert.Nil(t, err)

		want := [64]Piece{
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, PbP,
			00, 00, 00, 00, 00, 00, PbP, PbP,
			00, 00, 00, 00, 00, PbP, PbP, PbP,
			00, 00, 00, 00, PbP, PbP, PbP, PbP,
			00, 00, 00, PbP, PbP, PbP, PbP, PbP,
			00, 00, PbP, PbP, PbP, PbP, PbP, PbP,
			00, PbP, PbP, PbP, PbP, PbP, PbP, PbP,
		}
		assert.Equal(t, want, got)
	})
}

func Test_FenParseCastlePerms(t *testing.T) {
	t.Run("it can test castle perms", func(t *testing.T) {
		tests := []struct {
			permStr  string
			expected int
		}{
			{"K", CASTLE_PERMS_WK},
			{"Q", CASTLE_PERMS_WQ},
			{"k", CASTLE_PERMS_BK},
			{"q", CASTLE_PERMS_BQ},
			{"Kk", CASTLE_PERMS_BK + CASTLE_PERMS_WK},
			{"Qq", CASTLE_PERMS_BQ + CASTLE_PERMS_WQ},
			{"KkQq", CASTLE_PERMS_ALL},
			{"-", CASTLE_PERMS_NONE},
		}

		for _, tc := range tests {
			parsed := parseCastlePerms(tc.permStr)
			assert.Equal(t, tc.expected, parsed.val)
		}
	})
}

func Test_FromFen(t *testing.T) {
	t.Run("it can parse a starting position", func(t *testing.T) {
		fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		state, err := FromFen(fen)

		require.Nil(t, err)
		assert.Equal(t, WHITE, state.side)
		assert.Equal(t, &castlePerm{val: CASTLE_PERMS_ALL}, state.castlePerm)
		assert.Equal(t, -1, state.enPas)
		assert.Equal(t, 0, state.fiftyMove)
		assert.Equal(t, 0, state.hisPly)

		want := &board120{
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, PwR, PwN, PwB, PwQ, PwK, PwB, PwN, PwR, -1,
			-1, PwP, PwP, PwP, PwP, PwP, PwP, PwP, PwP, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, PbP, PbP, PbP, PbP, PbP, PbP, PbP, PbP, -1,
			-1, PbR, PbN, PbB, PbQ, PbK, PbB, PbN, PbR, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		}

		assert.Equal(t, want, state.pieces)
	})
}
