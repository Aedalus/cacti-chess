package search

import (
	"cacti-chess/engine/eval"
	"cacti-chess/engine/position"
	"fmt"
	"math"
	"time"
)

const maxDepth = 128
const mate = 29000 // score for mate

type SearchInfo struct {
	start time.Time
	stop  time.Time
	depth int
	nodes uint64 // the count of positions the engine visited

	scorer        Scorer
	searchPly     int
	searchHistory [13][120]int
	searchKillers [2][]int
	pvTable       *PrincipalVariationTable

	quit    bool // quit is set to true if forcefully exited
	stopped bool // stopped is more graceful

	fh  int // used to test alpha/beta efficiency
	fhf int // used to test alpha/beta efficiency

	// time controls
	depthset  int // max depth
	timeset   int // max time
	movestogo int
	infinite  int
}

func New() *SearchInfo {
	s := &SearchInfo{}
	s.scorer = &eval.PositionEvaluator{}
	s.Clear()
	return s
}

type Scorer interface {
	Evaluate(p *position.Position) float64
	EvaluateAbsolute(p *position.Position) float64
}

// Clear resets the search
func (s *SearchInfo) Clear() {
	// reset the search info
	s.start = time.Time{}
	s.stop = time.Time{}
	s.depth = 0
	s.nodes = 0

	s.searchPly = 0
	s.searchHistory = [13][120]int{}
	s.searchKillers = [2][]int{}
	s.pvTable = &PrincipalVariationTable{}

	s.quit = false
	s.stopped = false

	// time controls
	s.depthset = 0
	s.timeset = 0
	s.movestogo = 0
	s.infinite = 0

	s.fh = 0
	s.fhf = 0
}

func (s *SearchInfo) GetPrincipalVariationLine(p *position.Position) []position.Movekey {
	return s.pvTable.GetBestLine(p)
}

func (s *SearchInfo) GetPrincipalVariationTable() *PrincipalVariationTable {
	return s.pvTable
}

/*
AlphaBeta is a combination of Negamax search + AlphaBeta pruning.

When analyzing to a single depth, we're looking at the max(eval(depth)) where depth = 1.
This becomes more complicated as multiple depths are introduced. We're no longer looking for the best move we
can do, but rather the move that leaves the opponent with the worst options. i.e. it doesn't help to capture
a piece if you then lose your queen.

Assuming positive scores are good for white, white always wants the MAXIMUM value possible.
Because negative scores are good for black, black always wants the MINIMUM value possible.
We then have to alternate between two functions, min/max at each depth.

For example, assuming we're looking at a position where it's white's turn. If depth 1 is making a single move,
depth 0 is just the current state of the board.

- max depth 1
	- the depth 0 node is equal to the MAX of all depth 1 nodes, since white gets to choose
	- each depth 1 node is equal to the board evaluation.
- max depth 2
	- the depth 0 node is equal to the MAX of all depth 1 nodes, since white gets to choose
	- each depth 1 node is equal to the MIN of all depth 2 nodes UNDER IT, since black gets to choose
	- each depth 2 node is equal to the board evaluation.
- max depth 3
	- the depth 0 node is equal to the MAX of all depth 1 nodes, since white gets to choose
	- each depth 1 node is equal to the MIN of all depth 2 nodes UNDER IT, since black gets to choose
	- each depth 2 node is equal to the MAX of all depth 3 nodes UNDER IT, since white gets to choose
	- each depth 3 node is equal to the board evaluation.

We can see that there is a base case where if depth == max depth, we evaluate the board position. However
for each level above it we alternate between min/max.

We could write two separate min/max functions to achieve this. However there's a mathmatical relation
that gives us a trick. First, imagine scoring is done relative to the player. The value of a position
for one player is equal to the negation of the value for the other player.

value for white = -1 * value for black
value for black = -1 * value for white

i.e
	if it's +3 for me, it's -3 for you
	if it's -5 for me, it's +5 for you

This allows us to put it in more relative terms. On a players move, they're looking for a move that
MAXIMIZES the NEGATION of all subsequent moves. Let's say we have the choice of two moves, A + B.
A let's our opponent score (5,10,15) relative to THEM. B let's them score (1,2,3) relative to THEM.
This becomes (-5,-10,-15), (-1,-2,-3) relative to US. We would want to choose the move that MAXIMIZES
the score of negated nodes underneath.

Notice that this means the scoring function must always return the value relative for the current player.


// ALPHA BETA PRUNING
Alpha and Beta also take on relative values. So alpha is always equal to the best known score for the
current player. Beta is the best known for the opponent. If we are currently looking at

alpha = 20, beta = -10

and do a hypothetically neutral move, when we call the recursive function again we need to flip the values
and negate them, so the next iteration shows

alpha = 10, beta = -20
*/

func (s *SearchInfo) AlphaBeta(p *position.Position, alpha, beta float64, depth int, doNull bool) float64 {

	// base case it's a leaf node, we return the evaluation relative to the current player.
	if depth == 0 {
		s.nodes++
		return s.scorer.EvaluateAbsolute(p)
	}

	// edge cases for repetition, or if we are too far down return 0 for a draw
	if p.IsRepetition() || p.GetFiftyMove() >= 100 {
		return 0
	}

	// endgame scenario, few pieces so we're searching a lot of depth
	if s.searchPly > maxDepth {
		return s.scorer.EvaluateAbsolute(p)
	}

	movelist := p.GenerateAllMoves()

	legal := 0
	oldAlpha := alpha
	bestMove := position.Movekey(0)

	for _, mv := range *movelist {
		// if it's not legal, auto undo
		if !p.MakeMove(mv.Key) {
			continue
		}

		legal++

		// We call alphaBeta again to find the best response from our opponent
		// we negate the return value so it's relative to US.
		// We flip the alpha/beta order and sign as well for the same reason
		score := -s.AlphaBeta(p, -beta, -alpha, depth-1, true)
		p.UndoMove()

		// evaluate if this is better than what we've seen
		if score > alpha {
			if score >= beta {
				if legal == 1 {
					s.fhf++
				}
				s.fh++
				return beta
			}
			alpha = score
			bestMove = mv.Key
		}
	}

	// check for stalemate/checkmate
	if legal == 0 {
		// if we're mated, return the low mate score with the depth to mate
		// added. i.e. mate in 2 is -28998, 3 -28997
		if p.IsKingAttacked() {
			return float64(-mate + s.searchPly)
		} else {
			// stalemate
			return 0
		}
	}

	if alpha != oldAlpha {
		s.pvTable.Set(p, bestMove)
	}

	return alpha
}

func (s *SearchInfo) SearchPosition(p *position.Position) {

	bestMove := position.Movekey(0)
	bestScore := math.Inf(-1)
	currentDepth := 0
	var pvMoves []position.Movekey

	// todo - cleanup
	s.pvTable = &PrincipalVariationTable{}

	// iterative deepening
	for i := 1; i <= s.depth; i++ {
		bestScore = s.AlphaBeta(p, math.Inf(-1), math.Inf(1), currentDepth, true)
		fmt.Printf("depth: %v, side: %v, score: %v, move: %v, nodes: %v", currentDepth, p.GetSide(), bestScore, bestMove.ShortString(), s.nodes)

		for _, pvm := range pvMoves {
			fmt.Print(pvm)
		}
		fmt.Print("\n")
	}
}
