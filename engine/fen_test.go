package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFenParsePiecesStr(t *testing.T) {
	t.Run("it can parse starting position", func(t *testing.T) {
		fenStr := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
		got, err := parsePiecesStr(fenStr)

		assert.Nil(t, err)

		want := [64]piece{
			wR, wN, wB, wQ, wK, wB, wN, wR,
			wP, wP, wP, wP, wP, wP, wP, wP,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, 00,
			bP, bP, bP, bP, bP, bP, bP, bP,
			bR, bN, bB, bQ, bK, bB, bN, bR,
		}
		assert.Equal(t, want, got)
	})
	t.Run("it can parse gap numbers", func(t *testing.T) {
		fenStr := "1ppppppp/2pppppp/3ppppp/4pppp/5ppp/6pp/7p/8"
		got, err := parsePiecesStr(fenStr)

		assert.Nil(t, err)

		want := [64]piece{
			00, 00, 00, 00, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 00, 00, bP,
			00, 00, 00, 00, 00, 00, bP, bP,
			00, 00, 00, 00, 00, bP, bP, bP,
			00, 00, 00, 00, bP, bP, bP, bP,
			00, 00, 00, bP, bP, bP, bP, bP,
			00, 00, bP, bP, bP, bP, bP, bP,
			00, bP, bP, bP, bP, bP, bP, bP,
		}
		assert.Equal(t, want, got)
	})
}

func TestParseCastlePerms(t *testing.T) {
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

func TestParseFen(t *testing.T) {
	t.Run("it can parse a starting position", func(t *testing.T) {
		fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		state, err := ParseFen(fen)

		assert.Nil(t, err)
		assert.Equal(t, WHITE, state.side)
		assert.Equal(t, &castlePerm{val: CASTLE_PERMS_ALL}, state.castlePerm)
		assert.Equal(t, -1, state.enPas)
		assert.Equal(t, 0, state.fiftyMove)
		assert.Equal(t, 0, state.hisPly)

		want := &board120{
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, wR, wN, wB, wQ, wK, wB, wN, wR, -1,
			-1, wP, wP, wP, wP, wP, wP, wP, wP, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, bP, bP, bP, bP, bP, bP, bP, bP, -1,
			-1, bR, bN, bB, bQ, bK, bB, bN, bR, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		}

		assert.Equal(t, want, state.pieces)
	})
}
