package eval

import "cacti-chess/engine/position"

// EvalPosition returns the evaluation of a given
// position in hundredths of a pawn. Always returns
// the score as positive from the perspective of
// the current side, even if black.

type PositionScorer struct{}

func (s PositionScorer) Score(p *position.Position) int {
	return evalPosition(p)
}

func evalPosition(p *position.Position) int {
	// calculate initial material
	material := p.GetMaterial()
	pceCount := p.GetPieceCount()
	pceList := p.GetPieceList()

	score := material[position.WHITE] - material[position.BLACK]

	// calculate piece squares
	pce := position.PwP
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score += pawnTable[position.SQ64(sq120)]
	}

	pce = position.PbP
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score -= pawnTable[mirror64[position.SQ64(sq120)]]
	}

	pce = position.PwN
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score += knightTable[position.SQ64(sq120)]
	}

	pce = position.PbN
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score -= knightTable[mirror64[position.SQ64(sq120)]]
	}

	pce = position.PwB
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score += bishopTable[position.SQ64(sq120)]
	}

	pce = position.PbB
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score -= bishopTable[mirror64[position.SQ64(sq120)]]
	}

	pce = position.PwR
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score += rookTable[position.SQ64(sq120)]
	}

	pce = position.PbR
	for i := 0; i < pceCount[pce]; i++ {
		sq120 := pceList[pce][i]
		score -= rookTable[mirror64[position.SQ64(sq120)]]
	}

	return score
}
