package position

import (
	"fmt"
	"strings"
)

type Movescore struct {
	Key   Movekey
	Score int
}

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

// todo - clean these up? Might use them to speed up search?
func (list *Movelist) addQuietMove(move Movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addCaptureMove(move Movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addEnPasMove(move Movekey) {
	*list = append(*list, &Movescore{move, 0})
}

func (list *Movelist) addWhitePawnCaptureMove(from, to int, captured Piece) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PwQ))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PwR))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PwB))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PwN))
	} else {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured))

	}
}

func (list *Movelist) addWhitePawnMove(from, to int) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PwQ))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PwR))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PwB))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PwN))
	} else {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to))
	}
}

func (list *Movelist) addBlackPawnCaptureMove(from, to int, captured Piece) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PbQ))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PbR))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PbB))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured).setPromoted(PbN))
	} else {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setCaptured(captured))
	}
}

func (list *Movelist) addBlackPawnMove(from, to int) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PbQ))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PbR))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PbB))
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to).setPromoted(PbN))
	} else {
		list.addCaptureMove(Movekey(0).setFrom(from).setTo(to))
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
					list.addQuietMove(Movekey(0).setFrom(E1).setTo(G1).setCastle())
				}
			}
		}

		// white queenside
		if p.castlePerm.Has(CASTLE_PERMS_WQ) {
			if p.pieces[D1] == EMPTY && p.pieces[C1] == EMPTY && p.pieces[B1] == EMPTY {
				if !p.IsSquareAttacked(E1, BLACK) && !p.IsSquareAttacked(D1, BLACK) && !p.IsSquareAttacked(C1, BLACK) {
					list.addQuietMove(Movekey(0).setFrom(E1).setTo(C1).setCastle())
				}
			}
		}

	} else {
		// black kingside
		if p.castlePerm.Has(CASTLE_PERMS_BK) {
			if p.pieces[F8] == EMPTY && p.pieces[G8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(F8, WHITE) && !p.IsSquareAttacked(G8, WHITE) {
					// add a new castle move
					list.addQuietMove(Movekey(0).setFrom(E8).setTo(G8).setCastle())
				}
			}
		}

		// black queenside
		if p.castlePerm.Has(CASTLE_PERMS_BQ) {
			if p.pieces[D8] == EMPTY && p.pieces[C8] == EMPTY && p.pieces[B8] == EMPTY {
				if !p.IsSquareAttacked(E8, WHITE) && !p.IsSquareAttacked(D8, WHITE) && !p.IsSquareAttacked(C8, WHITE) {
					list.addQuietMove(Movekey(0).setFrom(E8).setTo(C8).setCastle())
				}
			}
		}
	}

	// pawns
	if p.side == WHITE {
		// iterate through each pawn
		for i := 0; i < p.pieceCount[PwP]; i++ {
			sq := p.pieceList[PwP][i]
			// todo - comment out to speed up
			if sqOffBoard(sq) {
				panic(fmt.Errorf("PwP position off board: %v", sq))
			}

			// check pawn movements
			if p.pieces[sq+10] == EMPTY {
				list.addWhitePawnMove(sq, sq+10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_2 && p.pieces[sq+20] == EMPTY {
					list.addQuietMove(Movekey(0).setFrom(sq).setTo(sq + 20).setPawnStart())
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
			if p.enPas != NO_SQ {
				if lSq == p.enPas {
					list.addCaptureMove(Movekey(0).setFrom(sq).setTo(lSq).setEnPas())
				}
				if rSq == p.enPas {
					list.addCaptureMove(Movekey(0).setFrom(sq).setTo(rSq).setEnPas())
				}
			}
		}
	} else { // black side
		// iterate through each pawn
		for i := 0; i < p.pieceCount[PbP]; i++ {
			sq := p.pieceList[PbP][i]
			// todo - comment out to speed up
			if sqOffBoard(sq) {
				panic(fmt.Errorf("PbP position off board: %v", sq))
			}

			// check pawn movements
			if p.pieces[sq-10] == EMPTY {
				list.addBlackPawnMove(sq, sq-10)

				// check for RANK_2 double move
				if rankLookups[sq] == RANK_7 && p.pieces[sq-20] == EMPTY {
					list.addQuietMove(Movekey(0).setFrom(sq).setTo(sq - 20).setPawnStart())
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
			if p.enPas != NO_SQ {
				if lSq == p.enPas {
					list.addCaptureMove(Movekey(0).setFrom(sq).setTo(lSq).setEnPas())
				}
				if rSq == p.enPas {
					list.addCaptureMove(Movekey(0).setFrom(sq).setTo(rSq).setEnPas())
				}
			}
		}
	}

	// set up all pieces per color
	var nonslidingPieces [2]Piece
	var slidingPieces [3]Piece
	if p.side == WHITE {
		nonslidingPieces = [2]Piece{PwN, PwK}
		slidingPieces = [3]Piece{PwB, PwR, PwQ}
	} else {
		nonslidingPieces = [2]Piece{PbN, PbK}
		slidingPieces = [3]Piece{PbB, PbR, PbQ}
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
					list.addQuietMove(Movekey(0).setFrom(fromSq).setTo(toSq))
				} else if pieceLookups[p.pieces[toSq]].color != p.side {
					list.addCaptureMove(Movekey(0).setFrom(fromSq).setTo(toSq).setCaptured(p.pieces[toSq]))
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
						list.addQuietMove(Movekey(0).setFrom(fromSq).setTo(toSq))
					} else {
						if pieceLookups[toPce].color != p.side {
							list.addCaptureMove(Movekey(0).setFrom(fromSq).setTo(toSq).setCaptured(toPce))
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
