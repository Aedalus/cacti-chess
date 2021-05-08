package search

import "cacti-chess/engine/position"

type PrincipalVariationEntry struct {
	PosKey uint64
	move   int
}

type PrincipalVariationTable map[uint64]position.Movekey

func (tb *PrincipalVariationTable) Set(p *position.Position, mv position.Movekey) {
	(*tb)[p.GetPosKey()] = mv
}

func (tb *PrincipalVariationTable) Probe(p *position.Position) position.Movekey {
	return (*tb)[p.GetPosKey()]
}

func (tb PrincipalVariationTable) GetPVLine(p *position.Position) []position.Movekey {
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
