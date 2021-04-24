package board

import "math/rand"

// these maps contain random uint64s for different combinations,
// and they might not be the same through reboots.

var hashPieceKeys [13][BOARD_SQ_NUMBER]uint64 // piece type/position
var hashSideKey uint64 // used if white's turn
var hashCastleKeys [16]uint64 // castleKeys

func init(){
	// initialize the pieces
	for piece := 0; piece < 13; piece++{
		for position := 0; position < BOARD_SQ_NUMBER; position++ {
			hashPieceKeys[piece][position] = rand.Uint64()
		}
	}

	hashSideKey = rand.Uint64()

	for i := 0; i < 16; i++ {
		hashCastleKeys[i] = rand.Uint64()
	}
}

// hashes the current state of the board
func (s state) genPosKey() uint64 {
	var finalKey uint64 = 0
	var p piece = piece(0)

	// pieces
	for sq := 0; sq < BOARD_SQ_NUMBER; sq++ {
		p = s.pieces[sq]
		if p != NO_SQ && p != EMPTY {
			finalKey ^= hashPieceKeys[p][sq]
		}
	}

	// side
	if s.side == WHITE {
		finalKey ^= hashSideKey
	}

	// enPas
	if s.enPas != NO_SQ {
		finalKey ^= hashPieceKeys[EMPTY][s.enPas]
	}

	// castle keys
	finalKey ^= hashCastleKeys[s.castlePerm.val]

	return finalKey
}
