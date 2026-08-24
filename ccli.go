package ccli

import (
	"fmt"
	"io"
	"os"
	"strings"
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

//nolint:gochecknoglobals // cached to avoid redundant allocations; color state is fixed after init
var (
	cachedDefaults     Options
	cachedDefaultsOnce sync.Once
	cachedRootTpl      string
	cachedCmdTpl       string
	cachedTplOnce      sync.Once
)

func buildOptions(opts []Option) Options {
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

	return options
}

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

//nolint:gocritic // named returns would conflict with nonamedreturns linter
func defaultTemplates() (string, string) {
	cachedTplOnce.Do(func() {
		opts := defaultOptions()
		cachedRootTpl = rootCommandHelpTemplate(opts)
		cachedCmdTpl = commandHelpTemplate(opts)
	})

	return cachedRootTpl, cachedCmdTpl
}

// NewCommand creates a new root command with colored help output using
// default colors. All help templates are set per-command via
// CustomRootCommandHelpTemplate and CustomHelpTemplate, avoiding
// global side effects.
func NewCommand() *cli.Command {
	rootTpl, cmdTpl := defaultTemplates()

	return &cli.Command{
		Writer:                        color.Output,
		ErrWriter:                     color.Error,
		CustomRootCommandHelpTemplate: rootTpl,
		CustomHelpTemplate:            cmdTpl,
	}
}

// NewCommandWith creates a new root command with colored help output,
// configured by functional options. Unset colors fall back to defaults.
// All help templates are set per-command, avoiding global side effects.
func NewCommandWith(opts ...Option) *cli.Command {
	return NewCommandWithOptions(buildOptions(opts))
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

// Apply recursively sets colored help templates on all subcommands of cmd
// using default colors. Call this after adding subcommands to ensure they
// get colored help output.
func Apply(cmd *cli.Command) {
	_, cmdTpl := defaultTemplates()

	writer := cmd.Writer
	if writer == nil {
		writer = color.Output
	}

	errWriter := cmd.ErrWriter
	if errWriter == nil {
		errWriter = color.Error
	}

	applyTemplates(cmd.Commands, cmdTpl, writer, errWriter)
}

// ApplyWith recursively sets colored help templates on all subcommands of cmd
// using functional options. Call this after adding subcommands.
func ApplyWith(cmd *cli.Command, opts ...Option) {
	ApplyWithOptions(cmd, buildOptions(opts))
}

// ApplyWithOptions recursively sets colored help templates on all subcommands
// of cmd using the provided color options. Call this after adding subcommands.
func ApplyWithOptions(cmd *cli.Command, opts Options) {
	opts = resolveOptions(opts)
	tpl := commandHelpTemplate(opts)

	var (
		writer    = color.Output
		errWriter = color.Error
	)

	if opts.Disable {
		writer = os.Stdout
		errWriter = os.Stderr
	}

	applyTemplates(cmd.Commands, tpl, writer, errWriter)
}

func applyTemplates(cmds []*cli.Command, tpl string, writer, errWriter io.Writer) {
	for _, sub := range cmds {
		if sub.CustomHelpTemplate == "" {
			sub.CustomHelpTemplate = tpl
		}

		if sub.Writer == nil {
			sub.Writer = writer
		}

		if sub.ErrWriter == nil {
			sub.ErrWriter = errWriter
		}

		applyTemplates(sub.Commands, tpl, writer, errWriter)
	}
}

func rootCommandHelpTemplate(opts Options) string {
	var buf strings.Builder

	buf.WriteString(opts.Yellow("NAME:"))
	buf.WriteString("\n   {{$v := offset .FullName 6}}")
	buf.WriteString(opts.Green("{{wrap .FullName 3}}"))
	buf.WriteString("{{if .Usage}} - {{wrap .Usage $v}}{{end}}\n\n")

	buf.WriteString(opts.Yellow("USAGE:"))
	buf.WriteString("\n   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}")
	buf.WriteString(opts.Cyan("{{.FullName}}"))
	buf.WriteString(" {{if .VisibleFlags}}[global options]{{end}}")
	buf.WriteString("{{if .VisibleCommands}} [command [command options]]{{end}} ")
	buf.WriteString("{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}")
	buf.WriteString("{{if .Arguments}}[arguments...]{{end}}{{end}}{{end}}")

	buf.WriteString("{{if .Version}}{{if not .HideVersion}}\n\n")
	buf.WriteString(opts.Yellow("VERSION:"))
	buf.WriteString("\n   {{.Version}}{{end}}{{end}}")

	buf.WriteString("{{if .Description}}\n\n")
	buf.WriteString(opts.Yellow("DESCRIPTION:"))
	buf.WriteString("\n   {{template \"descriptionTemplate\" .}}{{end}}")

	buf.WriteString("\n{{- if len .Authors}}\n\n")
	buf.WriteString(opts.Yellow("AUTHOR"))
	buf.WriteString("{{with $length := len .Authors}}{{if ne 1 $length}}")
	buf.WriteString(opts.Yellow("S"))
	buf.WriteString("{{end}}{{end}}")
	buf.WriteString(opts.Yellow(":"))
	buf.WriteString("\n   {{range $index, $author := .Authors}}{{if $index}}\n   {{end}}")
	buf.WriteString(opts.Blue("{{$author}}"))
	buf.WriteString("{{end}}{{end}}")

	buf.WriteString("{{if .VisibleCommands}}\n\n")
	buf.WriteString(opts.Yellow("COMMANDS:"))
	buf.WriteString("{{template \"visibleCommandCategoryTemplate\" .}}{{end}}")

	buf.WriteString("{{if .VisibleFlagCategories}}\n\n")
	buf.WriteString(opts.Yellow("GLOBAL OPTIONS:"))
	buf.WriteString("{{template \"visibleFlagCategoryTemplate\" .}}")
	buf.WriteString("{{else if .VisibleFlags}}\n\n")
	buf.WriteString(opts.Yellow("GLOBAL OPTIONS:"))
	buf.WriteString("{{template \"visibleFlagTemplate\" .}}{{end}}")

	buf.WriteString("{{if .Copyright}}\n\n")
	buf.WriteString(opts.Red("COPYRIGHT:"))
	buf.WriteString("\n   {{template \"copyrightTemplate\" .}}{{end}}\n")

	return buf.String()
}

func commandHelpTemplate(opts Options) string {
	var buf strings.Builder

	buf.WriteString(opts.Yellow("NAME:"))
	buf.WriteString("\n   {{$v := offset .FullName 6}}")
	buf.WriteString(opts.Green("{{wrap .FullName 3}}"))
	buf.WriteString("{{if .Usage}} - {{wrap .Usage $v}}{{end}}\n\n")

	buf.WriteString(opts.Yellow("USAGE:"))
	buf.WriteString("\n   {{template \"usageTemplate\" .}}")

	buf.WriteString("{{if .Category}}\n\n")
	buf.WriteString(opts.Yellow("CATEGORY:"))
	buf.WriteString("\n   {{.Category}}{{end}}")

	buf.WriteString("{{if .Description}}\n\n")
	buf.WriteString(opts.Yellow("DESCRIPTION:"))
	buf.WriteString("\n   {{template \"descriptionTemplate\" .}}{{end}}")

	buf.WriteString("{{if .VisibleCommands}}\n\n")
	buf.WriteString(opts.Yellow("COMMANDS:"))
	buf.WriteString("{{template \"visibleCommandCategoryTemplate\" .}}{{end}}")

	buf.WriteString("{{if .VisibleFlagCategories}}\n\n")
	buf.WriteString(opts.Yellow("OPTIONS:"))
	buf.WriteString("{{template \"visibleFlagCategoryTemplate\" .}}")
	buf.WriteString("{{else if .VisibleFlags}}\n\n")
	buf.WriteString(opts.Yellow("OPTIONS:"))
	buf.WriteString("{{template \"visibleFlagTemplate\" .}}{{end}}")

	buf.WriteString("{{if .VisiblePersistentFlags}}\n\n")
	buf.WriteString(opts.Yellow("GLOBAL OPTIONS:"))
	buf.WriteString("{{template \"visiblePersistentFlagTemplate\" .}}{{end}}\n")

	return buf.String()
}
