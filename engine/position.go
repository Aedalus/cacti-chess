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
func (p Position) GenPosKey() uint64 {
	var finalKey uint64 = 0
	var pce piece = piece(0)

	// pieces
	for sq := 0; sq < BOARD_SQ_NUMBER; sq++ {
		pce = p.pieces[sq]
		if pce != NO_SQ && pce != EMPTY {
			finalKey ^= hashPieceKeys[pce][sq]
		}
	}

	// side
	if p.side == WHITE {
		finalKey ^= hashSideKey
	}

	// enPas
	if p.enPas != NO_SQ {
		finalKey ^= hashPieceKeys[EMPTY][p.enPas]
	}

	// castle keys
	finalKey ^= hashCastleKeys[p.castlePerm.val]

	return finalKey
}

func (p *Position) updateListCaches() {
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		pce := p.pieces[i]
		if pce != NO_SQ && pce != EMPTY {
			color := pieceLookups[pce].color
			if pieceLookups[pce].isBig {
				p.bigPieceCount[color]++
			}
			if pieceLookups[pce].isMajor {
				p.majPieceCount[color]++
			}
			if pieceLookups[pce].isMinor {
				p.minPieceCount[color]++
			}

			p.materialCount[color] += pieceLookups[pce].value

			// update the pieceList, then increment the counter
			// conceptually like
			// [wP][0] = A2
			// [wP][1] = B2 etc
			curPieceCount := p.pieceCount[pce]
			p.pieceList[pce][curPieceCount] = i
			p.pieceCount[pce]++

			// update king positions
			if pce == wK {
				p.kingSq[WHITE] = i
			}
			if pce == bK {
				p.kingSq[BLACK] = i
			}

			// update pawn boards
			if pce == wP {
				p.pawns[WHITE].set(SQ64(i))
				p.pawns[BOTH].set(SQ64(i))
			}
			if pce == bP {
				p.pawns[BLACK].set(SQ64(i))
				p.pawns[BOTH].set(SQ64(i))
			}
		}
	}
}

// will panic if the cache isn't right. used for debugging
func (p *Position) assertCache() {
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
		pce := p.pieces[i]
		if pce != NO_SQ && pce != EMPTY {
			color := pieceLookups[pce].color
			if pieceLookups[pce].isBig {
				t_bigPieceCount[color]++
			}
			if pieceLookups[pce].isMajor {
				t_majPieceCount[color]++
			}
			if pieceLookups[pce].isMinor {
				t_minPieceCount[color]++
			}

			t_materialCount[color] += pieceLookups[pce].value

			// update the pieceList, then increment the counter
			// conceptually like
			// [wP][0] = A2
			// [wP][1] = B2 etc
			curPieceCount := t_pieceCount[pce]
			t_pieceList[pce][curPieceCount] = i
			t_pieceCount[pce]++

			// update king positions
			if pce == wK {
				t_kingSq[WHITE] = i
			}
			if pce == bK {
				t_kingSq[BLACK] = i
			}

			// update pawn boards
			if pce == wP {
				t_pawns[WHITE].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
			if pce == bP {
				t_pawns[BLACK].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
		}
	}

	if !reflect.DeepEqual(t_pieceCount, p.pieceCount) {
		panic(fmt.Errorf("pieceCount - got %v, want %v", p.pieceCount, t_pieceCount))
	}
	if !reflect.DeepEqual(t_pieceList, p.pieceList) {
		panic(fmt.Errorf("pieceList - got %v, want %v", p.pieceList, t_pieceList))
	}
	if !reflect.DeepEqual(t_bigPieceCount, p.bigPieceCount) {
		panic(fmt.Errorf("bigPieceCount - got %v, want %v", p.bigPieceCount, t_bigPieceCount))
	}
	if !reflect.DeepEqual(t_majPieceCount, p.majPieceCount) {
		panic(fmt.Errorf("majPieceCount - got %v want %v", p.majPieceCount, t_majPieceCount))
	}
	if !reflect.DeepEqual(t_minPieceCount, p.minPieceCount) {
		panic(fmt.Errorf("minPieceCount - got %v want %v", p.minPieceCount, t_minPieceCount))
	}
	if !reflect.DeepEqual(t_kingSq, p.kingSq) {
		panic(fmt.Errorf("kingSq - got %v want %v", p.kingSq, t_kingSq))
	}
	if !reflect.DeepEqual(t_materialCount, p.materialCount) {
		panic(fmt.Errorf("materialCount - got %v want %v", p.materialCount, t_materialCount))
	}
	if !reflect.DeepEqual(t_pawns, p.pawns) {
		panic(fmt.Errorf("pawns - got %v want %v", p.pawns, t_pawns))
	}
}

// Reset will re-initialize the engine to an empty state
func (p *Position) Reset() {
	// initialize piece array
	p.pieces = &board120{}
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		p.pieces[i] = NO_SQ
	}
	for i := 0; i < 64; i++ {
		p.pieces[SQ120(i)] = EMPTY
	}

	// castle perms
	p.castlePerm = &castlePerm{CASTLE_PERMS_NONE}

	// piece counts
	for i := 0; i < 13; i++ {
		p.pieceCount[i] = 0
	}
	for i := 0; i < 2; i++ {
		p.bigPieceCount[i] = 0
		p.majPieceCount[i] = 0
		p.minPieceCount[i] = 0
		p.materialCount[i] = 0
		p.pawns[i] = &bitboard64{0}
	}

	p.pawns[WHITE] = &bitboard64{}
	p.pawns[BLACK] = &bitboard64{}
	p.pawns[BOTH] = &bitboard64{}

	p.kingSq[WHITE] = NO_SQ
	p.kingSq[BLACK] = NO_SQ

	p.side = BOTH
	p.enPas = NO_SQ
	p.fiftyMove = 0
	p.searchPly = 0
	p.halfMoveCount = 0
	p.history = &[2048]undo{}
	p.posKey = 0
}

func (p *Position) IsSquareAttacked(sq, attackingColor int) bool {

	// todo - delete these asserts for optimization
	if fileLookups[sq] == NO_SQ {
		panic(fmt.Errorf("invalid square: %v", sq))
	}
	if attackingColor != WHITE && attackingColor != BLACK {
		panic(fmt.Errorf("unexpected color: %v", attackingColor))
	}
	p.assertCache()

	// pawns
	if attackingColor == WHITE {
		if p.pieces[sq-11] == wP || p.pieces[sq-9] == wP {
			return true
		}
	} else {
		if p.pieces[sq+11] == bP || p.pieces[sq+9] == bP {
			return true
		}
	}

	// knights
	for _, n := range dirKnight { // has the offsets for the knight jumps
		kSq := sq + n
		if attackingColor == WHITE && p.pieces[kSq] == wN {
			return true
		} else if attackingColor == BLACK && p.pieces[kSq] == bN {
			return true
		}
	}

	// rook/queen
	for _, dir := range dirRook {
		tsq := sq + dir
		pce := p.pieces[tsq]
		// only go until we're off the board
		for pce != NO_SQ {
			if pce != EMPTY {
				if pieceLookups[pce].isRookOrQueen && pieceLookups[pce].color == attackingColor {
					return true
				}
				break // break out if we hit something that wasn't a rook or queen
			}
			tsq += dir // move to the next square
			pce = p.pieces[tsq]
		}
	}

	// bishop/queen
	for _, dir := range dirBishop {
		tsq := sq + dir
		pce := p.pieces[tsq]
		// only go until we're off the board
		for pce != NO_SQ {
			if pce != EMPTY {
				if pieceLookups[pce].isBishopOrQueen && pieceLookups[pce].color == attackingColor {
					return true
				}
				break // break out if we hit something else
			}
			tsq += dir // move to the next square
			pce = p.pieces[tsq]
		}
	}

	for _, dir := range dirKing {
		tsq := sq + dir
		p := p.pieces[tsq]
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

func (p Position) String() string {
	output := strings.Builder{}
	output.WriteString(p.PrintBoard())
	output.WriteString("--------\n")
	output.WriteString(fmt.Sprintf("side: %v\n", p.side))
	output.WriteString(fmt.Sprintf("enPas: %v\n", p.enPas))
	output.WriteString(fmt.Sprintf("castle: %v\n", p.castlePerm))
	output.WriteString(fmt.Sprintf("posKey: %x\n", p.posKey))

	return output.String()
}
func (p Position) PrintBoard() string {
	output := strings.Builder{}
	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			sq := fileRankToSq(int(f), int(r))
			piece := p.pieces[sq]
			output.WriteString(piece.String())
		}
		output.WriteString("\n")
	}
	return output.String()
}

// PrintAttackBoard returns a string representation
// of all squares attacked by a given side
func (p Position) PrintAttackBoard(attackingSide int) string {
	output := strings.Builder{}

	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			sq := fileRankToSq(f, r)
			if p.IsSquareAttacked(sq, attackingSide) {
				output.WriteString("X")
			} else {
				output.WriteString("-")
			}
		}
		output.WriteString("\n")
	}
	return output.String()
}
