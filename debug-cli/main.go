package main

import (
	"cacti-chess/engine/position"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Name: "ccd",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "playground",
				Usage: "allows a game to be played adding moves for both sides",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "fen",
						Aliases:     []string{"f"},
						Usage:       "an optional fen to start from",
						DefaultText: "A",
					},
				},
				Action: func(c *cli.Context) error {
					fen := c.String("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
					p, err := position.FromFen(fen)
					if err != nil {
						log.Fatalf("could not parse fen: %v", err)
					}

					p.PrintBoard()
					return nil
				},
			},
			//{
			//	Name:    "draw-fen",
			//	Aliases: []string{"d"},
			//	Usage:   "draws a board to the screen from a given fen",
			//	Flags: []cli.Flag{
			//		&cli.StringFlag{
			//			Name:     "fen",
			//			Aliases:  []string{"f"},
			//			Usage:    "the fen string",
			//			Required: true,
			//		},
			//	},
			//	Action: func(c *cli.Context) error {
			//		fen := c.String("fen")
			//		state, err := position.ParseFen(fen)
			//		if err != nil {
			//			log.Fatalf("error parsing fen: %v", err)
			//		}
			//		fmt.Printf("%v", state.PrintBoard())
			//		return nil
			//	},
			//},

			{
				Name:    "draw-attacks",
				Aliases: []string{"d"},
				Usage:   "draws a board to the screen from a given fen",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "fen",
						Aliases:  []string{"f"},
						Usage:    "the fen string",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "color",
						Aliases:  []string{"c"},
						Usage:    "the attacking color to print",
						Value:    "white",
						Required: true,
					},
				},
				//Action: func(c *cli.Context) error {
				//	fen := c.String("fen")
				//	engine, err := engine.FromFen(fen)
				//	color := 0
				//	if c.String("color") == "black" || c.String("color") == "b" {
				//		color = 1
				//	}
				//	if err != nil {
				//		log.Fatalf("error parsing fen: %v", err)
				//	}
				//	fmt.Printf("%v", engine.PrintAttackBoard(color))
				//	return nil
				//},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
