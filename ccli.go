package ccli

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// ColorFunc is a function that colorizes a string.
type ColorFunc func(a ...any) string

// Option is a functional option for configuring a ccli application.
type Option func(*Options)

// WithBlue sets the color function for author names.
func WithBlue(fn ColorFunc) Option {
	return func(opts *Options) {
		opts.Blue = fn
	}
}

// WithCyan sets the color function for usage text.
func WithCyan(fn ColorFunc) Option {
	return func(opts *Options) {
		opts.Cyan = fn
	}
}

// WithGreen sets the color function for app/command names.
func WithGreen(fn ColorFunc) Option {
	return func(opts *Options) {
		opts.Green = fn
	}
}

// WithRed sets the color function for copyright text.
func WithRed(fn ColorFunc) Option {
	return func(opts *Options) {
		opts.Red = fn
	}
}

// WithYellow sets the color function for section headers.
func WithYellow(fn ColorFunc) Option {
	return func(opts *Options) {
		opts.Yellow = fn
	}
}

// WithDisable turns off all coloring.
func WithDisable() Option {
	return func(opts *Options) {
		opts.Disable = true
	}
}

// Options allows customizing the color scheme used in help output.
type Options struct {
	// Blue is used for author names. Defaults to color.FgBlue.
	Blue ColorFunc
	// Cyan is used for usage text. Defaults to color.FgCyan.
	Cyan ColorFunc
	// Green is used for app/command names. Defaults to color.FgGreen.
	Green ColorFunc
	// Red is used for copyright text. Defaults to color.FgRed.
	Red ColorFunc
	// Yellow is used for section headers. Defaults to color.FgYellow.
	Yellow ColorFunc
	// Disable turns off all coloring. When true, all color functions
	// are replaced with plain fmt.Sprint.
	Disable bool
}

//nolint:gochecknoglobals // cached defaults to avoid redundant allocations
var (
	cachedDefaults     Options
	cachedDefaultsOnce sync.Once
)

func resolveOptions(opts Options) Options {
	if opts.Disable {
		noColor := ColorFunc(fmt.Sprint)
		opts.Blue = noColor
		opts.Cyan = noColor
		opts.Green = noColor
		opts.Red = noColor
		opts.Yellow = noColor

		return opts
	}

	if opts.Blue != nil && opts.Cyan != nil && opts.Green != nil &&
		opts.Red != nil && opts.Yellow != nil {
		return opts
	}

	defaults := defaultOptions()
	opts.Blue = colorOrDefault(opts.Blue, defaults.Blue)
	opts.Cyan = colorOrDefault(opts.Cyan, defaults.Cyan)
	opts.Green = colorOrDefault(opts.Green, defaults.Green)
	opts.Red = colorOrDefault(opts.Red, defaults.Red)
	opts.Yellow = colorOrDefault(opts.Yellow, defaults.Yellow)

	return opts
}

func colorOrDefault(fn, fallback ColorFunc) ColorFunc {
	if fn != nil {
		return fn
	}

	return fallback
}

func defaultOptions() Options {
	cachedDefaultsOnce.Do(func() {
		cachedDefaults = Options{
			Blue:    color.New(color.FgBlue).SprintFunc(),
			Cyan:    color.New(color.FgCyan).SprintFunc(),
			Green:   color.New(color.FgGreen).SprintFunc(),
			Red:     color.New(color.FgRed).SprintFunc(),
			Yellow:  color.New(color.FgYellow).SprintFunc(),
			Disable: false,
		}
	})

	return cachedDefaults
}

// NewApp creates a new application with colored help output using default
// colors.
//
// This function sets cli.CommandHelpTemplate and cli.SubcommandHelpTemplate
// as package-level globals. The app help template is set per-app via
// CustomAppHelpTemplate to avoid global side effects for that template.
func NewApp() *cli.App {
	return NewAppWithOptions(defaultOptions())
}

// NewAppWith creates a new application with colored help output, configured
// by functional options. Unset colors fall back to defaults.
//
// This function sets cli.CommandHelpTemplate and cli.SubcommandHelpTemplate
// as package-level globals. The app help template is set per-app via
// CustomAppHelpTemplate to avoid global side effects for that template.
func NewAppWith(opts ...Option) *cli.App {
	options := Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return NewAppWithOptions(options)
}

// NewAppWithOptions creates a new application with colored help output using
// the provided color options. Any nil color function in opts falls back to
// the default color. Set Disable to true to turn off all coloring.
//
// This function sets cli.CommandHelpTemplate and cli.SubcommandHelpTemplate
// as package-level globals. The app help template is set per-app via
// CustomAppHelpTemplate to avoid global side effects for that template.
func NewAppWithOptions(opts Options) *cli.App {
	opts = resolveOptions(opts)

	app := cli.NewApp()

	if opts.Disable {
		app.Writer = os.Stdout
		app.ErrWriter = os.Stderr
	} else {
		app.Writer = color.Output
		app.ErrWriter = color.Error
	}

	app.CustomAppHelpTemplate = appHelpTemplate(opts)

	setCommandHelpTemplate(opts)
	setSubcommandHelpTemplate(opts)

	return app
}

func appHelpTemplate(opts Options) string {
	return fmt.Sprintf(
		`%s {{if .Version}}{{if not .HideVersion}}{{.Version}}{{end}}{{end}}
{{if .Usage}}{{.Usage}}{{end}}

%s
    %s {{if .VisibleFlags}}[global options]{{end}}`+
			`{{if .Commands}} command [command options]{{end}} `+
			`{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}`+
			`[arguments...]{{end}}{{end}}{{if .Description}}

%s
    {{.Description}}{{end}}{{if len .Authors}}

%s{{with $length := len .Authors}}`+
			`{{if ne 1 $length}}%s{{end}}{{end}}%s
    {{range $index, $author := .Authors}}{{if $index}}
    {{end}}%s{{end}}{{end}}{{if .VisibleCommands}}

%s{{range .VisibleCategories}}{{if .Name}}
    {{.Name}}:{{end}}{{range .VisibleCommands}}
    %s{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

%s
    {{range $index, $option := .VisibleFlags}}{{if $index}}
    {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

%s{{end}}
`, opts.Green("{{.Name}}"),
		opts.Yellow("USAGE:"),
		opts.Cyan("{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}"),
		opts.Yellow("DESCRIPTION:"),
		opts.Yellow("AUTHOR"),
		opts.Yellow("S"),
		opts.Yellow(":"),
		opts.Blue("{{$author}}"),
		opts.Yellow("COMMANDS:"),
		opts.Green(`{{join .Names ", "}}`),
		opts.Yellow("GLOBAL OPTIONS:"),
		opts.Red("{{.Copyright}}"),
	)
}

func setCommandHelpTemplate(opts Options) {
	cli.CommandHelpTemplate = fmt.Sprintf(`%s
    %s - {{.Usage}}

%s
    %s{{if .VisibleFlags}} [command options]{{end}} `+
		`{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}`+
		`[arguments...]{{end}}{{if .Category}}

%s
    {{.Category}}{{end}}{{if .Description}}

%s
    {{.Description}}{{end}}{{if .VisibleFlags}}

%s
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`, opts.Yellow("NAME:"),
		opts.Green("{{.HelpName}}"),
		opts.Yellow("USAGE:"),
		opts.Cyan("{{.HelpName}}"),
		opts.Yellow("CATEGORY:"),
		opts.Yellow("DESCRIPTION:"),
		opts.Yellow("OPTIONS:"),
	)
}

func setSubcommandHelpTemplate(opts Options) {
	cli.SubcommandHelpTemplate = fmt.Sprintf(`%s
    %s - `+
		`{{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

%s
    %s command{{if .VisibleFlags}} [command options]{{end}} `+
		`{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}`+
		`[arguments...]{{end}}

%s{{range .VisibleCategories}}{{if .Name}}
    {{.Name}}:{{end}}{{range .VisibleCommands}}
    %s{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
%s
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`, opts.Yellow("NAME:"),
		opts.Green("{{.HelpName}}"),
		opts.Yellow("USAGE:"),
		opts.Cyan("{{.HelpName}}"),
		opts.Yellow("COMMANDS:"),
		opts.Green(`{{join .Names ", "}}`),
		opts.Yellow("OPTIONS:"),
	)
}
