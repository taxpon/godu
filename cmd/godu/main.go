package main

import (
	"github.com/taxpon/godu"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "godu"
	app.Usage = "Get disk usage info with useful options"
	app.Author = "Takuro Wada"
	app.Email = "taxpon@gmail.com"
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "recursive, r", Usage: "calculate disk usage recursively"},
		cli.BoolFlag{Name: "absolute, a", Usage: "show result using absolute path"},
		cli.BoolFlag{Name: "dump, d", Usage: "dump resutl as message pack binary"},
		cli.BoolFlag{Name: "load, l", Usage: "load resutl stored as message pack binary"},
	}

	app.Action = func(c *cli.Context) error {
		targetPath := ""
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			os.Exit(2)
		}

		targetPath = c.Args().Get(0)
		err := godu.Run(targetPath,
			c.Bool("recursive"),
			c.Bool("absolute"),
			c.Bool("dump"),
			c.Bool("load"),
		)

		return err
	}

	app.Run(os.Args)
}
