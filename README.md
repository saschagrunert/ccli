# ccli
## Command line parsing in go, with coloring support

This project uses the already existing go cli parsing solution and adds additional coloring support.
Some strong defaults are provided as well. An integration example can be found at `cmd/main.go`.

## Usage
Install the package with:

```
go get github.com/saschagrunert/ccli
```

Afterwards it can be used like the `cli` package:

```go
package main

import (
	"fmt"
	"github.com/saschagrunert/ccli"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	app := ccli.NewApp()
	app.Name = "AppName"
	app.Usage = "App usage..."
	app.Version = "0.1.0"
	app.Description = "Application description"
	app.Copyright = fmt.Sprintf("© %d Some Company", time.Now().Year())
	app.Authors = []cli.Author{cli.Author{Name: "Name", Email: "e@mail.com"}}
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
```

If you run this application, the help output will look like this:

![screenshot](.github/screenshot.png)
