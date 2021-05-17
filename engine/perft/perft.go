package perft

import "cacti-chess/engine/position"

var perftCounts = map[string]int{}
var perftSection string

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
	movelist := p.GenerateAllMoves()

	// iterate all moves, depth first search
	for _, mv := range *movelist {
		if depth == totalDepth {
			perftSection = mv.Key.ShortString()
		}

		// if the move leaves us in check, forget it
		if !p.MakeMove(mv.Key) {
			continue
		}

		childNodes += perftRecursive(p, totalDepth, depth-1)
		p.UndoMove()
	}

	return childNodes
}

// Perft returns the total number of possible moves at a given depth.
// The standard starting position has 20 possible moves for the first moves,
// 400 possible moves for the response, etc. This is used to ensure the engine's
// move generation matches established results, eliminating edge cases around
// castling, en passant, etc.
func Perft(p *position.Position, depth int) int {
	return perftRecursive(p, depth, depth)
}
