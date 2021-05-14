package main

import (
	"bufio"
	"cacti-chess/engine/position"
	"fmt"
	"log"
	"os"
	"strings"
)

// https://www.shredderchess.com/chess-features/uci-universal-chess-interface.html

func main() {
	reader := bufio.NewReader(os.Stdin)
	client := &UCIClient{}

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %v", err)
		}

		text = strings.TrimSpace(text)
		client.parseLine(text)
	}
}

type UCIClient struct {
	position *position.Position
}

func (c *UCIClient) parseLine(line string) {
	segments := strings.Split(line, " ")

	switch segments[0] {
	case "isready":
		fmt.Println("readyok")
	case "position":
		c.parsePosition(segments)
	case "ucinewgame":
		c.parsePosition([]string{"position", "startpos"})
	case "go":
		//todo - parse go
	case "uci":
		fmt.Println("id name cacti-chess")
		fmt.Println("id author aedalus")
		fmt.Println("uciok")
	case "quit":
		os.Exit(0)
	case "":
	default:
		log.Fatalf("cmd not recognized: %q", segments[0])
	}
}

func (c *UCIClient) parsePosition(segments []string) {
	if len(segments) < 2 {
		log.Fatalf("error parsing position line %q. Expected > 2 segments", segments)
	}

	var fen string
	var moves []string

	// parse fen/startingpos
	if segments[1] == "startpos" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		moves = segments[2:]
	} else if segments[1] == "fen" {
		// fen has 6 parts
		fen = strings.Join(segments[2:8], " ")
		if len(segments) >= 10 {
			moves = segments[9:]
		}
	}

	fmt.Println(fen)
	fmt.Println(moves)

	// create the position
	pos, err := position.FromFen(fen)
	if err != nil {
		log.Fatalf("error parsing fen: %v", err)
	}

	// apply any existing moves
	for _, mv := range moves {
		movekey, err := pos.ParseMove(mv)
		if err != nil {
			log.Fatalf("error parsing move %q: %v", mv, err)
		}
		pos.MakeMove(movekey)
	}

	// set as the active position
	c.position = pos

	fmt.Println(c.position)
}
