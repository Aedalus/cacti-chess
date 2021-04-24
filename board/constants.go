package board

import "math/rand"

type piece int8
type file int8
type rank int8

const BOARD_SQ_NUMBER = 120
const NO_SQ = -1

const (
	EMPTY piece = iota
	wP
	wN
	wB
	wR
	wQ
	wK
	bP
	bN
	bB
	bR
	bQ
	bK
)

const (
	FILE_A file = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

const (
	RANK_1 rank = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
)

const (
	WHITE = iota
	BLACK
	BOTH
)

//const (
//	A1 = 22
//	B1 = 23
//	C1 = 24
//	D1 = 25
//	E1 = 26
//	F1 = 27
//	G1 = 28
//	H1 = 29
//	A2 int8 = iota
//)

// these maps contain random uint64s for different combinations,
// and they might not be the same through reboots. used to
// generate the posKey from a board state
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