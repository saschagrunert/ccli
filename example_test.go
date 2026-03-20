package ccli_test

import (
	"github.com/fatih/color"
	"github.com/saschagrunert/ccli/v2"
)

const (
	exampleAppName = "myapp"
	exampleUsage   = "does something useful"
	exampleVersion = "1.0.0"
)

func Example() {
	app := ccli.NewApp()
	app.Name = exampleAppName
	app.Usage = exampleUsage
	app.Version = exampleVersion

	// app.Run(os.Args)
}

func ExampleNewAppWithOptions() {
	app := ccli.NewAppWithOptions(ccli.Options{
		Green:   color.New(color.FgHiGreen).SprintFunc(),
		Yellow:  color.New(color.FgHiYellow).SprintFunc(),
		Blue:    nil,
		Cyan:    nil,
		Red:     nil,
		Disable: false,
	})
	app.Name = exampleAppName
	app.Usage = exampleUsage

	// app.Run(os.Args)
}

func ExampleNewAppWith() {
	app := ccli.NewAppWith(
		ccli.WithGreen(color.New(color.FgHiGreen).SprintFunc()),
		ccli.WithYellow(color.New(color.FgHiYellow).SprintFunc()),
	)
	app.Name = exampleAppName
	app.Usage = exampleUsage

	// app.Run(os.Args)
}

func ExampleNewAppWith_disable() {
	app := ccli.NewAppWith(ccli.WithDisable())
	app.Name = exampleAppName
	app.Usage = exampleUsage

	// app.Run(os.Args)
}
