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
		cli.BoolFlag{
			Name: "recursive, r",
			Usage: "calculate disk usage recursively"},
		cli.BoolFlag{
			Name: "absolute, a",
			Usage: "show result using absolute path"},
		cli.BoolFlag{
			Name: "dump, d",
			Usage: "dump resutl as message pack binary"},
	}

	app.Commands = []cli.Command{
		{
			Name: "compare",
			Aliases: []string{"c"},
			Usage: "compare 2 archived data",
			Action: func(c *cli.Context) error {
				err := godu.Compare(c.Args().Get(0), c.Args().Get(1))
				return err
			},
		},
		{
			Name:    "load",
			Aliases: []string{"l"},
			Usage:   "load saved usage information",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list, l",
					Usage: "list all saved file",
				},
				cli.StringFlag{
					Name:  "input, i",
					Value: "dump.bin",
					Usage: "input file to load",
				},
			},
			Action: func(c *cli.Context) error {
				err := godu.Load(
					c.Bool("list"),
					c.String("input"),
				)
				if err != nil {
					return err
				}
				return nil
			},
		},
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
		)

		return err
	}

	app.Run(os.Args)
}
