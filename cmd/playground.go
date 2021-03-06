package main

import (
	"bufio"
	"cacti-chess/engine/position"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

var playgroundCmd = &cli.Command{
	Name:  "playground",
	Usage: "allows a game to be played adding moves for both sides",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "fen",
			Aliases: []string{"f"},
			Usage:   "an optional fen to start from (default starting position)",
		},
	},
	Action: func(c *cli.Context) error {
		fen := c.String("fen")
		if fen == "" {
			fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		}
		fmt.Printf("loading game from: %v\n", fen)
		p, err := position.FromFen(fen)
		if err != nil {
			log.Fatalf("could not parse fen: %v", err)
		}

		reader := bufio.NewReader(os.Stdin)
		for true {
			// print the board
			fmt.Println(p)

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
		return nil
	},
}
