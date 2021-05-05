package position

import (
	"fmt"
)

/*
Making a move...

1. make move
2. get the from, to, cap from the move
3. store the current position in the p.history array
4. move the current piece from -> to
5. if a capture was made, remove it from the piece list
6. update 50 move rule based on pawn move
7. promotions
8. enPas captures
9. set enPas square if needed
10. for all pieces added, moved, removed, etc update position counters + piece lists
11. maintain hashkey
12. castle permissions
13. change side, increment ply + hisPly
*/

// hash keys?

// clear a single sq from the position
func (p *Position) clearPiece(sq int) {
	if sqOffBoard(sq) {
		panic(fmt.Sprintf("sq %v is off board", sq))
	}

	pce := p.pieces[sq]
	pceMeta := pieceLookups[pce]

	// todo - remove
	if pce == EMPTY {
		panic(fmt.Sprintf("tried to clean an empty piece on sq %d", sq))
	}

	// set square, subtract value
	p.pieces[sq] = EMPTY
	p.materialCount[pceMeta.color] -= pceMeta.value

	// update big/major/minor/pawns
	if pceMeta.isBig {
		p.bigPieceCount[pceMeta.color]--
		if pceMeta.isMajor {
			p.majPieceCount[pceMeta.color]--
		}
		if pceMeta.isMinor {
			p.minPieceCount[pceMeta.color]--
		}
	} else {
		p.pawns[pceMeta.color].clear(SQ64(sq))
		p.pawns[BOTH].clear(SQ64(sq))
	}

	// find where it is in the pieceList
	foundPieceIndex := -1
	for i := 0; i < p.pieceCount[pce]; i++ {
		if p.pieceList[pce][i] == sq {
			foundPieceIndex = i
			break
		}
	}

	// copy the last item in the list to the found index
	p.pieceList[pce][foundPieceIndex] = p.pieceList[pce][p.pieceCount[pce]-1]
	// delete the last index, since we copied it forward
	p.pieceList[pce][p.pieceCount[pce]-1] = 0
	// decrement the total piece Count to match
	p.pieceCount[pce]--
}

func (p *Position) addPiece(sq int, pce piece) {

	pceMeta := pieceLookups[pce]

	p.pieces[sq] = pce

	if pceMeta.isBig {
		p.bigPieceCount[pceMeta.color]++
		if pceMeta.isMajor {
			p.majPieceCount[pceMeta.color]++
		} else if pceMeta.isMinor {
			p.minPieceCount[pceMeta.color]++
		}
	} else {
		p.pawns[pceMeta.color].set(SQ64(sq))
		p.pawns[BOTH].set(SQ64(sq))
	}

	// add value
	p.materialCount[pceMeta.color] += pceMeta.value

	// update pieceLists
	p.pieceList[pce][p.pieceCount[pce]] = sq
	p.pieceCount[pce]++
}

func (p *Position) movePiece(from, to int) {
	pce := p.pieces[from]
	pceMeta := pieceLookups[pce]

	if pce == EMPTY {
		panic("woops!")
	}

	p.pieces[from] = EMPTY
	p.pieces[to] = pce

	// update the pawn boards as needed
	if !pceMeta.isBig {
		p.pawns[pceMeta.color].clear(SQ64(from))
		p.pawns[BOTH].clear(SQ64(from))
		p.pawns[pceMeta.color].set(SQ64(to))
		p.pawns[BOTH].set(SQ64(to))
	}

	found := false
	for i := 0; i < p.pieceCount[pce]; i++ {
		if p.pieceList[pce][i] == from {
			p.pieceList[pce][i] = to
			found = true
			break
		}
	}

	// todo - eliminate after perft
	if !found {
		fmt.Println("woops!")
		// something setting p.pieceList[bP][1] to 0
		panic("didnt find existing piece")
	}
}

// MakeMove updates the position for a newly made
// move. It returns false if a king is left in check.
func (p *Position) MakeMove(move *movekey) bool {
	err := p.AssertCache()
	if err != nil {
		panic(err)
	}

	from := move.getFrom()
	to := move.getTo()
	side := p.side
	captured := move.getCaptured()
	promoted := move.getPromoted()
	pce := p.pieces[from]

	p.history[p.hisPly].posKey = p.posKey

	// enPas need to remove an additional piece
	if move.isEnPas() {
		if side == WHITE {
			p.clearPiece(to - 10)
		} else {
			p.clearPiece(to + 10)
		}
	}

	// castling
	if move.isCastle() {
		switch to {
		case C1: // white queenside
			p.movePiece(A1, D1)
		case G1: // white kingside
			p.movePiece(H1, F1)
		case C8: // black queenside
			p.movePiece(A8, D8)
		case G8: // black kingside
			p.movePiece(H8, F8)
		default:
			panic(fmt.Sprintf("castle to sq not recognized: %v", to))
		}
	}

	// hash enPas?

	p.history[p.hisPly].move = *move
	p.history[p.hisPly].fiftyMove = p.fiftyMove
	p.history[p.hisPly].enPas = p.enPas
	p.history[p.hisPly].castlePerm = *p.castlePerm

	// hash castle out?

	// update castlePerms
	// todo - potentially speed this up
	switch from {
	case A1:
		p.castlePerm.Clear(CASTLE_PERMS_WQ)
	case E1:
		p.castlePerm.Clear(CASTLE_PERMS_WK)
		p.castlePerm.Clear(CASTLE_PERMS_WQ)
	case H1:
		p.castlePerm.Clear(CASTLE_PERMS_WK)
	case A8:
		p.castlePerm.Clear(CASTLE_PERMS_BQ)
	case E8:
		p.castlePerm.Clear(CASTLE_PERMS_BQ)
		p.castlePerm.Clear(CASTLE_PERMS_BK)
	case H8:
		p.castlePerm.Clear(CASTLE_PERMS_BK)
	}

	if move.isPawnStart() {
		if side == WHITE {
			p.enPas = to - 10
		} else {
			p.enPas = to + 10
		}
	} else {
		p.enPas = NO_SQ
	}

	// hash castle?

	// update history in general
	p.fiftyMove++
	p.hisPly++
	p.searchPly++

	// If it was a pawn move, reset to 0
	if pce == wP || pce == bP {
		p.fiftyMove = 0
	}

	// capture piece
	if captured != EMPTY {
		p.clearPiece(to)
		p.fiftyMove = 0

		// check castle perms on capture for rooks
		// todo - optimize
		switch to {
		case A1:
			p.castlePerm.Clear(CASTLE_PERMS_WQ)
		case H1:
			p.castlePerm.Clear(CASTLE_PERMS_WK)
		case A8:
			p.castlePerm.Clear(CASTLE_PERMS_BQ)
		case H8:
			p.castlePerm.Clear(CASTLE_PERMS_BK)
		}
	}

	// move piece very last, after clearing capture
	p.movePiece(from, to)

	// promoted
	if promoted != EMPTY {
		p.clearPiece(to)
		p.addPiece(to, promoted)
	}

	// update any king move
	if pce == wK || pce == bK {
		p.kingSq[p.side] = to
	}

	// update side
	p.side ^= 1 // flip from 0 <-> 1

	// todo - optimize hash
	p.posKey = p.GenPosKey()

	// assert we're set up right
	err = p.AssertCache()
	if err != nil {
		fmt.Errorf("woops!")
		panic(err)
	}

	// last check if king is now attacked
	if p.IsSquareAttacked(p.kingSq[side], p.side) {
		p.UndoMove()
		return false
	}

	if p.enPas == 55 {
		fmt.Println("huh")
		panic("uh oh")
	}
	return true
}

func (p *Position) UndoMove() {
	err := p.AssertCache()
	if err != nil {
		panic(err)
	}

	p.hisPly--
	p.searchPly--

	u := p.history[p.hisPly]
	move := u.move
	from := move.getFrom()
	to := move.getTo()
	captured := move.getCaptured()

	p.castlePerm = &u.castlePerm
	p.fiftyMove = u.fiftyMove
	p.enPas = u.enPas

	// switch sides
	p.side ^= 1

	// enPas need to remove an additional piece
	if move.isEnPas() {
		if p.side == WHITE {
			p.addPiece(to-10, bP)
		} else {
			p.addPiece(to+10, wP)
		}
	}

	// castling
	if move.isCastle() {
		switch to {
		case C1: // white queenside
			p.movePiece(D1, A1)
		case G1: // white kingside
			p.movePiece(F1, H1)
		case C8: // black queenside
			p.movePiece(D8, A8)
		case G8: // black kingside
			p.movePiece(F8, H8)
		default:
			panic(fmt.Sprintf("castle to sq not recognized: %v", to))
		}
	}

	// promoted
	if move.getPromoted() != EMPTY {
		p.clearPiece(to)
		if pieceLookups[move.getPromoted()].color == WHITE {
			p.addPiece(to, wP)
		} else {
			p.addPiece(to, bP)
		}
	}

	p.movePiece(to, from)

	// restore king lookup if needed
	pce := p.pieces[from]
	if pce == wK || pce == bK {
		p.kingSq[p.side] = from
	}

	// restore capture
	if captured != EMPTY {
		p.addPiece(to, captured)
	}

	// rehash
	p.posKey = p.GenPosKey()

	err = p.AssertCache()
	if err != nil {
		fmt.Println("woops")
		panic(err)
	}
}
