package main

import (
	"bufio"
	"cacti-chess/engine/position"
	"cacti-chess/engine/search"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// https://www.shredderchess.com/chess-features/uci-universal-chess-interface.html

var logFile *os.File

func init() {
	file, err := os.Create("uci.log")
	if err != nil {
		panic(err)
	}
	logFile = file
}
func main() {
	fmt.Fprintln(logFile, "starting")
	reader := bufio.NewReader(os.Stdin)
	client := &UCIClient{
		search: search.New(),
	}

	for {
		text, err := reader.ReadString('\n')
		//fmt.Printf("info input %q\n", text)
		if err != nil {
			log.Fatalf("error reading input: %v", err)
		}

		text = strings.TrimSpace(text)
		client.parseLine(text)
	}
}

type UCIClient struct {
	position *position.Position
	search   *search.SearchInfo
}

func (c *UCIClient) parseLine(line string) {
	segments := strings.Split(line, " ")

	switch segments[0] {
	case "isready":
		fmt.Println("readyok")
	case "position":
		fmt.Fprintln(logFile, strings.Join(segments, " "))
		c.parsePosition(segments)
	case "ucinewgame":
		c.parsePosition([]string{"position", "startpos"})
	case "go":
		c.parseGo(segments)
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

type GoCmdArgs struct {
	SearchMoves []string      // Only search specific moves
	Ponder      bool          // Keep thinking, even if checkmate
	Wtime       time.Duration // Time remaining for White
	Btime       time.Duration // Time remaining for Black
	Winc        time.Duration // Time increment for White
	Binc        time.Duration // Time increment for Black
	MovesToGo   int           // n moves til time control
	Depth       int           // search n plys
	Nodes       int           // search n nodes
	Mate        int           // search for a mate in x moves
	MoveTime    time.Duration // search an exact duration
	Infinite    bool          // search until 'stop' command
}

func parseGoCmdArgs(segments []string) GoCmdArgs {
	goCmdArgs := GoCmdArgs{
		SearchMoves: []string{},
		Ponder:      false,
		Wtime:       0,
		Btime:       0,
		Winc:        0,
		Binc:        0,
		MovesToGo:   0,
		Depth:       5,
		Nodes:       0,
		Mate:        0,
		MoveTime:    0,
		Infinite:    false,
	}

	for i, arg := range segments {
		switch arg {
		case "searchmoves":
			for j := i + 1; j < len(segments); j++ {
				goCmdArgs.SearchMoves = append(goCmdArgs.SearchMoves, segments[j])
			}
		case "ponder":
			goCmdArgs.Ponder = true
		case "wtime":
			t, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Wtime = time.Millisecond * time.Duration(t)
		case "btime":
			t, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Btime = time.Millisecond * time.Duration(t)
		case "winc":
			t, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Winc = time.Millisecond * time.Duration(t)
		case "binc":
			t, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Binc = time.Millisecond * time.Duration(t)
		case "movestogo":
			mvc, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.MovesToGo = mvc
		case "depth":
			n, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Depth = n
		case "nodes":
			n, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Nodes = n
		case "mate":
			n, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.Mate = n
		case "movetime":
			t, _ := strconv.Atoi(segments[i+1])
			goCmdArgs.MoveTime = time.Millisecond * time.Duration(t)
		case "infinite":
			goCmdArgs.Infinite = true
		}
	}
	return goCmdArgs
}

func (c *UCIClient) parseGo(segments []string) {
	goCmdArgs := parseGoCmdArgs(segments)

	c.search = search.New()

	_, line := c.search.SearchPosition(c.position, search.Options{
		Depth: goCmdArgs.Depth,
	})

	fmt.Fprintf(logFile, "found bestmove %s\n", line[0].ShortString())
	fmt.Printf("bestmove %v\n", line[0].ShortString())
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
		if len(segments) > 2 {
			moves = segments[3:]
		}
	} else if segments[1] == "fen" {
		// fen has 6 parts
		fen = strings.Join(segments[2:8], " ")
		if len(segments) >= 10 {
			moves = segments[9:]
		}
	}

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
}
