package main

import (
	"fmt"
	"github.com/saschagrunert/ccli"
	"github.com/urfave/cli"
	"os"
)

func main() {
	settings := ccli.AppSettings{
		Name:        "AppName",
		Usage:       "App usage...",
		Version:     "0.1.0",
		Description: "Application description",
		Company:     "Some company",
		Authors:     []cli.Author{cli.Author{Name: "Name", Email: "e@mail.com"}},
	}
	app := ccli.NewApp(settings)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}
	app.Run(os.Args)
}
