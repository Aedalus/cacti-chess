package main

import (
	"bufio"
	"cacti-chess/engine/position"
	"cacti-chess/engine/search"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
	"strings"
)

var playCmd = &cli.Command{
	Name:  "play",
	Usage: "plays a game of chess",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "fen",
			Aliases: []string{"f"},
			Usage:   "an optional fen to start from (default starting position)",
		},
		&cli.BoolFlag{
			Name:    "black",
			Aliases: []string{"b"},
			Usage:   "play black",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		// read in flags
		fen := c.String("fen")
		if fen == "" {
			fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		}
		playBlack := c.Bool("black")

		// initialize position
		fmt.Printf("loading game from: %v\n", fen)
		isPlayerTurn := !playBlack
		p, err := position.FromFen(fen)
		if err != nil {
			log.Fatalf("could not parse fen: %v", err)
		}

		for true {
			// print the board
			fmt.Println(p)

			// check for win conditions
			if p.IsLegalMove() == false {
				if p.IsStalemate() {
					fmt.Println("Stalemate!")
				} else if p.IsCheckmate() {
					fmt.Println("Checkmate!")
				}
			}

			// player/engine turn
			if isPlayerTurn {
				doPlayerTurn(p)
			} else {
				doEngineTurn(p)
			}

			// switch sides
			isPlayerTurn = !isPlayerTurn

		}
		return nil
	},
}

// doPlayerTurn reads from stdin and makes the given move
func doPlayerTurn(p *position.Position) {
	reader := bufio.NewReader(os.Stdin)
	// prompt until we get good input
	for true {
		// read in the next move
		fmt.Print("Enter move: ")
		mvStr, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading from stdin")
		}

		mvStr = strings.TrimSpace(mvStr)
		if mvStr == "quit" || mvStr == "exit" {
			fmt.Println("exiting...")
			os.Exit(0)
		}

		if mvStr == "undo" {
			p.UndoMove()
			p.UndoMove()
			break
		}

		// parse the given move
		movekey, err := p.ParseMove(mvStr)
		if err != nil || movekey.IsNoMove() {
			fmt.Println("move not recognized, try again")
		} else {
			p.MakeMove(movekey)
			break
		}
	}
}

// doEngineTurn searches for the best move and then performs it
func doEngineTurn(p *position.Position) {
	search := search.New()
	search.AlphaBeta(p, math.Inf(-1), math.Inf(1), 5, false)
	line := search.GetPrincipalVariationLine(p)
	p.MakeMove(line[0])
}
