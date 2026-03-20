package ccli_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/saschagrunert/ccli/v3"
	"github.com/urfave/cli/v3"
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

func TestNewCommand(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommand()
	if cmd == nil {
		t.Fatal("NewCommand returned nil")
	}

	if cmd.Writer != color.Output {
		t.Error("Writer should be color.Output")
	}

	if cmd.ErrWriter != color.Error {
		t.Error("ErrWriter should be color.Error")
	}

	if cmd.CustomRootCommandHelpTemplate == "" {
		t.Error("CustomRootCommandHelpTemplate should be set")
	}

	if cmd.CustomHelpTemplate == "" {
		t.Error("CustomHelpTemplate should be set")
	}
}

func TestNewCommandWith(t *testing.T) {
	t.Parallel()

	called := false
	custom := ccli.ColorFunc(func(_ ...any) string {
		called = true

		return "custom"
	})

	cmd := ccli.NewCommandWith(
		ccli.WithGreen(custom),
		ccli.WithYellow(custom),
	)
	if cmd == nil {
		t.Fatal("NewCommandWith returned nil")
	}

	if !called {
		t.Error("custom color function should have been called for templates")
	}
}

func TestNewCommandWithDisable(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWith(ccli.WithDisable())
	if cmd == nil {
		t.Fatal("NewCommandWith WithDisable returned nil")
	}

	if cmd.Writer != os.Stdout {
		t.Error("Writer should be os.Stdout when disabled")
	}

	if cmd.ErrWriter != os.Stderr {
		t.Error("ErrWriter should be os.Stderr when disabled")
	}
}

func TestNewCommandWithOptions(t *testing.T) {
	t.Parallel()

	called := false
	custom := ccli.ColorFunc(func(_ ...any) string {
		called = true

		return "custom"
	})

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    custom,
		Cyan:    custom,
		Green:   custom,
		Red:     custom,
		Yellow:  custom,
		Disable: false,
	})
	if cmd == nil {
		t.Fatal("NewCommandWithOptions returned nil")
	}

	if !called {
		t.Error("custom color function should have been called for templates")
	}
}

func TestNewCommandWithOptionsNilFallback(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: false,
	})
	if cmd == nil {
		t.Fatal("NewCommandWithOptions returned nil with empty options")
	}

	if cmd.CustomRootCommandHelpTemplate == "" {
		t.Error("CustomRootCommandHelpTemplate should be set even with empty options")
	}
}

func TestNewCommandWithOptionsPartialOverride(t *testing.T) {
	t.Parallel()

	greenCalled := false
	customGreen := ccli.ColorFunc(func(_ ...any) string {
		greenCalled = true

		return "green"
	})

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   customGreen,
		Red:     nil,
		Yellow:  nil,
		Disable: false,
	})
	if cmd == nil {
		t.Fatal("NewCommandWithOptions returned nil")
	}

	if !greenCalled {
		t.Error("custom green function should have been called")
	}

	if cmd.CustomRootCommandHelpTemplate == "" {
		t.Error("CustomRootCommandHelpTemplate should be set with partial options")
	}
}

func TestNewCommandWithOptionsDisable(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: true,
	})
	if cmd == nil {
		t.Fatal("NewCommandWithOptions returned nil with Disable")
	}

	if cmd.Writer != os.Stdout {
		t.Error("Writer should be os.Stdout when disabled")
	}

	if cmd.ErrWriter != os.Stderr {
		t.Error("ErrWriter should be os.Stderr when disabled")
	}

	if !strings.Contains(cmd.CustomRootCommandHelpTemplate, "USAGE:") {
		t.Error("disabled template should still contain USAGE section")
	}

	// Verify no ANSI escape codes in the template.
	if strings.Contains(cmd.CustomRootCommandHelpTemplate, "\033[") {
		t.Error("disabled template should not contain ANSI escape codes")
	}
}

func TestCommandHelpOutput(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName
	cmd.Usage = "a test application"
	cmd.Version = "1.0.0"
	cmd.Authors = []any{"Test Author <test@test.com>"}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testAppName, "a test application", "1.0.0", "Test Author"} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output missing %q, got:\n%s", expected, output)
		}
	}
}

func TestCommandHelpRenderedColorCodes(t *testing.T) {
	t.Parallel()

	green := color.New(color.FgGreen)
	green.EnableColor()

	yellow := color.New(color.FgYellow)
	yellow.EnableColor()

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    color.New(color.FgBlue).SprintFunc(),
		Cyan:    color.New(color.FgCyan).SprintFunc(),
		Green:   green.SprintFunc(),
		Red:     color.New(color.FgRed).SprintFunc(),
		Yellow:  yellow.SprintFunc(),
		Disable: false,
	})
	cmd.Name = "colorapp"
	cmd.Usage = "a colored app"
	cmd.Version = "1.0.0"

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{"colorapp", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Verify rendered output contains ANSI escape codes.
	if !strings.Contains(output, "\033[") {
		t.Error("rendered help output should contain ANSI escape codes")
	}
}

func TestCommandHelpInvalidFlag(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName

	var errBuf bytes.Buffer

	cmd.ErrWriter = &errBuf

	err := cmd.Run(context.Background(), []string{testAppName, "--nonexistent"})
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
}

func TestSubcommandHelpOutput(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName
	cmd.Commands = []*cli.Command{
		{
			Name:  "greet",
			Usage: "say hello",
			Action: func(_ context.Context, _ *cli.Command) error {
				return nil
			},
		},
	}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, "greet", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"greet", "say hello", "NAME:", "USAGE:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("subcommand help output missing %q, got:\n%s", expected, output)
		}
	}
}

func TestNestedSubcommandHelpOutput(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName
	cmd.Commands = []*cli.Command{
		{
			Name:  "parent",
			Usage: "parent command",
			Commands: []*cli.Command{
				{
					Name:  "child",
					Usage: "child command",
					Action: func(_ context.Context, _ *cli.Command) error {
						return nil
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, "parent", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"parent", "child", "NAME:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("nested subcommand help output missing %q, got:\n%s", expected, output)
		}
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommand()
	cmd.Name = testAppName
	cmd.Commands = []*cli.Command{
		{
			Name:  "sub1",
			Usage: "first subcommand",
			Commands: []*cli.Command{
				{
					Name:  "nested",
					Usage: "nested subcommand",
				},
			},
		},
		{
			Name:  "sub2",
			Usage: "second subcommand",
		},
	}

	ccli.Apply(cmd)

	for _, sub := range cmd.Commands {
		if sub.CustomHelpTemplate == "" {
			t.Errorf("subcommand %q should have CustomHelpTemplate set", sub.Name)
		}
	}

	nested := cmd.Commands[0].Commands[0]
	if nested.CustomHelpTemplate == "" {
		t.Error("nested subcommand should have CustomHelpTemplate set")
	}
}

func TestApplyPreservesExisting(t *testing.T) {
	t.Parallel()

	custom := "custom template"

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{
			Name:               "custom",
			CustomHelpTemplate: custom,
		},
		{
			Name: "default",
		},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].CustomHelpTemplate != custom {
		t.Error("Apply should not overwrite existing CustomHelpTemplate")
	}

	if cmd.Commands[1].CustomHelpTemplate == "" {
		t.Error("Apply should set CustomHelpTemplate on subcommands without one")
	}
}

func TestApplyWithOptions(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: true,
	})
	cmd.Commands = []*cli.Command{
		{Name: "sub"},
	}

	ccli.ApplyWithOptions(cmd, ccli.Options{
		Blue:    nil,
		Cyan:    nil,
		Green:   nil,
		Red:     nil,
		Yellow:  nil,
		Disable: true,
	})

	tpl := cmd.Commands[0].CustomHelpTemplate
	if tpl == "" {
		t.Fatal("subcommand should have CustomHelpTemplate set")
	}

	if strings.Contains(tpl, "\033[") {
		t.Error("disabled template should not contain ANSI escape codes")
	}
}

func TestApplyWith(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{Name: "sub"},
	}

	ccli.ApplyWith(cmd, ccli.WithDisable())

	tpl := cmd.Commands[0].CustomHelpTemplate
	if tpl == "" {
		t.Fatal("subcommand should have CustomHelpTemplate set")
	}

	if strings.Contains(tpl, "\033[") {
		t.Error("disabled template should not contain ANSI escape codes")
	}
}

func TestSubcommandHelpOutputWithApply(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName
	cmd.Commands = []*cli.Command{
		{
			Name:  "greet",
			Usage: "say hello",
			Commands: []*cli.Command{
				{
					Name:  "world",
					Usage: "greet the world",
				},
			},
		},
	}

	ccli.ApplyWithOptions(cmd, noColorOptions())

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, "greet", "world", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"world", "greet the world", "NAME:", "USAGE:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("nested subcommand help missing %q, got:\n%s", expected, output)
		}
	}
}

func TestNoGlobalTemplateMutation(t *testing.T) {
	t.Parallel()

	originalRoot := cli.RootCommandHelpTemplate
	originalCmd := cli.CommandHelpTemplate
	originalSub := cli.SubcommandHelpTemplate

	_ = ccli.NewCommand()

	if cli.RootCommandHelpTemplate != originalRoot {
		t.Error("NewCommand should not mutate cli.RootCommandHelpTemplate")
	}

	if cli.CommandHelpTemplate != originalCmd {
		t.Error("NewCommand should not mutate cli.CommandHelpTemplate")
	}

	if cli.SubcommandHelpTemplate != originalSub {
		t.Error("NewCommand should not mutate cli.SubcommandHelpTemplate")
	}
}

func BenchmarkNewCommand(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = ccli.NewCommand()
	}
}

func BenchmarkNewCommandWithOptions(b *testing.B) {
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
		_ = ccli.NewCommandWithOptions(opts)
	}
}

func BenchmarkNewCommandWithDisable(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_ = ccli.NewCommandWith(ccli.WithDisable())
	}
}
