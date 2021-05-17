package eval

import "cacti-chess/engine/position"

type PositionEvaluator struct{}

// Evaluate takes a position and returns a value based on the material
// and positional advantages. A positive number is an advantage for
// white, negative for black. Returned unit is in 100th of a pawn.
func (s PositionEvaluator) Evaluate(p *position.Position) float64 {
	// calculate initial material
	material := p.GetMaterial()
	pceCount := p.GetPieceCount()
	pceList := p.GetPieceList()

	score := material[position.WHITE] - material[position.BLACK]

	//calculate
	//piece
	//squares
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

	return float64(score)
}

// EvaluateAbsolute returns the same evaluation as Evaluate, but will
// always return a positive value even for black. This can be used
// for negamax implementations of minimax
func (s PositionEvaluator) EvaluateAbsolute(p *position.Position) float64 {
	if p.GetSide() == position.WHITE {
		return s.Evaluate(p)
	} else {
		return -s.Evaluate(p)
	}
}
