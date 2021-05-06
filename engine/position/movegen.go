package position

import (
	"fmt"
	"strings"
)

type Movescore struct {
	Key   movekey
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
func (list *Movelist) addQuietMove(move movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addCaptureMove(move movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addEnPasMove(move movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addWhitePawnCaptureMove(from, to int, captured piece) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(wQ))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(wR))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(wB))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(wN))
	} else {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured))

	}
}

func (list *Movelist) addWhitePawnMove(from, to int) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(wQ))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(wR))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(wB))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(wN))
	} else {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to))
	}
}

func (list *Movelist) addBlackPawnCaptureMove(from, to int, captured piece) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(bQ))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(bR))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(bB))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(bN))
	} else {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setCaptured(captured))
	}
}

func (list *Movelist) addBlackPawnMove(from, to int) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(bQ))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(bR))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(bB))
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to).setPromoted(bN))
	} else {
		list.addCaptureMove(movekey(0).setFrom(from).setTo(to))
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
					list.addQuietMove(movekey(0).setFrom(E1).setTo(G1).setCastle())
				}
			}
		}

		// white queenside
		if p.castlePerm.Has(CASTLE_PERMS_WQ) {
			if p.pieces[D1] == EMPTY && p.pieces[C1] == EMPTY && p.pieces[B1] == EMPTY {
				if !p.IsSquareAttacked(E1, BLACK) && !p.IsSquareAttacked(D1, BLACK) && !p.IsSquareAttacked(C1, BLACK) {
					list.addQuietMove(movekey(0).setFrom(E1).setTo(C1).setCastle())
				}
			}
		}

	} else {
		// black kingside
		if p.castlePerm.Has(CASTLE_PERMS_BK) {
			if p.pieces[F8] == EMPTY && p.pieces[G8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(F8, WHITE) && !p.IsSquareAttacked(G8, WHITE) {
					// add a new castle move
					list.addQuietMove(movekey(0).setFrom(E8).setTo(G8).setCastle())
				}
			}
		}

		// black queenside
		if p.castlePerm.Has(CASTLE_PERMS_BQ) {
			if p.pieces[D8] == EMPTY && p.pieces[C8] == EMPTY && p.pieces[B8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(D8, WHITE) && !p.IsSquareAttacked(C8, WHITE) {
					list.addQuietMove(movekey(0).setFrom(E8).setTo(C8).setCastle())
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
				list.addWhitePawnMove(sq, sq+10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_2 && p.pieces[sq+20] == EMPTY {
					list.addQuietMove(movekey(0).setFrom(sq).setTo(sq + 20).setPawnStart())
				}
			}

			// check pawn captures
			lSq := sq + 9
			rSq := sq + 11
			if !sqOffBoard(lSq) && pieceLookups[p.pieces[lSq]].color == BLACK {
				list.addWhitePawnCaptureMove(sq, lSq, p.pieces[lSq])
			}
			if !sqOffBoard(rSq) && pieceLookups[p.pieces[rSq]].color == BLACK {
				list.addWhitePawnCaptureMove(sq, rSq, p.pieces[rSq])
			}

			// enPas captures
			if lSq == p.enPas {
				list.addCaptureMove(movekey(0).setFrom(sq).setTo(lSq).setEnPas())
			}
			if rSq == p.enPas {
				list.addCaptureMove(movekey(0).setFrom(sq).setTo(rSq).setEnPas())
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
				list.addBlackPawnMove(sq, sq-10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_7 && p.pieces[sq-20] == EMPTY {
					list.addQuietMove(movekey(0).setFrom(sq).setTo(sq - 20).setPawnStart())
				}
			}

			// check pawn captures
			lSq := sq - 9
			rSq := sq - 11
			if !sqOffBoard(lSq) && pieceLookups[p.pieces[lSq]].color == WHITE {
				list.addBlackPawnCaptureMove(sq, lSq, p.pieces[lSq])
			}
			if !sqOffBoard(rSq) && pieceLookups[p.pieces[rSq]].color == WHITE {
				list.addBlackPawnCaptureMove(sq, rSq, p.pieces[rSq])
			}

			// enPas captures
			if lSq == p.enPas {
				list.addCaptureMove(movekey(0).setFrom(sq).setTo(lSq).setEnPas())
			}
			if rSq == p.enPas {
				list.addCaptureMove(movekey(0).setFrom(sq).setTo(rSq).setEnPas())
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
					list.addQuietMove(movekey(0).setFrom(fromSq).setTo(toSq))
				} else if pieceLookups[p.pieces[toSq]].color != p.side {
					list.addCaptureMove(movekey(0).setFrom(fromSq).setTo(toSq).setCaptured(p.pieces[toSq]))
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
						list.addQuietMove(movekey(0).setFrom(fromSq).setTo(toSq))
					} else {
						if pieceLookups[toPce].color != p.side {
							list.addCaptureMove(movekey(0).setFrom(fromSq).setTo(toSq).setCaptured(toPce))
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
