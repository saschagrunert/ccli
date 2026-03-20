package ccli

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
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

// NewCommand creates a new root command with colored help output using
// default colors. All help templates are set per-command via
// CustomRootCommandHelpTemplate and CustomHelpTemplate, avoiding
// global side effects.
func NewCommand() *cli.Command {
	return NewCommandWithOptions(defaultOptions())
}

// NewCommandWith creates a new root command with colored help output,
// configured by functional options. Unset colors fall back to defaults.
// All help templates are set per-command, avoiding global side effects.
func NewCommandWith(opts ...Option) *cli.Command {
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

	return NewCommandWithOptions(options)
}

// NewCommandWithOptions creates a new root command with colored help output
// using the provided color options. Any nil color function in opts falls
// back to the default color. Set Disable to true to turn off all coloring.
// All help templates are set per-command, avoiding global side effects.
func NewCommandWithOptions(opts Options) *cli.Command {
	opts = resolveOptions(opts)

	cmd := &cli.Command{}

	if opts.Disable {
		cmd.Writer = os.Stdout
		cmd.ErrWriter = os.Stderr
	} else {
		cmd.Writer = color.Output
		cmd.ErrWriter = color.Error
	}

	cmd.CustomRootCommandHelpTemplate = rootCommandHelpTemplate(opts)
	cmd.CustomHelpTemplate = commandHelpTemplate(opts)

	return cmd
}

func rootCommandHelpTemplate(opts Options) string {
	return fmt.Sprintf(
		`%s
   {{$v := offset .FullName 6}}%s{{if .Usage}} - {{wrap .Usage $v}}{{end}}

%s
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}%s `+
			`{{if .VisibleFlags}}[global options]{{end}}`+
			`{{if .VisibleCommands}} [command [command options]]{{end}} `+
			`{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}`+
			`{{if .Arguments}}[arguments...]{{end}}{{end}}{{end}}`+
			`{{if .Version}}{{if not .HideVersion}}

%s
   {{.Version}}{{end}}{{end}}{{if .Description}}

%s
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

%s{{with $length := len .Authors}}`+
			`{{if ne 1 $length}}%s{{end}}{{end}}%s
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}%s{{end}}{{end}}{{if .VisibleCommands}}

%s{{template "visibleCommandCategoryTemplate" .}}{{end}}`+
			`{{if .VisibleFlagCategories}}

%s{{template "visibleFlagCategoryTemplate" .}}`+
			`{{else if .VisibleFlags}}

%s{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}

%s
   {{template "copyrightTemplate" .}}{{end}}
`, opts.Yellow("NAME:"),
		opts.Green("{{wrap .FullName 3}}"),
		opts.Yellow("USAGE:"),
		opts.Cyan("{{.FullName}}"),
		opts.Yellow("VERSION:"),
		opts.Yellow("DESCRIPTION:"),
		opts.Yellow("AUTHOR"),
		opts.Yellow("S"),
		opts.Yellow(":"),
		opts.Blue("{{$author}}"),
		opts.Yellow("COMMANDS:"),
		opts.Yellow("GLOBAL OPTIONS:"),
		opts.Yellow("GLOBAL OPTIONS:"),
		opts.Yellow("COPYRIGHT:"),
	)
}

func commandHelpTemplate(opts Options) string {
	return fmt.Sprintf(
		`%s
   {{$v := offset .FullName 6}}%s{{if .Usage}} - {{wrap .Usage $v}}{{end}}

%s
   {{template "usageTemplate" .}}{{if .Category}}

%s
   {{.Category}}{{end}}{{if .Description}}

%s
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

%s{{template "visibleFlagCategoryTemplate" .}}`+
			`{{else if .VisibleFlags}}

%s{{template "visibleFlagTemplate" .}}{{end}}`+
			`{{if .VisiblePersistentFlags}}

%s{{template "visiblePersistentFlagTemplate" .}}{{end}}
`, opts.Yellow("NAME:"),
		opts.Green("{{wrap .FullName 3}}"),
		opts.Yellow("USAGE:"),
		opts.Yellow("CATEGORY:"),
		opts.Yellow("DESCRIPTION:"),
		opts.Yellow("OPTIONS:"),
		opts.Yellow("OPTIONS:"),
		opts.Yellow("GLOBAL OPTIONS:"),
	)
}
