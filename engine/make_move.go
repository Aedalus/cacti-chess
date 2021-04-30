package engine

import "fmt"

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

// clear a single sq from the position
func (p *Position) clearPiece(sq int) {
	if sqOffBoard(sq) {
		panic(fmt.Sprintf("sq %v is off board", sq))
	}

	pce := p.pieces[sq]
	pceMeta := pieceLookups[pce]

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

	if foundPieceIndex == -1 {
		panic("piece should have been found in pieceList")
	}

	// copy the last item in the list to the found index
	p.pieceList[pce][foundPieceIndex] = p.pieceList[pce][p.pieceCount[pce]]
	// delete the last index, since we copied it forward
	p.pieceList[pce][p.pieceCount[pce]] = 0
	// decrement the total piece count to match
	p.pieceCount[pce]--
}
