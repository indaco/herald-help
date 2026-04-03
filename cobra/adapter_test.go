package heraldcobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTestCommand() *cobra.Command {
	root := &cobra.Command{
		Use:        "myapp [flags] <command>",
		Short:      "A sample app",
		Long:       "A sample application for testing the cobra adapter.",
		Aliases:    []string{"ma"},
		Deprecated: "use newapp instead",
		Example: `  # Start the server
  $ myapp serve --port 9090

  # Build with verbose output
  $ myapp build -v`,
	}

	root.Flags().StringP("output", "o", "stdout", "Output destination")
	root.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	root.Flags().Int("port", 8080, "Port number")

	root.PersistentFlags().String("config", "", "Config file")

	serve := &cobra.Command{
		Use:     "serve",
		Short:   "Start the server",
		Aliases: []string{"s"},
	}
	build := &cobra.Command{
		Use:   "build",
		Short: "Build the project",
	}
	hidden := &cobra.Command{
		Use:    "internal",
		Short:  "Internal command",
		Hidden: true,
	}

	root.AddCommand(serve, build, hidden)

	return root
}

func TestFromCobra(t *testing.T) {
	cmd := newTestCommand()
	hc := FromCobra(cmd)

	if hc.Name != "myapp" {
		t.Errorf("Name = %q, want %q", hc.Name, "myapp")
	}

	if hc.Description != "A sample application for testing the cobra adapter." {
		t.Errorf("Description = %q, want long description", hc.Description)
	}

	if len(hc.Aliases) != 1 || hc.Aliases[0] != "ma" {
		t.Errorf("Aliases = %v, want [ma]", hc.Aliases)
	}

	if hc.Deprecated != "use newapp instead" {
		t.Errorf("Deprecated = %q, want %q", hc.Deprecated, "use newapp instead")
	}
}

func TestFromCobraFlags(t *testing.T) {
	cmd := newTestCommand()
	hc := FromCobra(cmd)

	flagMap := make(map[string]struct{ Short, Type, Default string })
	for _, f := range hc.Flags {
		flagMap[f.Long] = struct{ Short, Type, Default string }{f.Short, f.Type, f.Default}
	}

	t.Run("output flag", func(t *testing.T) {
		f, ok := flagMap["--output"]
		if !ok {
			t.Fatal("--output not found")
		}
		if f.Short != "-o" {
			t.Errorf("output short = %q, want %q", f.Short, "-o")
		}
		if f.Type != "string" {
			t.Errorf("output type = %q, want %q", f.Type, "string")
		}
	})

	t.Run("verbose flag", func(t *testing.T) {
		f, ok := flagMap["--verbose"]
		if !ok {
			t.Fatal("--verbose not found")
		}
		if f.Short != "-v" {
			t.Errorf("verbose short = %q, want %q", f.Short, "-v")
		}
	})

	t.Run("port flag", func(t *testing.T) {
		f, ok := flagMap["--port"]
		if !ok {
			t.Fatal("--port not found")
		}
		if f.Short != "" {
			t.Errorf("port should have no short, got %q", f.Short)
		}
	})
}

func TestFromCobraInheritedFlags(t *testing.T) {
	root := newTestCommand()
	// Access a child to trigger persistent flag inheritance.
	child := root.Commands()[0] // serve
	hc := FromCobra(child)

	var found bool
	for _, f := range hc.Flags {
		if f.Long == "--config" && f.Inherited {
			found = true
			break
		}
	}

	if !found {
		t.Error("inherited --config flag not found")
	}
}

func TestFromCobraSubcommands(t *testing.T) {
	cmd := newTestCommand()
	hc := FromCobra(cmd)

	// Should include serve, build, and help (auto-added), but not hidden.
	names := make(map[string]bool)
	for _, c := range hc.Commands {
		names[c.Name] = true
	}

	if !names["serve"] {
		t.Error("missing serve command")
	}
	if !names["build"] {
		t.Error("missing build command")
	}
	if names["internal"] {
		t.Error("hidden command should not appear")
	}
}

func TestFromCobraExamples(t *testing.T) {
	cmd := newTestCommand()
	hc := FromCobra(cmd)

	if len(hc.Examples) != 2 {
		t.Fatalf("Examples count = %d, want 2", len(hc.Examples))
	}

	if hc.Examples[0].Desc != "# Start the server" {
		t.Errorf("first example desc = %q", hc.Examples[0].Desc)
	}
	if hc.Examples[0].Command != "myapp serve --port 9090" {
		t.Errorf("first example command = %q", hc.Examples[0].Command)
	}
}

func TestFromCobraNoLongDesc(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "simple",
		Short: "A simple command",
	}
	hc := FromCobra(cmd)

	if hc.Description != "A simple command" {
		t.Errorf("Description should fall back to Short, got %q", hc.Description)
	}
}

func TestFromCobraNoExamples(t *testing.T) {
	cmd := &cobra.Command{Use: "bare"}
	hc := FromCobra(cmd)

	if len(hc.Examples) != 0 {
		t.Errorf("Examples should be empty, got %d", len(hc.Examples))
	}
}

func TestFromCobraHiddenFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("visible", "", "visible flag")
	f := cmd.Flags().String("hidden", "", "hidden flag")
	_ = f
	_ = cmd.Flags().MarkHidden("hidden")

	hc := FromCobra(cmd)

	for _, flag := range hc.Flags {
		if flag.Long == "--hidden" && !flag.Hidden {
			t.Error("hidden flag should have Hidden=true")
		}
	}
}

func TestParseExamplesTrailingDesc(t *testing.T) {
	raw := "Just a description with no command"
	examples := parseExamples(raw)

	if len(examples) != 1 {
		t.Fatalf("examples count = %d, want 1", len(examples))
	}
	if examples[0].Desc != "Just a description with no command" {
		t.Errorf("desc = %q", examples[0].Desc)
	}
	if examples[0].Command != "" {
		t.Errorf("command should be empty, got %q", examples[0].Command)
	}
}

func TestParseExamplesEmpty(t *testing.T) {
	examples := parseExamples("")
	if len(examples) != 0 {
		t.Errorf("empty input should produce no examples, got %d", len(examples))
	}
}

func TestFromCobraRequiredFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "User name")
	_ = cmd.MarkFlagRequired("name")

	hc := FromCobra(cmd)

	for _, f := range hc.Flags {
		if f.Long == "--name" {
			if !f.Required {
				t.Error("--name should be marked required")
			}
			return
		}
	}
	t.Error("--name flag not found")
}

func TestParseExamplesMultilineDesc(t *testing.T) {
	raw := "First line of description\nSecond line of description\n$ myapp run"
	examples := parseExamples(raw)

	if len(examples) != 1 {
		t.Fatalf("examples count = %d, want 1", len(examples))
	}
	if examples[0].Command != "myapp run" {
		t.Errorf("command = %q, want %q", examples[0].Command, "myapp run")
	}
	if examples[0].Desc != "First line of description Second line of description" {
		t.Errorf("desc = %q", examples[0].Desc)
	}
}

func TestIsCommandLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"$ myapp run", true},
		{"$\tmyapp run", true},
		{"myapp run", false},
		{"# comment", false},
		{"", false},
	}

	for _, tc := range tests {
		if got := isCommandLine(tc.line); got != tc.want {
			t.Errorf("isCommandLine(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
