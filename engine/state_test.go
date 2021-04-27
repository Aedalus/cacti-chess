package engine

import (
	"fmt"
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

func TestState_Reset(t *testing.T) {
	want := &State{
		pieces: &board120{
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, 00, 00, 00, 00, 00, 00, 00, 00, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		},
		side: BOTH,
		pawns: [3]*bitboard64{
			&bitboard64{},
			&bitboard64{},
			&bitboard64{},
		},
		pieceList:     [13][10]int{},
		kingSq:        [2]int{NO_SQ, NO_SQ},
		castlePerm:    &castlePerm{CASTLE_PERMS_NONE},
		enPas:         NO_SQ,
		fiftyMove:     0,
		searchPly:     0,
		posKey:        0,
		pieceCount:    [13]int{},
		bigPieceCount: [2]int{},
		majPieceCount: [2]int{},
		minPieceCount: [2]int{},
		materialCount: [2]int{},
		halfMoveCount: 0,
		history:       &[2048]undo{},
	}
	sample := &State{
		pieces:        &board120{},
		side:          BOTH,
		pawns:         [3]*bitboard64{},
		pieceList:     [13][10]int{},
		kingSq:        [2]int{},
		castlePerm:    &castlePerm{CASTLE_PERMS_NONE},
		enPas:         B2,
		fiftyMove:     20,
		searchPly:     70,
		posKey:        3,
		pieceCount:    [13]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		bigPieceCount: [2]int{1, 1},
		majPieceCount: [2]int{1, 1},
		minPieceCount: [2]int{1, 1},
		materialCount: [2]int{1, 1},
		halfMoveCount: 70,
		history:       &[2048]undo{},
	}
	sample.Reset()
	assert.Equal(t, sample, want)
}

func TestState_String(t *testing.T) {
	fen := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"
	state, err := parseFen(fen)

	assert.Nil(t, err)
	fmt.Println(state)
}

func TestUpdateListsMaterial(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	state, err := parseFen(fen)

	assert.Nil(t, err)

	state.updateListCaches()
	assert.Equal(t, [2]int{8, 8}, state.bigPieceCount)
	assert.Equal(t, [2]int{4, 4}, state.majPieceCount)
	assert.Equal(t, [2]int{4, 4}, state.minPieceCount)
	assert.Equal(t, [2]int{54220, 54220}, state.materialCount)

	// check piece counts empty, wP, wB, wN, etc
	assert.Equal(t, [13]int{
		0, 8, 2, 2, 2, 1, 1, 8, 2, 2, 2, 1, 1,
	}, state.pieceCount)

	// spot check some white pieces
	assert.Equal(t, [10]int{
		31, 32, 33, 34, 35, 36, 37, 38, 0, 0,
	}, state.pieceList[wP])

	assert.Equal(t, [10]int{
		21, 28, 0, 0, 0, 0, 0, 0, 0, 0,
	}, state.pieceList[wR])
}