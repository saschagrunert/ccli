package ccli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"time"
)

// AppSettings specifies the application options
type AppSettings struct {
	Name        string
	Usage       string
	Version     string
	Description string
	Company     string
	Authors     []cli.Author
}

// NewApp creates a new applications with the given settings
func NewApp(settings AppSettings) *cli.App {
	cli.AppHelpTemplate = getHelpTemplate()

	app := cli.NewApp()
	app.Authors = settings.Authors
	app.Compiled = time.Now()
	app.Copyright = fmt.Sprintf("© %d %s", time.Now().Year(), settings.Company)
	app.Description = settings.Description
	app.Name = settings.Name
	app.Usage = settings.Usage
	app.Version = settings.Version

	return app
}

func getHelpTemplate() string {
	// Set the colors
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// Generate the template string
	a := green("{{.Name}}") + " {{.Version}}\n"
	a += "{{if .Usage}}{{.Usage}}{{end}}\n\n"
	a += yellow("USAGE:\n")
	a += "    {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} "
	a += "{{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}\n\n"
	a += "{{end}}{{end}}{{if .Description}}"
	a += yellow("DESCRIPTION:\n")
	a += "    {{.Description}}\n\n{{end}}{{if len .Authors}}"
	a += yellow("AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:\n")
	a += "    {{range $index, $author := .Authors}}{{if $index}}\n"
	a += "    {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}\n\n"
	a += yellow("COMMANDS:") + "{{range .VisibleCategories}}{{if .Name}}\n"
	a += green("    {{.Name}}")
	a += ":{{end}}{{range .VisibleCommands}}\n"
	a += green(`    {{join .Names ", "}}{{"\t"}}`)
	a += "{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}\n\n"
	a += yellow("GLOBAL OPTIONS:")
	a += `
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

`
	a += blue("{{.Copyright}}{{end}}\n")
	return a
}
