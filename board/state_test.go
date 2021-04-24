package board

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_GenPosKey(t *testing.T) {
	t.Run("it generates a unique id for each item in the piece array", func(t *testing.T) {

		// store the key -> input to check for unique keys
		result := map[uint64]int{}

		// Check for all pieces
		for p := 1; p < 13; p++ {

			for sq := 0; sq < BOARD_SQ_NUMBER; sq++ {
				s := State{
					pieces:     &board120{},
					castlePerm: &castlePerm{CASTLE_PERMS_ALL},
				}
				s.pieces[sq] = piece(p)
				key := s.GenPosKey()
				result[key] = sq
			}

			// should have 120 unique keys per piece:pos combo
			assert.Equal(t, 120*p, len(result))
		}
	})

	t.Run("it generates a unique key based on current side", func(t *testing.T) {
		// store the key -> input to check for unique keys
		result := map[uint64]int{}
		s := State{
			pieces:     &board120{},
			castlePerm: &castlePerm{CASTLE_PERMS_ALL},
		}

		whiteKey := s.GenPosKey()
		result[whiteKey] = WHITE

		s.side = BLACK
		blackKey := s.GenPosKey()
		result[blackKey] = BLACK

		assert.Equal(t, 2, len(result))
	})
	t.Run("it generates a unique key based on castlePerms", func(t *testing.T) {
		// store the key -> input to check for unique keys
		result := map[uint64]int{}

		for cp := 0; cp < 16; cp++ {
			s := State{
				pieces:     &board120{},
				castlePerm: &castlePerm{cp},
			}
			key := s.GenPosKey()
			result[key] = cp
		}
		assert.Equal(t, 16, len(result))
	})

}
