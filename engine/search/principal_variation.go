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

func (tb PrincipalVariationTable) GetPVLine(depth int, p *position.Position) {
	move := tb[p.GetPosKey()]

	for move != position.Movekey(0) {
		if p.MoveExists(move) {
			p.MakeMove(move)
		} else {
			break
		}

		move = tb[p.GetPosKey()]
	}

	for p.GetSearchPly() > 0 {
		p.UndoMove()
	}
}
