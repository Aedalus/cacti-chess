package engine

import (
	"fmt"
	"strings"
)

type movescore struct {
	key   *movekey
	score int
}

type movelist struct {
	moves [256]movescore
	count int
}

func (list *movelist) String() string {
	b := strings.Builder{}

	b.WriteString("movelist: \n")
	for i := 0; i < list.count; i++ {
		move := list.moves[i].key
		score := list.moves[i].score

		b.WriteString(fmt.Sprintf("Move:%d > %v (score:%d)\n", i, move.ShortString(), score))
	}
	b.WriteString(fmt.Sprintf("movelist total: %d", list.count))

	return b.String()
}

func (list *movelist) addQuietMove(p *Position, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func (list *movelist) addCaptureMove(p *Position, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func (list *movelist) addEnPasMove(p *Position, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func (list *movelist) addWhitePawnCaptureMove(p *Position, from, to int, captured piece) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(p, newMovekey(from, to, captured, wQ, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, wR, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, wB, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, wN, false, false))
	} else {
		list.addCaptureMove(p, newMovekey(from, to, captured, EMPTY, false, false))
	}
}

func (list *movelist) addWhitePawnMove(p *Position, from, to int) {
	if rankLookups[from] == RANK_7 {
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, wQ, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, wR, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, wB, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, wN, false, false))
	} else {
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, EMPTY, false, false))
	}
}

func (list *movelist) addBlackPawnCaptureMove(p *Position, from, to int, captured piece) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(p, newMovekey(from, to, captured, bQ, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, bR, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, bB, false, false))
		list.addCaptureMove(p, newMovekey(from, to, captured, bN, false, false))
	} else {
		list.addCaptureMove(p, newMovekey(from, to, captured, EMPTY, false, false))
	}
}

func (list *movelist) addBlackPawnMove(p *Position, from, to int) {
	if rankLookups[from] == RANK_2 {
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, bQ, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, bR, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, bB, false, false))
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, bN, false, false))
	} else {
		list.addCaptureMove(p, newMovekey(from, to, EMPTY, EMPTY, false, false))
	}
}

func generateAllMoves(p *Position, list *movelist) {
	// initialize
	list.count = 0

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
					list.addQuietMove(p, newMovekey(sq, sq+20, EMPTY, EMPTY, false, true))
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
				list.addCaptureMove(p, newMovekey(sq, lSq, EMPTY, EMPTY, true, false))
			}
			if rSq == p.enPas {
				list.addCaptureMove(p, newMovekey(sq, rSq, EMPTY, EMPTY, true, false))
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
					list.addQuietMove(p, newMovekey(sq, sq-20, EMPTY, EMPTY, false, true))
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
				list.addCaptureMove(p, newMovekey(sq, lSq, EMPTY, EMPTY, true, false))
			}
			if rSq == p.enPas {
				list.addCaptureMove(p, newMovekey(sq, rSq, EMPTY, EMPTY, true, false))
			}
		}
	}

	// set up all pieces per color
	var nonslidingPieces [2]piece
	//var slidePieces []piece
	if p.side == WHITE {
		nonslidingPieces = [2]piece{wN, wK}
		//slidePieces = []piece{wB, wR, wQ}
	} else {
		nonslidingPieces = [2]piece{bN, bK}
		//slidePieces = []piece{bB, bR, bQ}
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
					list.addQuietMove(p, newMovekey(fromSq, toSq, EMPTY, EMPTY, false, false))
				} else if pieceLookups[p.pieces[toSq]].color != p.side {
					list.addCaptureMove(p, newMovekey(fromSq, toSq, p.pieces[toSq], EMPTY, false, false))
				}
			}
		}
	}

	// sliding pieces
	//for _, pce := slidePieces {
	//	for pceNum := 0; pceNum < p.pieceCount[pce]; pceNum++ {
	//		fromSq := p.pieceList[pce][pceNum]
	//
	//		for _, dir := range dirKing {
	//			toSq := fromSq + dir
	//			if sqOffBoard(toSq) {
	//				continue
	//			}
	//			if p.pieces[toSq] == EMPTY {
	//				list.addQuietMove(p, newMovekey(fromSq, toSq, EMPTY, EMPTY, false, false))
	//			} else if pieceLookups[p.pieces[toSq]].color != p.side {
	//				list.addCaptureMove(p, newMovekey(fromSq, toSq, p.pieces[toSq], EMPTY, false, false))
	//			}
	//		}
	//	}
	//}
}
