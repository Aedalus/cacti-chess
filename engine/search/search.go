package search

import (
	"cacti-chess/engine/position"
	"time"
)

type SearchInfo struct {
	start time.Time
	stop  time.Time
	depth int
	nodes uint64 // the count of positions the engine visited

	searchPly     int
	searchHistory [13][120]int
	searchKillers [2][]int
	pvTable       *PrincipalVariationTable

	quit    bool // quit is set to true if forcefully exited
	stopped bool // stopped is more graceful

	// time controls
	depthset  int // max depth
	timeset   int // max time
	movestogo int
	infinite  int
}

func New() *SearchInfo {
	s := &SearchInfo{}
	s.Clear()
	return s
}

type Scorer interface {
	Score(p *position.Position) int
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
}

func (s *SearchInfo) SearchPosition(p *position.Position) {

}

func (s *SearchInfo) CheckUp() {
	// checks if it needs to report back
}

func (s *SearchInfo) Quiescence(p *position.Position, alpha, beta int) int {
	return 0
}

func (s *SearchInfo) AlphaBeta(p *position.Position, alpha, beta, depth, doNull int) int {
	return 0
}
