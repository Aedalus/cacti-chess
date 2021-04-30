package engine

type piece int8

const PIECE_COUNT = 13
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

type pieceMetadata struct {
	isBig           bool  // everything but pawns + empty
	isMajor         bool  // rooks, queens, kings
	isMinor         bool  // bishops, knights
	value           int   // material value
	color           int   // WHITE, BLACK, BOTH
	slides          bool  // bishop/rook/queen
	isRookOrQueen   bool  // if it's a rook or queen
	isBishopOrQueen bool  // if it's a bishop or queen
	dir             []int // movement directions
}

// used for fast lookups without calling functions
var pieceLookups = map[piece]pieceMetadata{
	EMPTY: {
		color: BOTH,
	},
	wP: {
		color: WHITE,
		value: 100,
	},
	wN: {
		color:   WHITE,
		isBig:   true,
		isMinor: true,
		value:   325,
		dir:     []int{-8, -19, -21, -12, 8, 19, 21, 12},
	},
	wB: {
		color:           WHITE,
		isBig:           true,
		isMinor:         true,
		value:           330,
		slides:          true,
		isBishopOrQueen: true,
		dir:             []int{-9, -11, 9, 11},
	},
	wR: {
		color:         WHITE,
		isBig:         true,
		isMajor:       true,
		value:         550,
		slides:        true,
		isRookOrQueen: true,
		dir:           []int{-1, -10, 1, 10},
	},
	wQ: {
		color:           WHITE,
		isBig:           true,
		isMajor:         true,
		value:           1000,
		slides:          true,
		isRookOrQueen:   true,
		isBishopOrQueen: true,
		dir:             []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	wK: {
		color:   WHITE,
		isBig:   true,
		isMajor: true,
		value:   50000,
		dir:     []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	bP: {
		color: BLACK,
		value: 100,
	},
	bN: {
		color:   BLACK,
		isBig:   true,
		isMinor: true,
		value:   325,
		dir:     []int{-8, -19, -21, -12, 8, 19, 21, 12},
	},
	bB: {
		color:           BLACK,
		isBig:           true,
		isMinor:         true,
		value:           330,
		slides:          true,
		isBishopOrQueen: true,
		dir:             []int{-9, -11, 9, 11},
	},
	bR: {
		color:         BLACK,
		isBig:         true,
		isMajor:       true,
		value:         550,
		slides:        true,
		isRookOrQueen: true,
		dir:           []int{-1, -10, 1, 10},
	},
	bQ: {
		color:           BLACK,
		isBig:           true,
		isMajor:         true,
		value:           1000,
		slides:          true,
		isRookOrQueen:   true,
		isBishopOrQueen: true,
		dir:             []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	bK: {
		color:   BLACK,
		isBig:   true,
		isMajor: true,
		value:   50000,
		dir:     []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
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
