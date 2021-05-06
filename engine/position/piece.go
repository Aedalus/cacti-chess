package position

type Piece int8

const PIECE_COUNT = 13
const (
	EMPTY Piece = iota
	PwP
	PwN
	PwB
	PwR
	PwQ
	PwK
	PbP
	PbN
	PbB
	PbR
	PbQ
	PbK
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
var pieceLookups = map[Piece]pieceMetadata{
	EMPTY: {
		color: BOTH,
	},
	PwP: {
		color: WHITE,
		value: 100,
	},
	PwN: {
		color:   WHITE,
		isBig:   true,
		isMinor: true,
		value:   325,
		dir:     []int{-8, -19, -21, -12, 8, 19, 21, 12},
	},
	PwB: {
		color:           WHITE,
		isBig:           true,
		isMinor:         true,
		value:           330,
		slides:          true,
		isBishopOrQueen: true,
		dir:             []int{-9, -11, 9, 11},
	},
	PwR: {
		color:         WHITE,
		isBig:         true,
		isMajor:       true,
		value:         550,
		slides:        true,
		isRookOrQueen: true,
		dir:           []int{-1, -10, 1, 10},
	},
	PwQ: {
		color:           WHITE,
		isBig:           true,
		isMajor:         true,
		value:           1000,
		slides:          true,
		isRookOrQueen:   true,
		isBishopOrQueen: true,
		dir:             []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	PwK: {
		color:   WHITE,
		isBig:   true,
		isMajor: true,
		value:   50000,
		dir:     []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	PbP: {
		color: BLACK,
		value: 100,
	},
	PbN: {
		color:   BLACK,
		isBig:   true,
		isMinor: true,
		value:   325,
		dir:     []int{-8, -19, -21, -12, 8, 19, 21, 12},
	},
	PbB: {
		color:           BLACK,
		isBig:           true,
		isMinor:         true,
		value:           330,
		slides:          true,
		isBishopOrQueen: true,
		dir:             []int{-9, -11, 9, 11},
	},
	PbR: {
		color:         BLACK,
		isBig:         true,
		isMajor:       true,
		value:         550,
		slides:        true,
		isRookOrQueen: true,
		dir:           []int{-1, -10, 1, 10},
	},
	PbQ: {
		color:           BLACK,
		isBig:           true,
		isMajor:         true,
		value:           1000,
		slides:          true,
		isRookOrQueen:   true,
		isBishopOrQueen: true,
		dir:             []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
	PbK: {
		color:   BLACK,
		isBig:   true,
		isMajor: true,
		value:   50000,
		dir:     []int{-1, -10, 1, 10, -9, -11, 11, 9},
	},
}

func (p Piece) String() string {
	switch p {
	case EMPTY:
		return "."
	case PwR:
		return "R"
	case PwP:
		return "P"
	case PwN:
		return "N"
	case PwB:
		return "B"
	case PwQ:
		return "Q"
	case PwK:
		return "K"
	case PbP:
		return "p"
	case PbN:
		return "n"
	case PbB:
		return "b"
	case PbR:
		return "r"
	case PbQ:
		return "q"
	case PbK:
		return "k"
	default:
		return "UNKNOWN"
	}
}
