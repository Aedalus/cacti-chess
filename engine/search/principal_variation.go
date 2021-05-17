package search

import "cacti-chess/engine/position"

// PrincipalVariationTable can be used to reconstruct the best
// line found from a position. As the AlphaBeta function evaluates,
// it stores key/values for the poskey of a position with the best known
// move. The "principal variation" is the best known line for a given position
type PrincipalVariationTable map[uint64]position.Movekey

// Set adds a new movekey, retrieving the poskey from the position
func (tb *PrincipalVariationTable) Set(p *position.Position, mv position.Movekey) {
	(*tb)[p.GetPosKey()] = mv
}

// Probe gets the movekey from a given position
func (tb *PrincipalVariationTable) Probe(p *position.Position) position.Movekey {
	return (*tb)[p.GetPosKey()]
}

// GetBestLine returns the best known line as a slice of movekeys. It will also
// undo all the moves back to before the alpha/beta started.
func (tb PrincipalVariationTable) GetBestLine(p *position.Position) []position.Movekey {
	var moves []position.Movekey

	move := tb[p.GetPosKey()]

	for move != position.Movekey(0) {
		if p.MoveExists(move) {
			moves = append(moves, move)
			p.MakeMove(move)
			move = tb[p.GetPosKey()]
		} else {
			panic("move does not exist!")
		}
	}

	for i := 0; i < len(moves); i++ {
		p.UndoMove()
	}

	return moves
}
