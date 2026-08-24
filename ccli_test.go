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

const (
	testAppName    = "testapp"
	testCustom     = "custom"
	testHelpFlag   = "--help"
	testGreetName  = "greet"
	testGreetUsage = "say hello"
	testParentName = "parent"
	testChildName  = "child"
	testChildUsage = "child command"
	testSubName    = "sub"
	testWorldName  = "world"
	testNameHeader = "NAME:"
)

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

		return testCustom
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

func TestNewCommandWithAllColorOptions(t *testing.T) {
	t.Parallel()

	var blueCalled, cyanCalled, redCalled bool

	cmd := ccli.NewCommandWith(
		ccli.WithBlue(func(_ ...any) string {
			blueCalled = true

			return "blue"
		}),
		ccli.WithCyan(func(_ ...any) string {
			cyanCalled = true

			return "cyan"
		}),
		ccli.WithGreen(func(_ ...any) string {
			return "green"
		}),
		ccli.WithRed(func(_ ...any) string {
			redCalled = true

			return "red"
		}),
		ccli.WithYellow(func(_ ...any) string {
			return "yellow"
		}),
	)
	if cmd == nil {
		t.Fatal("NewCommandWith returned nil")
	}

	if !blueCalled {
		t.Error("WithBlue color function should have been called")
	}

	if !cyanCalled {
		t.Error("WithCyan color function should have been called")
	}

	if !redCalled {
		t.Error("WithRed color function should have been called")
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

		return testCustom
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
	cmd.Version = exampleVersion
	cmd.Authors = []any{"Test Author <test@test.com>"}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, testHelpFlag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testAppName, "a test application", exampleVersion, "Test Author"} {
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
	cmd.Version = exampleVersion

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{"colorapp", testHelpFlag})
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
			Name:  testGreetName,
			Usage: testGreetUsage,
			Action: func(_ context.Context, _ *cli.Command) error {
				return nil
			},
		},
	}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, testGreetName, testHelpFlag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testGreetName, testGreetUsage, testNameHeader, "USAGE:"} {
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
			Name:  testParentName,
			Usage: "parent command",
			Commands: []*cli.Command{
				{
					Name:  testChildName,
					Usage: testChildUsage,
					Action: func(_ context.Context, _ *cli.Command) error {
						return nil
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, testParentName, testHelpFlag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testParentName, testChildName, testNameHeader} {
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

	customTpl := "custom template"

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{
			Name:               testCustom,
			CustomHelpTemplate: customTpl,
		},
		{
			Name: "default",
		},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].CustomHelpTemplate != customTpl {
		t.Error("Apply should not overwrite existing CustomHelpTemplate")
	}

	if cmd.Commands[1].CustomHelpTemplate == "" {
		t.Error("Apply should set CustomHelpTemplate on subcommands without one")
	}
}

func TestApplyPropagatesWriter(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{
			Name: testSubName,
			Commands: []*cli.Command{
				{Name: "nested"},
			},
		},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].Writer != color.Output {
		t.Error("Apply should propagate Writer to subcommands")
	}

	if cmd.Commands[0].ErrWriter != color.Error {
		t.Error("Apply should propagate ErrWriter to subcommands")
	}

	nested := cmd.Commands[0].Commands[0]
	if nested.Writer != color.Output {
		t.Error("Apply should propagate Writer to nested subcommands")
	}

	if nested.ErrWriter != color.Error {
		t.Error("Apply should propagate ErrWriter to nested subcommands")
	}
}

func TestApplyFallsBackToColorWriters(t *testing.T) {
	t.Parallel()

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{Name: testSubName},
		},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].Writer != color.Output {
		t.Error("Apply should fall back to color.Output when root Writer is nil")
	}

	if cmd.Commands[0].ErrWriter != color.Error {
		t.Error("Apply should fall back to color.Error when root ErrWriter is nil")
	}
}

func TestApplyPropagatesDisabledWriter(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWith(ccli.WithDisable())
	cmd.Commands = []*cli.Command{
		{Name: testSubName},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].Writer != os.Stdout {
		t.Error("Apply should propagate disabled root's os.Stdout to subcommands")
	}

	if cmd.Commands[0].ErrWriter != os.Stderr {
		t.Error("Apply should propagate disabled root's os.Stderr to subcommands")
	}
}

func TestApplyPreservesExistingWriter(t *testing.T) {
	t.Parallel()

	var custom bytes.Buffer

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{
			Name:      testCustom,
			Writer:    &custom,
			ErrWriter: &custom,
		},
		{
			Name: "default",
		},
	}

	ccli.Apply(cmd)

	if cmd.Commands[0].Writer != &custom {
		t.Error("Apply should not overwrite existing Writer")
	}

	if cmd.Commands[0].ErrWriter != &custom {
		t.Error("Apply should not overwrite existing ErrWriter")
	}

	if cmd.Commands[1].Writer != color.Output {
		t.Error("Apply should set Writer on subcommands without one")
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
		{Name: testSubName},
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

	if cmd.Commands[0].Writer != os.Stdout {
		t.Error("disabled Apply should propagate os.Stdout to subcommands")
	}

	if cmd.Commands[0].ErrWriter != os.Stderr {
		t.Error("disabled Apply should propagate os.Stderr to subcommands")
	}
}

func TestApplyWith(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommand()
	cmd.Commands = []*cli.Command{
		{Name: testSubName},
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
			Name:  testGreetName,
			Usage: testGreetUsage,
			Commands: []*cli.Command{
				{
					Name:  testWorldName,
					Usage: "greet the world",
				},
			},
		},
	}

	ccli.ApplyWithOptions(cmd, noColorOptions())

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(
		context.Background(),
		[]string{testAppName, testGreetName, testWorldName, testHelpFlag},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{testWorldName, "greet the world", testNameHeader, "USAGE:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("nested subcommand help missing %q, got:\n%s", expected, output)
		}
	}
}

func TestCommandHelpWithSubcommands(t *testing.T) {
	t.Parallel()

	cmd := ccli.NewCommandWithOptions(noColorOptions())
	cmd.Name = testAppName
	cmd.Commands = []*cli.Command{
		{
			Name:  testParentName,
			Usage: "parent command",
			Commands: []*cli.Command{
				{
					Name:  testChildName,
					Usage: testChildUsage,
				},
			},
		},
	}

	ccli.ApplyWithOptions(cmd, noColorOptions())

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, testParentName, testHelpFlag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"COMMANDS:", testChildName, testChildUsage} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output missing %q, got:\n%s", expected, output)
		}
	}
}

func TestCopyrightUsesRedColor(t *testing.T) {
	t.Parallel()

	red := color.New(color.FgRed)
	red.EnableColor()

	cmd := ccli.NewCommandWithOptions(ccli.Options{
		Blue:    ccli.ColorFunc(fmt.Sprint),
		Cyan:    ccli.ColorFunc(fmt.Sprint),
		Green:   ccli.ColorFunc(fmt.Sprint),
		Red:     red.SprintFunc(),
		Yellow:  ccli.ColorFunc(fmt.Sprint),
		Disable: false,
	})
	cmd.Name = testAppName
	cmd.Copyright = "2026 Test Corp"

	var buf bytes.Buffer

	cmd.Writer = &buf

	err := cmd.Run(context.Background(), []string{testAppName, testHelpFlag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "2026 Test Corp") {
		t.Errorf("help output missing copyright text, got:\n%s", output)
	}

	if !strings.Contains(output, "\033[31m") {
		t.Errorf("COPYRIGHT header should use red ANSI color, got:\n%s", output)
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
