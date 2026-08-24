package ccli_test

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/saschagrunert/ccli/v3"
	"github.com/urfave/cli/v3"
)

const (
	exampleAppName = "myapp"
	exampleUsage   = "does something useful"
	exampleVersion = "1.0.0"
)

func Example() {
	cmd := ccli.NewCommand()
	cmd.Name = exampleAppName
	cmd.Usage = exampleUsage
	cmd.Version = exampleVersion

	// cmd.Run(context.Background(), os.Args)
}

func ExampleNewCommandWithOptions() {
	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Green:   color.New(color.FgHiGreen).SprintFunc(),
		Yellow:  color.New(color.FgHiYellow).SprintFunc(),
		Blue:    nil,
		Cyan:    nil,
		Red:     nil,
		Disable: false,
	})
	cmd.Name = exampleAppName
	cmd.Usage = exampleUsage

	// cmd.Run(context.Background(), os.Args)
}

func ExampleNewCommandWith() {
	cmd := ccli.NewCommandWith(
		ccli.WithGreen(color.New(color.FgHiGreen).SprintFunc()),
		ccli.WithYellow(color.New(color.FgHiYellow).SprintFunc()),
	)
	cmd.Name = exampleAppName
	cmd.Usage = exampleUsage

	// cmd.Run(context.Background(), os.Args)
}

func ExampleApply() {
	cmd := ccli.NewCommand()
	cmd.Name = exampleAppName
	cmd.Usage = exampleUsage
	cmd.Commands = []*cli.Command{
		{Name: "serve", Usage: "start the server"},
	}

	ccli.Apply(cmd)

	// cmd.Run(context.Background(), os.Args)
}

func ExampleNewCommandWith_disable() {
	cmd := ccli.NewCommandWith(ccli.WithDisable())
	cmd.Name = exampleAppName
	cmd.Usage = exampleUsage

	// cmd.Run(context.Background(), os.Args)
}

func ExampleNewCommandWith_run() {
	cmd := ccli.NewCommandWith(ccli.WithDisable())
	cmd.Name = "demo"
	cmd.Action = func(_ context.Context, _ *cli.Command) error {
		fmt.Println("hello from demo")

		return nil
	}

	err := cmd.Run(context.Background(), []string{"demo"})
	if err != nil {
		return
	}
	// Output: hello from demo
}
