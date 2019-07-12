package main

import (
	"fmt"
	"os"
	"time"

	ccli "github.com/saschagrunert/ccli/pkg"
	"github.com/urfave/cli"
)

func main() {
	app := ccli.NewApp()
	app.Name = "AppName"
	app.Usage = "App usage..."
	app.Version = "0.1.0"
	app.Description = "Application description"
	app.Copyright = fmt.Sprintf("© %d Some Company", time.Now().Year())
	app.Authors = []cli.Author{{Name: "Name", Email: "e@mail.com"}}
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
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
