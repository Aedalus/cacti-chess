package engine

type movescore struct {
	key   *movekey
	score int
}

type movelist struct {
	moves [256]movescore
	count int
}

func (list *movelist) addQuietMove(s *State, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func (list *movelist) addCaptureMove(s *State, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func (list *movelist) addEnPasMove(s *State, move *movekey) {
	list.moves[list.count].key = move
	list.moves[list.count].score = 0
	list.count++
}

func generateAllMoves(s *State, list *movelist) {

	list.count = 0
}
