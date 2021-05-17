package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Name: "cacti-chess",
		Commands: []*cli.Command{
			playgroundCmd,
			playCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
