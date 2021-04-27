package engine

import "math/rand"

type piece int8

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

// everything but pawns and empty should be big
var pieceBigMap = map[piece]bool{
	EMPTY: false,
	wP:    false,
	wN:    true,
	wB:    true,
	wR:    true,
	wQ:    true,
	wK:    true,
	bP:    false,
	bN:    true,
	bB:    true,
	bR:    true,
	bQ:    true,
	bK:    true,
}

// major pieces are rooks/queens/kings
var pieceMajorMap = map[piece]bool{

	EMPTY: false,
	wP:    false,
	wN:    false,
	wB:    false,
	wR:    true,
	wQ:    true,
	wK:    true,
	bP:    false,
	bN:    false,
	bB:    false,
	bR:    true,
	bQ:    true,
	bK:    true,
}

// minor pieces are bishops/knights
var pieceMinorMap = map[piece]bool{

	EMPTY: false,
	wP:    false,
	wN:    true,
	wB:    true,
	wR:    false,
	wQ:    false,
	wK:    false,
	bP:    false,
	bN:    true,
	bB:    true,
	bR:    false,
	bQ:    false,
	bK:    false,
}

// piece values
var pieceValueMap = map[piece]int{
	EMPTY: 0,
	wP:    100,
	wN:    325,
	wB:    335,
	wR:    550,
	wQ:    1000,
	wK:    50000,
	bP:    100,
	bN:    325,
	bB:    335,
	bR:    550,
	bQ:    1000,
	bK:    50000,
}

var pieceColorMap = map[piece]int{
	EMPTY: BOTH,
	wP:    WHITE,
	wN:    WHITE,
	wB:    WHITE,
	wR:    WHITE,
	wQ:    WHITE,
	wK:    WHITE,
	bP:    BLACK,
	bN:    BLACK,
	bB:    BLACK,
	bR:    BLACK,
	bQ:    BLACK,
	bK:    BLACK,
}

func (p piece) String() string {
	switch p {
	case EMPTY:
		return "."
	case wR:
		return "R"
	case wP:
		return "P"
	case wN:
		return "N"
	case wB:
		return "B"
	case wQ:
		return "Q"
	case wK:
		return "K"
	case bP:
		return "p"
	case bN:
		return "n"
	case bB:
		return "b"
	case bR:
		return "r"
	case bQ:
		return "q"
	case bK:
		return "k"
	default:
		return "UNKNOWN"
	}
}

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

// these maps contain random uint64s for different combinations,
// and they might not be the same through reboots. used to
// generate the posKey from a engine state
var hashPieceKeys [13][BOARD_SQ_NUMBER]uint64 // piece type/position
var hashSideKey uint64                        // used if white's turn
var hashCastleKeys [16]uint64                 // castleKeys

// used for lookups of rank/file quickly given a board number
// i.e. rankLookups[21] = 1, fileLookups = 1 (A1)
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
	for r := RANK_1; r <= RANK_8; r++ {
		for f := FILE_A; f <= FILE_H; f++ {
			sq120 := fileRankToSq(f, r)
			fileLookups[sq120] = f
			rankLookups[sq120] = r
		}
	}
}
