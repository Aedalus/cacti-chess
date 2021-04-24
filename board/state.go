package board

type undo struct {
	move int
	castlePerm castlePerm
	enPas uint
	fiftyMove uint
	posKey uint64
}

type state struct {
	pieces *board120 // Source of truth for pieces
	side int // white/black

	pawnsWhite *bitboard64 // pawn quick lookups
	pawnsBlack *bitboard64
	pawnsBoth *bitboard64

	// piece list for fast lookup
	// pieceList[wN][0] = E1 etc
	pieceList [13][10]int64

	kingSqWhite uint // king quick lookups
	kingSqBlack uint

	castlePerm *castlePerm // permissions to castle

	enPas int // if en passant is available

	fiftyMove uint // 50 move counter (100 since we're using half moves)

	ply uint // how far we are in the search

	posKey uint64 // unique index for position

	pieceCount [13]int // count of all pieces on the board

	bigPieceWhiteCount int // Count of r + k + b + q
	bigPieceBlackCount int
	bigPieceBothCount int

	majPieceWhiteCount int // count of r + q
	majPieceBlackCount int
	majPieceBothCount int

	minPieceWhiteCount int // count of b + k
	minPieceBlackCount int
	minPieceBothCount int

	// history
	hisPly uint // how many half moves have been made in the whole game
	history *[2048]undo
}
