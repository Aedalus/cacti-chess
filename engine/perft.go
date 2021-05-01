package engine

func (p *Position) Perft(depth int) int {
	err := p.assertCache()
	if err != nil {
		panic(err)
	}

	if depth == 0 {
		return 1
	}

	childNodes := 0
	mlist := p.generateAllMoves()

	// iterate all moves, depth first search
	for i := 0; i < mlist.count; i++ {

		//fmt.Printf("move: %v\n", mlist.moves[i].key.ShortString())
		// if the move leaves us in check, forget it
		if !p.MakeMove(mlist.moves[i].key) {
			continue
		}

		childNodes += p.Perft(depth - 1)
		p.UndoMove()
	}

	return childNodes
}
