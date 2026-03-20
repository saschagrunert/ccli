# ccli

[![ci](https://github.com/saschagrunert/ccli/actions/workflows/test.yml/badge.svg)](https://github.com/saschagrunert/ccli/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/saschagrunert/ccli/v2.svg)](https://pkg.go.dev/github.com/saschagrunert/ccli/v2)

## Command line parsing in Go, with coloring support

This package wraps [urfave/cli](https://github.com/urfave/cli) and adds
coloring to help output. Section headers, command names, author info, and
copyright each get their own color.

![screenshot](.github/screenshot.png)

## Usage

Install the package with:

```shell
go get github.com/saschagrunert/ccli/v2
```

Then use it like the `cli` package:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/saschagrunert/ccli/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	app := ccli.NewApp()
	app.Name = "AppName"
	app.Usage = "App usage..."
	app.Version = "0.1.0"
	app.Description = "Application description"
	app.Copyright = fmt.Sprintf("(c) %d Some Company", time.Now().Year())
	app.Authors = []*cli.Author{{Name: "Name", Email: "e@mail.com"}}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}
	app.Action = func(_ *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
```

## Custom colors

Use functional options with `NewAppWith`:

```go
import (
	"github.com/fatih/color"
	"github.com/saschagrunert/ccli/v2"
)

app := ccli.NewAppWith(
	ccli.WithGreen(color.New(color.FgHiGreen).SprintFunc()),
	ccli.WithYellow(color.New(color.FgHiYellow).SprintFunc()),
)
```

Or pass an `Options` struct to `NewAppWithOptions`:

```go
app := ccli.NewAppWithOptions(ccli.Options{
	Green:  color.New(color.FgHiGreen).SprintFunc(),
	Yellow: color.New(color.FgHiYellow).SprintFunc(),
})
```

Any color left unset falls back to its default. To turn off all coloring:

```go
app := ccli.NewAppWith(ccli.WithDisable())
```

## Notes

Calling `NewApp` or `NewAppWithOptions` sets `cli.CommandHelpTemplate` and
`cli.SubcommandHelpTemplate` as package-level globals. The app help template
is set per-app via `CustomAppHelpTemplate` and does not affect other apps.

## License

[MIT](LICENSE)
