package position

import (
	"math/rand"
)

const BOARD_SQ_NUMBER = 120
const NO_SQ = -1

const (
	FILE_A = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

const (
	RANK_1 = iota
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

// Piece positions for a 120 engine
const (
	// A1 through H1
	A1 = 21
	B1 = 22
	C1 = 23
	D1 = 24
	E1 = 25
	F1 = 26
	G1 = 27
	H1 = 28

	// A2 through H2
	A2 = 31
	B2 = 32
	C2 = 33
	D2 = 34
	E2 = 35
	F2 = 36
	G2 = 37
	H2 = 38

	// A3 through H3
	A3 = 41
	B3 = 42
	C3 = 43
	D3 = 44
	E3 = 45
	F3 = 46
	G3 = 47
	H3 = 48

	// A4 through H4
	A4 = 51
	B4 = 52
	C4 = 53
	D4 = 54
	E4 = 55
	F4 = 56
	G4 = 57
	H4 = 58

	// A5 through H5
	A5 = 61
	B5 = 62
	C5 = 63
	D5 = 64
	E5 = 65
	F5 = 66
	G5 = 67
	H5 = 68

	// A6 through H6
	A6 = 71
	B6 = 72
	C6 = 73
	D6 = 74
	E6 = 75
	F6 = 76
	G6 = 77
	H6 = 78

	// A7 through H7
	A7 = 81
	B7 = 82
	C7 = 83
	D7 = 84
	E7 = 85
	F7 = 86
	G7 = 87
	H7 = 88

	// A8 through H8
	A8 = 91
	B8 = 92
	C8 = 93
	D8 = 94
	E8 = 95
	F8 = 96
	G8 = 97
	H8 = 98
)

// directions for Piece movement, used to calculate attacks
var dirKnight = [8]int{-8, -19, -21, -12, 8, 19, 21, 12}
var dirRook = [4]int{-1, -10, 1, 10}
var dirBishop = [4]int{-9, -11, 11, 9}
var dirKing = [8]int{-1, -10, 1, 10, -9, -11, 11, 9}

// these maps contain random uint64s for different combinations,
// and they might not be the same through reboots. used to
// generate the posKey from a engine state
var hashPieceKeys [13][BOARD_SQ_NUMBER]uint64 // Piece type/position
var hashSideKey uint64                        // used if white's turn
var hashCastleKeys [16]uint64                 // castleKeys

// used for lookups of rank/file quickly given a board number
// i.e. rankLookups[21] = 1, fileLookups = 1 (A1)
// NO_SQ if not on the core 8x8
var rankLookups [BOARD_SQ_NUMBER]int
var fileLookups [BOARD_SQ_NUMBER]int

func init() {
	// initialize the hash lookups pieces
	for piece := 0; piece < 13; piece++ {
		for position := 0; position < BOARD_SQ_NUMBER; position++ {
			hashPieceKeys[piece][position] = rand.Uint64()
		}
	}

	hashSideKey = rand.Uint64()

	for i := 0; i < 16; i++ {
		hashCastleKeys[i] = rand.Uint64()
	}

	// initialize the rank/file lookups
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		rankLookups[i] = NO_SQ
		fileLookups[i] = NO_SQ
	}

	for r := RANK_1; r <= RANK_8; r++ {
		for f := FILE_A; f <= FILE_H; f++ {
			sq120 := fileRankToSq(f, r)
			fileLookups[sq120] = f
			rankLookups[sq120] = r
		}
	}
}
