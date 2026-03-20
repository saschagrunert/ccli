package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/saschagrunert/ccli/v3"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := ccli.NewCommand()
	cmd.Name = "AppName"
	cmd.Usage = "App usage..."
	cmd.Version = "0.1.0"
	cmd.Description = "Application description"
	cmd.Copyright = fmt.Sprintf("(c) %d Some Company", time.Now().Year())
	cmd.Authors = []any{"Name <e@mail.com>"}
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	cmd.Action = func(_ context.Context, _ *cli.Command) error {
		fmt.Println("boom! I say!")

		return nil
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(1)
	}
}
