package engine

import (
	"fmt"
	"reflect"
	"strings"
)

type undo struct {
	move       int
	castlePerm castlePerm
	enPas      uint
	fiftyMove  uint
	posKey     uint64
}

type Position struct {
	pieces *board120 // Source of truth for pieces
	side   int       // white/black

	pawns [3]*bitboard64 // white/black/both quick pawn lookups

	// piece list for fast lookup
	// pieceList[wN][0] = E1 etc
	pieceList [13][10]int

	kingSq [2]int // king quick lookups

	castlePerm *castlePerm // permissions to castle

	enPas int // if en passant is available

	fiftyMove int // 50 move counter (100 since we're using half moves)

	searchPly uint // how far we are in the search

	posKey uint64 // unique index for position

	pieceCount    [13]int // count of all pieces on the board
	bigPieceCount [2]int  // white/black non-pawn pieces
	majPieceCount [2]int  // queen/rook
	minPieceCount [2]int  // bishop/knight
	materialCount [2]int  // count of material

	// history
	halfMoveCount int // how many half moves have been made in the whole game
	history       *[2048]undo
}

// GenPosKey generates a statistically unique uint64
// key for the current state of the engine
func (s Position) GenPosKey() uint64 {
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

func (s *Position) updateListCaches() {
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		piece := s.pieces[i]
		if piece != NO_SQ && piece != EMPTY {
			color := pieceLookups[piece].color
			if pieceLookups[piece].isBig {
				s.bigPieceCount[color]++
			}
			if pieceLookups[piece].isMajor {
				s.majPieceCount[color]++
			}
			if pieceLookups[piece].isMinor {
				s.minPieceCount[color]++
			}

			s.materialCount[color] += pieceLookups[piece].value

			// update the pieceList, then increment the counter
			// conceptually like
			// [wP][0] = A2
			// [wP][1] = B2 etc
			curPieceCount := s.pieceCount[piece]
			s.pieceList[piece][curPieceCount] = i
			s.pieceCount[piece]++

			// update king positions
			if piece == wK {
				s.kingSq[WHITE] = i
			}
			if piece == bK {
				s.kingSq[BLACK] = i
			}

			// update pawn boards
			if piece == wP {
				s.pawns[WHITE].set(SQ64(i))
				s.pawns[BOTH].set(SQ64(i))
			}
			if piece == bP {
				s.pawns[BLACK].set(SQ64(i))
				s.pawns[BOTH].set(SQ64(i))
			}
		}
	}
}

// will panic if the cache isn't right. used for debugging
func (s *Position) assertCache() {
	// temporary values we recompute to check against
	t_pieceCount := [13]int{}
	t_pieceList := [13][10]int{}
	t_bigPieceCount := [2]int{}
	t_majPieceCount := [2]int{}
	t_minPieceCount := [2]int{}
	t_kingSq := [2]int{}
	t_materialCount := [2]int{}

	t_pawns := [3]*bitboard64{
		&bitboard64{},
		&bitboard64{},
		&bitboard64{},
	}

	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		p := s.pieces[i]
		if p != NO_SQ && p != EMPTY {
			color := pieceLookups[p].color
			if pieceLookups[p].isBig {
				t_bigPieceCount[color]++
			}
			if pieceLookups[p].isMajor {
				t_majPieceCount[color]++
			}
			if pieceLookups[p].isMinor {
				t_minPieceCount[color]++
			}

			t_materialCount[color] += pieceLookups[p].value

			// update the pieceList, then increment the counter
			// conceptually like
			// [wP][0] = A2
			// [wP][1] = B2 etc
			curPieceCount := t_pieceCount[p]
			t_pieceList[p][curPieceCount] = i
			t_pieceCount[p]++

			// update king positions
			if p == wK {
				t_kingSq[WHITE] = i
			}
			if p == bK {
				t_kingSq[BLACK] = i
			}

			// update pawn boards
			if p == wP {
				t_pawns[WHITE].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
			if p == bP {
				t_pawns[BLACK].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
		}
	}

	if !reflect.DeepEqual(t_pieceCount, s.pieceCount) {
		panic(fmt.Errorf("pieceCount - got %v, want %v", s.pieceCount, t_pieceCount))
	}
	if !reflect.DeepEqual(t_pieceList, s.pieceList) {
		panic(fmt.Errorf("pieceList - got %v, want %v", s.pieceList, t_pieceList))
	}
	if !reflect.DeepEqual(t_bigPieceCount, s.bigPieceCount) {
		panic(fmt.Errorf("bigPieceCount - got %v, want %v", s.bigPieceCount, t_bigPieceCount))
	}
	if !reflect.DeepEqual(t_majPieceCount, s.majPieceCount) {
		panic(fmt.Errorf("majPieceCount - got %v want %v", s.majPieceCount, t_majPieceCount))
	}
	if !reflect.DeepEqual(t_minPieceCount, s.minPieceCount) {
		panic(fmt.Errorf("minPieceCount - got %v want %v", s.minPieceCount, t_minPieceCount))
	}
	if !reflect.DeepEqual(t_kingSq, s.kingSq) {
		panic(fmt.Errorf("kingSq - got %v want %v", s.kingSq, t_kingSq))
	}
	if !reflect.DeepEqual(t_materialCount, s.materialCount) {
		panic(fmt.Errorf("materialCount - got %v want %v", s.materialCount, t_materialCount))
	}
	if !reflect.DeepEqual(t_pawns, s.pawns) {
		panic(fmt.Errorf("pawns - got %v want %v", s.pawns, t_pawns))
	}
}

// Reset will re-initialize the engine to an empty state
func (s *Position) Reset() {
	// initialize piece array
	s.pieces = &board120{}
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		s.pieces[i] = NO_SQ
	}
	for i := 0; i < 64; i++ {
		s.pieces[SQ120(i)] = EMPTY
	}

	// castle perms
	s.castlePerm = &castlePerm{CASTLE_PERMS_NONE}

	// piece counts
	for i := 0; i < 13; i++ {
		s.pieceCount[i] = 0
	}
	for i := 0; i < 2; i++ {
		s.bigPieceCount[i] = 0
		s.majPieceCount[i] = 0
		s.minPieceCount[i] = 0
		s.materialCount[i] = 0
		s.pawns[i] = &bitboard64{0}
	}

	s.pawns[WHITE] = &bitboard64{}
	s.pawns[BLACK] = &bitboard64{}
	s.pawns[BOTH] = &bitboard64{}

	s.kingSq[WHITE] = NO_SQ
	s.kingSq[BLACK] = NO_SQ

	s.side = BOTH
	s.enPas = NO_SQ
	s.fiftyMove = 0
	s.searchPly = 0
	s.halfMoveCount = 0
	s.history = &[2048]undo{}
	s.posKey = 0
}

func (s *Position) IsSquareAttacked(sq, attackingColor int) bool {

	// todo - delete these asserts for optimization
	if fileLookups[sq] == NO_SQ {
		panic(fmt.Errorf("invalid square: %v", sq))
	}
	if attackingColor != WHITE && attackingColor != BLACK {
		panic(fmt.Errorf("unexpected color: %v", attackingColor))
	}
	s.assertCache()

	// pawns
	if attackingColor == WHITE {
		if s.pieces[sq-11] == wP || s.pieces[sq-9] == wP {
			return true
		}
	} else {
		if s.pieces[sq+11] == bP || s.pieces[sq+9] == bP {
			return true
		}
	}

	// knights
	for _, n := range dirKnight { // has the offsets for the knight jumps
		kSq := sq + n
		if attackingColor == WHITE && s.pieces[kSq] == wN {
			return true
		} else if attackingColor == BLACK && s.pieces[kSq] == bN {
			return true
		}
	}

	// rook/queen
	for _, dir := range dirRook {
		tsq := sq + dir
		p := s.pieces[tsq]
		// only go until we're off the board
		for p != NO_SQ {
			if p != EMPTY {
				if pieceLookups[p].isRookOrQueen && pieceLookups[p].color == attackingColor {
					return true
				}
				break // break out if we hit something that wasn't a rook or queen
			}
			tsq += dir // move to the next square
			p = s.pieces[tsq]
		}
	}

	// bishop/queen
	for _, dir := range dirBishop {
		tsq := sq + dir
		p := s.pieces[tsq]
		// only go until we're off the board
		for p != NO_SQ {
			if p != EMPTY {
				if pieceLookups[p].isBishopOrQueen && pieceLookups[p].color == attackingColor {
					return true
				}
				break // break out if we hit something else
			}
			tsq += dir // move to the next square
			p = s.pieces[tsq]
		}
	}

	for _, dir := range dirKing {
		tsq := sq + dir
		p := s.pieces[tsq]
		if p == wK && attackingColor == WHITE {
			return true
		}
		if p == bK && attackingColor == BLACK {
			return true
		}
	}

	// only return false once we've ruled out everything
	return false
}

func (s Position) String() string {
	output := strings.Builder{}
	output.WriteString(s.PrintBoard())
	output.WriteString("--------\n")
	output.WriteString(fmt.Sprintf("side: %v\n", s.side))
	output.WriteString(fmt.Sprintf("enPas: %v\n", s.enPas))
	output.WriteString(fmt.Sprintf("castle: %v\n", s.castlePerm))
	output.WriteString(fmt.Sprintf("posKey: %x\n", s.posKey))

	return output.String()
}
func (s Position) PrintBoard() string {
	output := strings.Builder{}
	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			sq := fileRankToSq(int(f), int(r))
			piece := s.pieces[sq]
			output.WriteString(piece.String())
		}
		output.WriteString("\n")
	}
	return output.String()
}

// PrintSQsAttackedBySide returns a string representation
// of all squares attacked by a given side
func (s Position) PrintAttackBoard(attackingSide int) string {
	output := strings.Builder{}

	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			sq := fileRankToSq(f, r)
			if s.IsSquareAttacked(sq, attackingSide) {
				output.WriteString("X")
			} else {
				output.WriteString("-")
			}
		}
		output.WriteString("\n")
	}
	return output.String()
}
