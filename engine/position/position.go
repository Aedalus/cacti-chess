package position

import (
	"fmt"
	"reflect"
	"strings"
)

// undo tracks all data needed to revert the board to the previous state
type undo struct {
	move       Movekey
	castlePerm castlePerm
	enPas      int
	fiftyMove  int
	posKey     uint64
}

// Position is a given state of the board. The pieces field is the source of truth,
// with other fields allowing for quick lookups via caching
type Position struct {
	pieces *board120 // Source of truth for pieces
	side   int       // white/black

	pawns [3]*bitboard64 // white/black/both quick pawn lookups

	// Piece list for fast lookup
	// pieceList[PwN][0] = E1 etc
	pieceList [13][10]int

	kingSq [2]int // king quick lookups

	castlePerm *castlePerm // permissions to castle

	enPas int // if en passant is available

	fiftyMove int // 50 move counter (100 since we're using half Moves)

	searchPly int // how far we are in the search

	posKey uint64 // unique index for position

	pieceCount    [13]int // Count of all pieces on the board
	bigPieceCount [2]int  // white/black non-pawn pieces
	majPieceCount [2]int  // queen/rook
	minPieceCount [2]int  // bishop/knight
	materialCount [2]int  // Count of material

	// history
	hisPly  int // how many half Moves have been made in the whole game
	history []undo
}

func (p *Position) GetPosKey() uint64 {
	return p.posKey
}

func (p *Position) GetSearchPly() int {
	return p.searchPly
}

func (p *Position) GetMaterial() [2]int {
	return p.materialCount
}

func (p *Position) GetPieceCount() [13]int {
	return p.pieceCount
}

func (p *Position) GetPieceList() [13][10]int {
	return p.pieceList
}

func (p *Position) GetFiftyMove() int {
	return p.fiftyMove
}

func (p *Position) GetSide() int {
	return p.side
}

// GenPosKey generates a statistically unique uint64
// for the current state of the position
func (p Position) GenPosKey() uint64 {
	var finalKey uint64 = 0
	var pce Piece = Piece(0)

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

// updateListCaches updates all piece caches based on pieces
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
			// [PwP][0] = A2
			// [PwP][1] = B2 etc
			curPieceCount := p.pieceCount[pce]
			p.pieceList[pce][curPieceCount] = i
			p.pieceCount[pce]++

			// update king positions
			if pce == PwK {
				p.kingSq[WHITE] = i
			}
			if pce == PbK {
				p.kingSq[BLACK] = i
			}

			// update pawn boards
			if pce == PwP {
				p.pawns[WHITE].set(SQ64(i))
				p.pawns[BOTH].set(SQ64(i))
			}
			if pce == PbP {
				p.pawns[BLACK].set(SQ64(i))
				p.pawns[BOTH].set(SQ64(i))
			}
		}
	}
}

// AssertCache was used heavily during initial development,
// but slows down computation by ~2x. It recalculates the cache
// from scratch and asserts the existing cached values are as
// they should be
func (p *Position) AssertCache() error {
	// Comment out if debugging
	return nil

	// temporary values we recompute to check against
	t_pieceCount := [13]int{}
	t_pieceList := [13][10]int{}
	t_bigPieceCount := [2]int{}
	t_majPieceCount := [2]int{}
	t_minPieceCount := [2]int{}
	t_kingSq := [2]int{}
	t_materialCount := [2]int{}
	t_posKey := p.GenPosKey()

	t_pawns := [3]*bitboard64{
		&bitboard64{},
		&bitboard64{},
		&bitboard64{},
	}

	// check Piece lists
	for pce := PwP; pce <= PbK; pce++ {
		for i := 0; i < p.pieceCount[pce]; i++ {
			sq := p.pieceList[pce][i]
			if sq == 0 {
				return fmt.Errorf("pieceList error: (%d)[%d]: sq should not be 0", int(pce), i)
			}
			if sqOffBoard(sq) {
				return fmt.Errorf("pieceList error: (%d)[%d]: sq %d is off board", int(pce), i, sq)
			}
			if p.pieces[sq] != pce {
				return fmt.Errorf("pieceList error: (%d)[%d]: sq %d had %d", int(pce), i, sq, p.pieces[sq])
			}

		}
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
			// [PwP][0] = A2
			// [PwP][1] = B2 etc
			curPieceCount := t_pieceCount[pce]
			t_pieceList[pce][curPieceCount] = i
			t_pieceCount[pce]++

			// update king positions
			if pce == PwK {
				t_kingSq[WHITE] = i
			}
			if pce == PbK {
				t_kingSq[BLACK] = i
			}

			// update pawn boards
			if pce == PwP {
				t_pawns[WHITE].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
			if pce == PbP {
				t_pawns[BLACK].set(SQ64(i))
				t_pawns[BOTH].set(SQ64(i))
			}
		}
	}

	if !reflect.DeepEqual(t_pieceCount, p.pieceCount) {
		return fmt.Errorf("pieceCount - got %v, want %v", p.pieceCount, t_pieceCount)
	}
	if !reflect.DeepEqual(t_bigPieceCount, p.bigPieceCount) {
		return fmt.Errorf("bigPieceCount - got %v, want %v", p.bigPieceCount, t_bigPieceCount)
	}
	if !reflect.DeepEqual(t_majPieceCount, p.majPieceCount) {
		return fmt.Errorf("majPieceCount - got %v want %v", p.majPieceCount, t_majPieceCount)
	}
	if !reflect.DeepEqual(t_minPieceCount, p.minPieceCount) {
		return fmt.Errorf("minPieceCount - got %v want %v", p.minPieceCount, t_minPieceCount)
	}
	if !reflect.DeepEqual(t_kingSq, p.kingSq) {
		return fmt.Errorf("kingSq - got %v want %v", p.kingSq, t_kingSq)
	}
	if !reflect.DeepEqual(t_materialCount, p.materialCount) {
		return fmt.Errorf("materialCount - got %v want %v", p.materialCount, t_materialCount)
	}
	if !reflect.DeepEqual(t_pawns, p.pawns) {
		return fmt.Errorf("pawns - got %v want %v", p.pawns, t_pawns)
	}
	if p.posKey != t_posKey {
		return fmt.Errorf("posKey - got %v want %v", p.posKey, t_posKey)
	}

	return nil
}

// Reset will re-initialize the engine to an empty state
func (p *Position) Reset() {
	// initialize Piece array
	p.pieces = &board120{}
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		p.pieces[i] = NO_SQ
	}
	for i := 0; i < 64; i++ {
		p.pieces[SQ120(i)] = EMPTY
	}

	// castle perms
	p.castlePerm = &castlePerm{CASTLE_PERMS_NONE}

	// Piece counts
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
	p.hisPly = 0
	p.history = []undo{}
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
	err := p.AssertCache()
	if err != nil {
		panic(err)
	}

	// pawns
	if attackingColor == WHITE {
		if p.pieces[sq-11] == PwP || p.pieces[sq-9] == PwP {
			return true
		}
	} else {
		if p.pieces[sq+11] == PbP || p.pieces[sq+9] == PbP {
			return true
		}
	}

	// knights
	for _, n := range dirKnight { // has the offsets for the knight jumps
		kSq := sq + n
		if attackingColor == WHITE && p.pieces[kSq] == PwN {
			return true
		} else if attackingColor == BLACK && p.pieces[kSq] == PbN {
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
		if p == PwK && attackingColor == WHITE {
			return true
		}
		if p == PbK && attackingColor == BLACK {
			return true
		}
	}

	// only return false once we've ruled out everything
	return false
}

func (p *Position) IsKingAttacked() bool {
	return p.IsSquareAttacked(p.kingSq[p.side], p.side^1)
}

func (p *Position) IsLegalMove() bool {
	mlist := p.GenerateAllMoves()

	legal := 0

	for _, m := range *mlist {
		if valid := p.MakeMove(m.Key); valid == true {
			legal++
			p.UndoMove()
		}
	}

	return legal > 0
}

func (p *Position) IsStalemate() bool {
	kingAttacked := p.IsKingAttacked()
	legalMove := p.IsLegalMove()
	return !kingAttacked && !legalMove
}

func (p *Position) IsCheckmate() bool {
	kingAttacked := p.IsKingAttacked()
	legalMove := p.IsLegalMove()

	return kingAttacked && !legalMove
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

func (p *Position) IsRepetition() bool {
	for _, his := range p.history {
		if his.posKey == p.posKey {
			return true
		}
	}
	return false
}
