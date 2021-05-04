package engine

var perftCounts = map[string]int{}
var perftSection string
var perftSubmoves = map[string]string{}

func (p *Position) perftRecursive(totalDepth, depth int) int {
	err := p.assertCache()
	if err != nil {
		panic(err)
	}

	if depth == 0 {
		perftCounts[perftSection] += 1
		return 1
	}

	childNodes := 0
	mlist := p.generateAllMoves()

	// iterate all moves, depth first search
	for i := 0; i < mlist.count; i++ {
		if depth == totalDepth {
			perftSection = mlist.moves[i].key.ShortString()

		}

		// if the move leaves us in check, forget it
		if !p.MakeMove(mlist.moves[i].key) {
			continue
		}

		if depth != totalDepth {
			perftSubmoves[perftSection] += mlist.moves[i].key.ShortString() + ","
		}

		childNodes += p.perftRecursive(totalDepth, depth-1)
		p.UndoMove()
	}

	return childNodes
}

func (p *Position) Perft(depth int) int {
	return p.perftRecursive(depth, depth)
}
