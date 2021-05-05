package position

import (
	"fmt"
	"strings"
)

type Movescore struct {
	Key   *movekey
	Score int
}

// todo - turn this into a slice
type Movelist []*Movescore

func (list *Movelist) String() string {
	b := strings.Builder{}

	b.WriteString("movelist: \n")
	for i, mv := range *list {
		move := mv.Key
		score := mv.Score
		b.WriteString(fmt.Sprintf("Move:%d > %v (Score:%d)\n", i, move.ShortString(), score))
	}

	b.WriteString(fmt.Sprintf("movelist total: %d", len(*list)))

	return b.String()
}

// todo - clean these up?
func (list *Movelist) addQuietMove(move *movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addCaptureMove(move *movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addEnPasMove(move *movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addWhitePawnCaptureMove(p *Position, from, to int, captured piece) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(newMovekey(from, to, captured, wQ, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, wR, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, wB, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, wN, false, false))
	} else {
		list.addCaptureMove(newMovekey(from, to, captured, EMPTY, false, false))
	}
}

func (list *Movelist) addWhitePawnMove(p *Position, from, to int) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(newMovekey(from, to, EMPTY, wQ, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, wR, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, wB, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, wN, false, false))
	} else {
		list.addCaptureMove(newMovekey(from, to, EMPTY, EMPTY, false, false))
	}
}

func (list *Movelist) addBlackPawnCaptureMove(p *Position, from, to int, captured piece) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(newMovekey(from, to, captured, bQ, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, bR, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, bB, false, false))
		list.addCaptureMove(newMovekey(from, to, captured, bN, false, false))
	} else {
		list.addCaptureMove(newMovekey(from, to, captured, EMPTY, false, false))
	}
}

func (list *Movelist) addBlackPawnMove(p *Position, from, to int) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(newMovekey(from, to, EMPTY, bQ, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, bR, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, bB, false, false))
		list.addCaptureMove(newMovekey(from, to, EMPTY, bN, false, false))
	} else {
		list.addCaptureMove(newMovekey(from, to, EMPTY, EMPTY, false, false))
	}
}

func (p *Position) GenerateAllMoves() *Movelist {

	list := &Movelist{}

	// castling
	if p == nil {
		panic(fmt.Errorf("p should not be nil"))
	}

	if p.side == WHITE {

		// white kingside
		if p.castlePerm.Has(CASTLE_PERMS_WK) {
			if p.pieces[F1] == EMPTY && p.pieces[G1] == EMPTY {
				if !p.IsSquareAttacked(E1, BLACK) && !p.IsSquareAttacked(F1, BLACK) && !p.IsSquareAttacked(G1, BLACK) {
					// add a new castle move
					list.addQuietMove(newCastleMoveKey(E1, G1))
				}
			}
		}

		// white queenside
		if p.castlePerm.Has(CASTLE_PERMS_WQ) {
			if p.pieces[D1] == EMPTY && p.pieces[C1] == EMPTY && p.pieces[B1] == EMPTY {
				if !p.IsSquareAttacked(E1, BLACK) && !p.IsSquareAttacked(D1, BLACK) && !p.IsSquareAttacked(C1, BLACK) {
					list.addQuietMove(newCastleMoveKey(E1, C1))
				}
			}
		}

	} else {
		// black kingside
		if p.castlePerm.Has(CASTLE_PERMS_BK) {
			if p.pieces[F8] == EMPTY && p.pieces[G8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(F8, WHITE) && !p.IsSquareAttacked(G8, WHITE) {
					// add a new castle move
					list.addQuietMove(newCastleMoveKey(E8, G8))
				}
			}
		}

		// black queenside
		if p.castlePerm.Has(CASTLE_PERMS_BQ) {
			if p.pieces[D8] == EMPTY && p.pieces[C8] == EMPTY && p.pieces[B8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(D8, WHITE) && !p.IsSquareAttacked(C8, WHITE) {
					list.addQuietMove(newCastleMoveKey(E8, C8))
				}
			}
		}
	}

	// pawns
	if p.side == WHITE {
		// iterate through each pawn
		for i := 0; i < p.pieceCount[wP]; i++ {
			sq := p.pieceList[wP][i]
			// todo - comment out to speed up
			if sqOffBoard(sq) {
				panic(fmt.Errorf("wP position off board: %v", sq))
			}

			// check pawn movements
			if p.pieces[sq+10] == EMPTY {
				list.addWhitePawnMove(p, sq, sq+10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_2 && p.pieces[sq+20] == EMPTY {
					list.addQuietMove(newMovekey(sq, sq+20, EMPTY, EMPTY, false, true))
				}
			}

			// check pawn captures
			lSq := sq + 9
			rSq := sq + 11
			if !sqOffBoard(lSq) && pieceLookups[p.pieces[lSq]].color == BLACK {
				list.addWhitePawnCaptureMove(p, sq, lSq, p.pieces[lSq])
			}
			if !sqOffBoard(rSq) && pieceLookups[p.pieces[rSq]].color == BLACK {
				list.addWhitePawnCaptureMove(p, sq, rSq, p.pieces[rSq])
			}

			// enPas captures
			if lSq == p.enPas {
				list.addCaptureMove(newMovekey(sq, lSq, EMPTY, EMPTY, true, false))
			}
			if rSq == p.enPas {
				list.addCaptureMove(newMovekey(sq, rSq, EMPTY, EMPTY, true, false))
			}
		}
	} else { // black side
		// iterate through each pawn
		for i := 0; i < p.pieceCount[bP]; i++ {
			sq := p.pieceList[bP][i]
			// todo - comment out to speed up
			if sqOffBoard(sq) {
				panic(fmt.Errorf("bP position off board: %v", sq))
			}

			// check pawn movements
			if p.pieces[sq-10] == EMPTY {
				list.addBlackPawnMove(p, sq, sq-10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_7 && p.pieces[sq-20] == EMPTY {
					list.addQuietMove(newMovekey(sq, sq-20, EMPTY, EMPTY, false, true))
				}
			}

			// check pawn captures
			lSq := sq - 9
			rSq := sq - 11
			if !sqOffBoard(lSq) && pieceLookups[p.pieces[lSq]].color == WHITE {
				list.addBlackPawnCaptureMove(p, sq, lSq, p.pieces[lSq])
			}
			if !sqOffBoard(rSq) && pieceLookups[p.pieces[rSq]].color == WHITE {
				list.addBlackPawnCaptureMove(p, sq, rSq, p.pieces[rSq])
			}

			// enPas captures
			if lSq == p.enPas {
				list.addCaptureMove(newMovekey(sq, lSq, EMPTY, EMPTY, true, false))
			}
			if rSq == p.enPas {
				list.addCaptureMove(newMovekey(sq, rSq, EMPTY, EMPTY, true, false))
			}
		}
	}

	// set up all pieces per color
	var nonslidingPieces [2]piece
	var slidingPieces [3]piece
	if p.side == WHITE {
		nonslidingPieces = [2]piece{wN, wK}
		slidingPieces = [3]piece{wB, wR, wQ}
	} else {
		nonslidingPieces = [2]piece{bN, bK}
		slidingPieces = [3]piece{bB, bR, bQ}
	}

	// knights/kings
	for _, pce := range nonslidingPieces {
		for pceNum := 0; pceNum < p.pieceCount[pce]; pceNum++ {
			fromSq := p.pieceList[pce][pceNum]

			for _, dir := range pieceLookups[pce].dir {
				toSq := fromSq + dir
				if sqOffBoard(toSq) {
					continue
				}
				if p.pieces[toSq] == EMPTY {
					list.addQuietMove(newMovekey(fromSq, toSq, EMPTY, EMPTY, false, false))
				} else if pieceLookups[p.pieces[toSq]].color != p.side {
					list.addCaptureMove(newMovekey(fromSq, toSq, p.pieces[toSq], EMPTY, false, false))
				}
			}
		}
	}

	// sliding pieces
	for _, pce := range slidingPieces {
		for pceNum := 0; pceNum < p.pieceCount[pce]; pceNum++ {
			fromSq := p.pieceList[pce][pceNum]

			for _, dir := range pieceLookups[pce].dir {
				toSq := fromSq + dir
				toPce := p.pieces[toSq]
				for toPce != NO_SQ {
					if toPce == EMPTY {
						list.addQuietMove(newMovekey(fromSq, toSq, EMPTY, EMPTY, false, false))
					} else {
						if pieceLookups[toPce].color != p.side {
							list.addCaptureMove(newMovekey(fromSq, toSq, toPce, EMPTY, false, false))
						}
						break // go to next direction
					}
					toSq += dir
					toPce = p.pieces[toSq]
				}
			}
		}
	}

	return list
}
