package engine

type piece int8

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

// returns if a piece is major (rook/king/queen)
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

// returns if a piece is minor (bishop/knight)
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

// returns the value of a piece
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

// if a piece is a knight
//var isKnight = map[piece]bool{
//
//	EMPTY: false,
//	wP:    false,
//	wN:    true,
//	wB:    false,
//	wR:    false,
//	wQ:    false,
//	wK:    false,
//	bP:    false,
//	bN:    true,
//	bB:    false,
//	bR:    false,
//	bQ:    false,
//	bK:    false,
//}

// returns if a piece is a king
//var isKing = map[piece]bool{
//
//	EMPTY: false,
//	wP:    false,
//	wN:    false,
//	wB:    false,
//	wR:    false,
//	wQ:    false,
//	wK:    true,
//	bP:    false,
//	bN:    false,
//	bB:    false,
//	bR:    false,
//	bQ:    false,
//	bK:    true,
//}

// returns if a piece is a rook/queen
var isRookOrQueen = map[piece]bool{

	EMPTY: false,
	wP:    false,
	wN:    false,
	wB:    false,
	wR:    true,
	wQ:    true,
	wK:    false,
	bP:    false,
	bN:    false,
	bB:    false,
	bR:    true,
	bQ:    true,
	bK:    false,
}

// returns if a piece is a bishop/queen
var isBishopOrQueen = map[piece]bool{

	EMPTY: false,
	wP:    false,
	wN:    false,
	wB:    true,
	wR:    false,
	wQ:    true,
	wK:    false,
	bP:    false,
	bN:    false,
	bB:    true,
	bR:    false,
	bQ:    true,
	bK:    false,
}
