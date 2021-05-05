package engine

import (
	"cacti-chess/engine/position"
	"fmt"
)

// Engine provides methods used to analyze chess games
type Engine struct {
	position *position.Position
}

// FromFen creates a new engine object initialized to a given fen
func FromFen(fen string) (*Engine, error) {
	engine := &Engine{}
	pos, err := position.FromFen(fen)
	if err != nil {
		return nil, fmt.Errorf("error parsing fen: %v", err)
	}
	engine.position = pos
	return engine, nil
}
