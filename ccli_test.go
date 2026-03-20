package ccli_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/saschagrunert/ccli/v2"
	"github.com/urfave/cli/v2"
)

const testAppName = "testapp"

func noColorOptions() ccli.Options {
	return ccli.Options{
		Blue:    ccli.ColorFunc(fmt.Sprint),
		Cyan:    ccli.ColorFunc(fmt.Sprint),
		Green:   ccli.ColorFunc(fmt.Sprint),
		Red:     ccli.ColorFunc(fmt.Sprint),
		Yellow:  ccli.ColorFunc(fmt.Sprint),
		Disable: false,
	}
}

func TestNewApp(t *testing.T) {
	t.Parallel()

	app := ccli.NewApp()
	if app == nil {
		t.Fatal("NewApp returned nil")
	}

	if app.Writer != color.Output {
		t.Error("Writer should be color.Output")
	}

	if app.ErrWriter != color.Error {
		t.Error("ErrWriter should be color.Error")
	}

	if app.CustomAppHelpTemplate == "" {
		t.Error("CustomAppHelpTemplate should be set")
	}
}

func TestNewAppWith(t *testing.T) {
	t.Parallel()

	called := false
	custom := ccli.ColorFunc(func(_ ...any) string {
		called = true

		return "custom"
	})

	app := ccli.NewAppWith(
		ccli.WithGreen(custom),
		ccli.WithYellow(custom),
	)
	if app == nil {
		t.Fatal("NewAppWith returned nil")
	}

	if !called {
		t.Error("custom color function should have been called for templates")
	}
}

func TestNewAppWithDisable(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWith(ccli.WithDisable())
	if app == nil {
		t.Fatal("NewAppWith WithDisable returned nil")
	}

	if app.Writer != os.Stdout {
		t.Error("Writer should be os.Stdout when disabled")
	}

	if app.ErrWriter != os.Stderr {
		t.Error("ErrWriter should be os.Stderr when disabled")
	}
}

func TestNewAppWithOptions(t *testing.T) {
	t.Parallel()

	called := false
	custom := ccli.ColorFunc(func(_ ...any) string {
		called = true

		return "custom"
	})

	app := ccli.NewAppWithOptions(ccli.Options{
		Blue:    custom,
		Cyan:    custom,
		Green:   custom,
		Red:     custom,
		Yellow:  custom,
		Disable: false,
	})
	if app == nil {
		t.Fatal("NewAppWithOptions returned nil")
	}

	if !called {
		t.Error("custom color function should have been called for templates")
	}
}

func TestNewAppWithOptionsNilFallback(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: false,
	})
	if app == nil {
		t.Fatal("NewAppWithOptions returned nil with empty options")
	}

	if app.CustomAppHelpTemplate == "" {
		t.Error("CustomAppHelpTemplate should be set even with empty options")
	}
}

func TestNewAppWithOptionsPartialOverride(t *testing.T) {
	t.Parallel()

	greenCalled := false
	customGreen := ccli.ColorFunc(func(_ ...any) string {
		greenCalled = true

		return "green"
	})

	app := ccli.NewAppWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   customGreen,
		Red:     nil,
		Yellow:  nil,
		Disable: false,
	})
	if app == nil {
		t.Fatal("NewAppWithOptions returned nil")
	}

	if !greenCalled {
		t.Error("custom green function should have been called")
	}

	if app.CustomAppHelpTemplate == "" {
		t.Error("CustomAppHelpTemplate should be set with partial options")
	}
}

func TestNewAppWithOptionsDisable(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: true,
	})
	if app == nil {
		t.Fatal("NewAppWithOptions returned nil with Disable")
	}

	if app.Writer != os.Stdout {
		t.Error("Writer should be os.Stdout when disabled")
	}

	if app.ErrWriter != os.Stderr {
		t.Error("ErrWriter should be os.Stderr when disabled")
	}

	if !strings.Contains(app.CustomAppHelpTemplate, "USAGE:") {
		t.Error("disabled template should still contain USAGE section")
	}

	// Verify no ANSI escape codes in the template.
	if strings.Contains(app.CustomAppHelpTemplate, "\033[") {
		t.Error("disabled template should not contain ANSI escape codes")
	}
}

func TestAppHelpOutput(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(noColorOptions())
	app.Name = testAppName
	app.Usage = "a test application"
	app.Version = "1.0.0"
	app.Authors = []*cli.Author{{Name: "Test Author", Email: "test@test.com"}}

	var buf bytes.Buffer

	app.Writer = &buf

	err := app.Run([]string{testAppName, "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testAppName, "a test application", "1.0.0", "Test Author"} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output missing %q", expected)
		}
	}
}

func TestAppHelpRenderedColorCodes(t *testing.T) {
	t.Parallel()

	green := color.New(color.FgGreen)
	green.EnableColor()

	yellow := color.New(color.FgYellow)
	yellow.EnableColor()

	app := ccli.NewAppWithOptions(ccli.Options{
		Blue:    color.New(color.FgBlue).SprintFunc(),
		Cyan:    color.New(color.FgCyan).SprintFunc(),
		Green:   green.SprintFunc(),
		Red:     color.New(color.FgRed).SprintFunc(),
		Yellow:  yellow.SprintFunc(),
		Disable: false,
	})
	app.Name = "colorapp"
	app.Usage = "a colored app"
	app.Version = "1.0.0"

	var buf bytes.Buffer

	app.Writer = &buf

	err := app.Run([]string{"colorapp", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Verify rendered output contains ANSI escape codes.
	if !strings.Contains(output, "\033[") {
		t.Error("rendered help output should contain ANSI escape codes")
	}
}

func TestAppHelpInvalidFlag(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(noColorOptions())
	app.Name = testAppName

	var errBuf bytes.Buffer

	app.ErrWriter = &errBuf

	err := app.Run([]string{testAppName, "--nonexistent"})
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
}

func TestCommandHelpOutput(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(noColorOptions())
	app.Name = testAppName
	app.Commands = []*cli.Command{
		{
			Name:  "greet",
			Usage: "say hello",
			Action: func(_ *cli.Context) error {
				return nil
			},
		},
	}

	var buf bytes.Buffer

	app.Writer = &buf

	err := app.Run([]string{testAppName, "greet", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"greet", "say hello", "NAME:", "USAGE:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("command help output missing %q", expected)
		}
	}
}

func TestSubcommandHelpOutput(t *testing.T) {
	t.Parallel()

	app := ccli.NewAppWithOptions(noColorOptions())
	app.Name = testAppName
	app.Commands = []*cli.Command{
		{
			Name:  "parent",
			Usage: "parent command",
			Subcommands: []*cli.Command{
				{
					Name:  "child",
					Usage: "child command",
					Action: func(_ *cli.Context) error {
						return nil
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	app.Writer = &buf

	err := app.Run([]string{testAppName, "parent", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"parent", "child", "COMMANDS:", "NAME:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("subcommand help output missing %q", expected)
		}
	}
}

func TestCommandHelpTemplate(t *testing.T) {
	t.Parallel()

	_ = ccli.NewApp()

	if cli.CommandHelpTemplate == "" {
		t.Error("CommandHelpTemplate should be set after NewApp")
	}
}

func TestSubcommandHelpTemplate(t *testing.T) {
	t.Parallel()

	_ = ccli.NewApp()

	if cli.SubcommandHelpTemplate == "" {
		t.Error("SubcommandHelpTemplate should be set after NewApp")
	}
}

func BenchmarkNewApp(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = ccli.NewApp()
	}
}

func BenchmarkNewAppWithOptions(b *testing.B) {
	b.ReportAllocs()

	opts := ccli.Options{
		Blue:    ccli.ColorFunc(fmt.Sprint),
		Cyan:    ccli.ColorFunc(fmt.Sprint),
		Green:   ccli.ColorFunc(fmt.Sprint),
		Red:     ccli.ColorFunc(fmt.Sprint),
		Yellow:  ccli.ColorFunc(fmt.Sprint),
		Disable: false,
	}

	for b.Loop() {
		_ = ccli.NewAppWithOptions(opts)
	}
}

func BenchmarkNewAppWithDisable(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = ccli.NewAppWith(ccli.WithDisable())
	}
}
