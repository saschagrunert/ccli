package ccli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"time"
)

// NewApp creates a new applications with the given settings
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Compiled = time.Now()
	setAppTemplates()
	return app
}

func setAppTemplates() {
	// Set the colors
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// Set the application help template
	cli.AppHelpTemplate = fmt.Sprintf(`%s {{if .Version}}{{if not .HideVersion}}{{.Version}}{{end}}{{end}}
{{if .Usage}}{{.Usage}}{{end}}

%s
    {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Description}}

%s
    {{.Description}}{{end}}{{if len .Authors}}

%s{{with $length := len .Authors}}{{if ne 1 $length}}%s{{end}}{{end}}%s
    {{range $index, $author := .Authors}}{{if $index}}
    {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

%s{{range .VisibleCategories}}{{if .Name}}
    {{.Name}}:{{end}}{{range .VisibleCommands}}
    %s{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

%s
    {{range $index, $option := .VisibleFlags}}{{if $index}}
    {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

%s{{end}}
`, green("{{.Name}}"),
		yellow("USAGE:"),
		yellow("DESCRIPTION:"),
		yellow("AUTHOR"),
		yellow("S"),
		yellow(":"),
		yellow("COMMANDS:"),
		green(`{{join .Names ", "}}`),
		yellow("GLOBAL OPTIONS:"),
		blue("{{.Copyright}}"))

	// Set the command help template
	cli.CommandHelpTemplate = fmt.Sprintf(`%s
    {{.HelpName}} - {{.Usage}}

%s
    {{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}

%s
    {{.Category}}{{end}}{{if .Description}}

%s
    {{.Description}}{{end}}{{if .VisibleFlags}}

%s
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`, yellow("NAME:"),
		yellow("USAGE:"),
		yellow("CATEGORY:"),
		yellow("DESCRIPTION:"),
		yellow("OPTIONS:"))

	// Set the subcommand help template
	cli.SubcommandHelpTemplate = fmt.Sprintf(`%s
    {{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

%s
    {{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

%s{{range .VisibleCategories}}{{if .Name}}
    {{.Name}}:{{end}}{{range .VisibleCommands}}
    {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
%s
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`, yellow("NAME:"),
		yellow("USAGE:"),
		yellow("COMMANDS:"),
		yellow("OPTIONS:"))
}
