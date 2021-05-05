package perft

import "cacti-chess/engine/position"

var perftCounts = map[string]int{}
var perftSection string
var perftSubmoves = map[string]string{}

func perftRecursive(p *position.Position, totalDepth, depth int) int {

	err := p.AssertCache()
	if err != nil {
		panic(err)
	}

	if depth == 0 {
		perftCounts[perftSection] += 1
		return 1
	}

	childNodes := 0
	mlist := p.GenerateAllMoves()

	// iterate all moves, depth first search
	for i := 0; i < mlist.Count; i++ {
		if depth == totalDepth {
			perftSection = mlist.Moves[i].Key.ShortString()

		}

		// if the move leaves us in check, forget it
		if !p.MakeMove(mlist.Moves[i].Key) {
			continue
		}

		//if depth != totalDepth {
		//	perftSubmoves[perftSection] += mlist.moves[i].key.ShortString() + ","
		//}

		childNodes += perftRecursive(p, totalDepth, depth-1)
		p.UndoMove()
	}

	return childNodes
}

func Perft(p *position.Position, depth int) int {
	return perftRecursive(p, depth, depth)
}
