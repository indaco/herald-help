package heraldhelp

import (
	"bytes"
	"errors"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/indaco/herald"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func newTestTypography() *herald.Typography {
	return herald.New()
}

func fullCommand() Command {
	return Command{
		Name:        "myapp",
		Synopsis:    "myapp [flags] <subcommand>",
		Description: "A sample application for testing herald-help rendering.",
		Aliases:     []string{"ma"},
		Deprecated:  "use newapp instead",
		Flags: []Flag{
			{Long: "--output", Short: "-o", Type: "string", Default: "stdout", Desc: "Output destination"},
			{Long: "--verbose", Short: "-v", Type: "bool", Default: "false", Desc: "Enable verbose output"},
			{Long: "--port", Type: "int", Default: "8080", Desc: "Port number", EnvVars: []string{"PORT"}},
			{Long: "--secret", Type: "string", Desc: "Secret key", Hidden: true},
			{Long: "--config", Type: "string", Desc: "Config file", Inherited: true},
			{Long: "--format", Type: "string", Desc: "Output format", Enum: []string{"json", "yaml", "toml"}},
			{Long: "--old-flag", Type: "string", Desc: "Old flag", Deprecated: "use --format"},
			{Long: "--count", Type: "int", Desc: "Count", Required: true},
		},
		FlagGroups: []FlagGroup{
			{
				Name: "Server Options",
				Flags: []Flag{
					{Long: "--host", Type: "string", Default: "localhost", Desc: "Host address"},
					{Long: "--tls", Type: "bool", Default: "false", Desc: "Enable TLS"},
					{Long: "--inherited-group", Type: "string", Desc: "Group inherited", Inherited: true},
				},
			},
		},
		Args: []Arg{
			{Name: "<file>", Desc: "Input file", Required: true},
			{Name: "<output>", Desc: "Output file", Default: "out.txt"},
		},
		Commands: []CommandRef{
			{Name: "serve", Aliases: []string{"s"}, Desc: "Start the server"},
			{Name: "build", Desc: "Build the project"},
		},
		CommandGroups: []CommandGroup{
			{
				Name: "Management",
				Commands: []CommandRef{
					{Name: "config", Desc: "Manage configuration"},
					{Name: "plugin", Aliases: []string{"p"}, Desc: "Manage plugins"},
				},
			},
		},
		Examples: []Example{
			{Desc: "Run the server on port 9090:", Command: "myapp serve --port 9090"},
			{Desc: "Build with verbose output:", Command: "myapp build -v"},
			{Command: "myapp --help"},
		},
		SeeAlso: []string{"newapp", "myapp-serve(1)"},
		Footer:  "myapp version 1.2.3",
	}
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

func TestRenderNilTypography(t *testing.T) {
	result := Render(nil, Command{Name: "test"})
	if result != "" {
		t.Errorf("Render(nil, ...) = %q, want empty", result)
	}
}

func TestRenderEmptyCommand(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{})
	if result != "" {
		t.Errorf("Render(empty command) = %q, want empty", result)
	}
}

func TestRenderFullCommand(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	checks := []string{
		"myapp",
		"Deprecated",
		"use newapp instead",
		"myapp [flags] <subcommand>",
		"A sample application",
		"Arguments",
		"<file>",
		"Input file",
		"Flags",
		"--output",
		"-o",
		"--verbose",
		"--port",
		"$PORT",
		"Server Options",
		"--host",
		"Commands",
		"serve",
		"build",
		"Management",
		"config",
		"Examples",
		"myapp serve --port 9090",
		"See Also",
		"newapp",
		"myapp version 1.2.3",
	}

	for _, want := range checks {
		if !strings.Contains(result, want) {
			t.Errorf("Render(fullCommand) missing %q in:\n%s", want, result)
		}
	}
}

func TestRenderHiddenFlagsExcluded(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if strings.Contains(result, "--secret") {
		t.Errorf("hidden flag --secret should not appear in output:\n%s", result)
	}
}

func TestRenderHiddenFlagsIncluded(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich), WithShowHidden(true)))

	if !strings.Contains(result, "--secret") {
		t.Errorf("hidden flag --secret should appear with WithShowHidden:\n%s", result)
	}
}

func TestRenderInheritedFlags(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Inherited Flags") {
		t.Errorf("should contain 'Inherited Flags' section:\n%s", result)
	}
	if !strings.Contains(result, "--config") {
		t.Errorf("inherited flag --config should appear in output:\n%s", result)
	}
	if !strings.Contains(result, "--inherited-group") {
		t.Errorf("inherited flag --inherited-group should appear in output:\n%s", result)
	}
}

func TestRenderEnvVarDisplayOff(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich), WithEnvVarDisplay(false)))

	if strings.Contains(result, "$PORT") {
		t.Errorf("env var $PORT should not appear with WithEnvVarDisplay(false):\n%s", result)
	}
}

func TestRenderEnumFlags(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "enum: json, yaml, toml") {
		t.Errorf("enum values should appear in output:\n%s", result)
	}
}

func TestRenderDeprecatedFlag(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "DEPRECATED: use --format") {
		t.Errorf("deprecated flag notice should appear in output:\n%s", result)
	}
}

func TestRenderRequiredFlag(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "(required)") {
		t.Errorf("required flag annotation should appear in output:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Section ordering
// ---------------------------------------------------------------------------

func TestRenderCustomSectionOrder(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich), WithSectionOrder(SectionName, SectionFooter)))

	if !strings.Contains(result, "myapp") {
		t.Errorf("should contain name:\n%s", result)
	}
	if !strings.Contains(result, "version 1.2.3") {
		t.Errorf("should contain footer:\n%s", result)
	}
	// Should not contain other sections.
	if strings.Contains(result, "Arguments") {
		t.Errorf("should not contain Arguments section:\n%s", result)
	}
	if strings.Contains(result, "Examples") {
		t.Errorf("should not contain Examples section:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Individual section rendering
// ---------------------------------------------------------------------------

func TestRenderNameOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "testcmd") {
		t.Errorf("should contain command name:\n%s", result)
	}
}

func TestRenderSynopsisOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Synopsis: "cmd [flags]"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "cmd [flags]") {
		t.Errorf("should contain synopsis:\n%s", result)
	}
}

func TestRenderDescriptionOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Description: "A great tool."}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "A great tool.") {
		t.Errorf("should contain description:\n%s", result)
	}
}

func TestRenderDeprecatedOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Deprecated: "use other"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Deprecated") {
		t.Errorf("should contain deprecation warning:\n%s", result)
	}
}

func TestRenderArgsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Args: []Arg{
			{Name: "<src>", Desc: "Source file", Required: true},
			{Name: "<dst>", Desc: "Destination", Default: "/tmp"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Arguments") {
		t.Errorf("should contain Arguments heading:\n%s", result)
	}
	if !strings.Contains(result, "<src>") {
		t.Errorf("should contain arg name:\n%s", result)
	}
	if !strings.Contains(result, "yes") {
		t.Errorf("should contain required marker:\n%s", result)
	}
}

func TestRenderFlagsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Flags: []Flag{
			{Long: "--debug", Short: "-d", Type: "bool", Default: "false", Desc: "Debug mode"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Flags") {
		t.Errorf("should contain Flags heading:\n%s", result)
	}
	if !strings.Contains(result, "-d, --debug") {
		t.Errorf("should contain GNU-style flag format:\n%s", result)
	}
}

func TestRenderCommandsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Commands: []CommandRef{
			{Name: "init", Desc: "Initialize project"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Commands") {
		t.Errorf("should contain Commands heading:\n%s", result)
	}
	if !strings.Contains(result, "init") {
		t.Errorf("should contain subcommand:\n%s", result)
	}
}

func TestRenderExamplesOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Examples: []Example{
			{Desc: "Basic usage:", Command: "cmd run"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Examples") {
		t.Errorf("should contain Examples heading:\n%s", result)
	}
	if !strings.Contains(result, "Basic usage:") {
		t.Errorf("should contain example description:\n%s", result)
	}
	if !strings.Contains(result, "cmd run") {
		t.Errorf("should contain example command:\n%s", result)
	}
}

func TestRenderExampleNoDesc(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Examples: []Example{
			{Command: "cmd --help"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "cmd --help") {
		t.Errorf("should contain example command:\n%s", result)
	}
}

func TestRenderSeeAlsoOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:    "cmd",
		SeeAlso: []string{"other-cmd", "docs"},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "See Also") {
		t.Errorf("should contain See Also heading:\n%s", result)
	}
	if !strings.Contains(result, "other-cmd") {
		t.Errorf("should contain see-also entry:\n%s", result)
	}
}

func TestRenderFooterOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Footer: "v1.0.0"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "v1.0.0") {
		t.Errorf("should contain footer:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// RenderTo
// ---------------------------------------------------------------------------

func TestRenderTo(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd", Synopsis: "testcmd [flags]"}
	var buf bytes.Buffer
	err := RenderTo(&buf, ty, cmd)
	if err != nil {
		t.Fatalf("RenderTo returned error: %v", err)
	}
	result := stripANSI(buf.String())
	if !strings.Contains(result, "testcmd") {
		t.Errorf("RenderTo output missing command name:\n%s", result)
	}
}

func TestRenderToNilTypography(t *testing.T) {
	var buf bytes.Buffer
	err := RenderTo(&buf, nil, Command{Name: "test"})
	if err != nil {
		t.Fatalf("RenderTo(nil) returned error: %v", err)
	}
	if buf.String() != "\n" {
		t.Errorf("RenderTo(nil) should produce only newline, got %q", buf.String())
	}
}

// ---------------------------------------------------------------------------
// Flag name formatting
// ---------------------------------------------------------------------------

func TestFormatFlagName(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{"short and long", Flag{Short: "-o", Long: "--output"}, "-o, --output"},
		{"long only", Flag{Long: "--output"}, "    --output"},
		{"short only", Flag{Short: "-o"}, "-o"},
		{"empty", Flag{}, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatFlagName(tc.flag)
			if got != tc.want {
				t.Errorf("formatFlagName(%+v) = %q, want %q", tc.flag, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Env var formatting
// ---------------------------------------------------------------------------

func TestFormatEnvVars(t *testing.T) {
	tests := []struct {
		name string
		vars []string
		want string
	}{
		{"single", []string{"PORT"}, "$PORT"},
		{"multiple", []string{"PORT", "APP_PORT"}, "$PORT, $APP_PORT"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatEnvVars(tc.vars)
			if got != tc.want {
				t.Errorf("formatEnvVars(%v) = %q, want %q", tc.vars, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// FormatVersion
// ---------------------------------------------------------------------------

func TestFormatVersion(t *testing.T) {
	got := FormatVersion("myapp", "1.2.3")
	if got != "myapp version 1.2.3" {
		t.Errorf("FormatVersion = %q, want %q", got, "myapp version 1.2.3")
	}
}

// ---------------------------------------------------------------------------
// Flag groups only (no flat flags)
// ---------------------------------------------------------------------------

func TestRenderFlagGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		FlagGroups: []FlagGroup{
			{
				Name: "Output",
				Flags: []Flag{
					{Long: "--json", Type: "bool", Default: "false", Desc: "JSON output"},
				},
			},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Output") {
		t.Errorf("should contain flag group name:\n%s", result)
	}
	if !strings.Contains(result, "--json") {
		t.Errorf("should contain grouped flag:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Command groups only (no flat commands)
// ---------------------------------------------------------------------------

func TestRenderCommandGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		CommandGroups: []CommandGroup{
			{
				Name: "Advanced",
				Commands: []CommandRef{
					{Name: "migrate", Desc: "Run migrations"},
				},
			},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if !strings.Contains(result, "Commands") {
		t.Errorf("should contain Commands heading:\n%s", result)
	}
	if !strings.Contains(result, "Advanced") {
		t.Errorf("should contain group name:\n%s", result)
	}
	if !strings.Contains(result, "migrate") {
		t.Errorf("should contain grouped command:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Empty flag groups and command groups
// ---------------------------------------------------------------------------

func TestRenderEmptyFlagGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		FlagGroups: []FlagGroup{
			{Name: "Empty", Flags: nil},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if strings.Contains(result, "Flags") {
		t.Errorf("empty flag group should not produce Flags section:\n%s", result)
	}
}

func TestRenderEmptyCommandGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		CommandGroups: []CommandGroup{
			{Name: "Empty", Commands: nil},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich)))

	if strings.Contains(result, "Commands") {
		t.Errorf("empty command group should not produce Commands section:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// renderSection unknown section
// ---------------------------------------------------------------------------

func TestRenderSectionUnknown(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd"}
	result := renderSection(ty, cmd, Section(999), &RenderConfig{Style: StyleRich})
	if result != "" {
		t.Errorf("unknown section should return empty, got %q", result)
	}
}

func TestRenderNameEmpty(t *testing.T) {
	ty := newTestTypography()
	result := renderName(ty, Command{})
	if result != "" {
		t.Errorf("renderName(empty) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// buildConfig defaults
// ---------------------------------------------------------------------------

func TestBuildConfigDefaults(t *testing.T) {
	cfg := buildConfig(nil)
	if !cfg.EnvVarDisplay {
		t.Error("EnvVarDisplay should default to true")
	}
	if cfg.ShowHidden {
		t.Error("ShowHidden should default to false")
	}
	if cfg.SectionOrder != nil {
		t.Error("SectionOrder should default to nil")
	}
}

// ---------------------------------------------------------------------------
// filterFlags
// ---------------------------------------------------------------------------

func TestFilterFlags(t *testing.T) {
	flags := []Flag{
		{Long: "--visible", Hidden: false, Inherited: false},
		{Long: "--hidden", Hidden: true, Inherited: false},
		{Long: "--inherited", Hidden: false, Inherited: true},
		{Long: "--hidden-inherited", Hidden: true, Inherited: true},
	}

	tests := []struct {
		name       string
		showHidden bool
		inherited  bool
		wantCount  int
	}{
		{"non-inherited visible", false, false, 1},
		{"non-inherited all", true, false, 2},
		{"inherited visible", false, true, 1},
		{"inherited all", true, true, 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterFlags(flags, tc.showHidden, tc.inherited)
			if len(got) != tc.wantCount {
				t.Errorf("filterFlags(%v, %v) = %d flags, want %d", tc.showHidden, tc.inherited, len(got), tc.wantCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// RenderTo error path
// ---------------------------------------------------------------------------

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestRenderToWriteError(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd"}
	err := RenderTo(failWriter{}, ty, cmd)
	if err == nil {
		t.Error("RenderTo should return error on writer failure")
	}
}

// ---------------------------------------------------------------------------
// DefaultSectionOrder immutability
// ---------------------------------------------------------------------------

func TestDefaultSectionOrderReturnsACopy(t *testing.T) {
	order1 := DefaultSectionOrder()
	original := order1[0]
	// Mutate the returned slice.
	order1[0] = SectionFooter
	// A second call must return the original, unmodified order.
	order2 := DefaultSectionOrder()
	if order2[0] != original {
		t.Errorf("DefaultSectionOrder is not returning a copy: got %d, want %d", order2[0], original)
	}
}

// ---------------------------------------------------------------------------
// WithoutSections
// ---------------------------------------------------------------------------

func TestWithoutSections(t *testing.T) {
	tests := []struct {
		name    string
		exclude []Section
		check   func([]Section) bool
		desc    string
	}{
		{
			name:    "exclude one section",
			exclude: []Section{SectionInheritedFlags},
			check: func(order []Section) bool {
				return !slices.Contains(order, SectionInheritedFlags)
			},
			desc: "SectionInheritedFlags should be excluded",
		},
		{
			name:    "exclude multiple sections",
			exclude: []Section{SectionFlags, SectionCommands},
			check: func(order []Section) bool {
				for _, s := range order {
					if s == SectionFlags || s == SectionCommands {
						return false
					}
				}
				return true
			},
			desc: "SectionFlags and SectionCommands should be excluded",
		},
		{
			name:    "no args leaves order unchanged",
			exclude: nil,
			check: func(order []Section) bool {
				def := DefaultSectionOrder()
				if len(order) != len(def) {
					return false
				}
				for i := range order {
					if order[i] != def[i] {
						return false
					}
				}
				return true
			},
			desc: "order should equal DefaultSectionOrder()",
		},
		{
			name:    "unknown section is no-op",
			exclude: []Section{Section(999)},
			check: func(order []Section) bool {
				def := DefaultSectionOrder()
				return len(order) == len(def)
			},
			desc: "order length should be unchanged",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &RenderConfig{}
			WithoutSections(tc.exclude...)(cfg)
			if !tc.check(cfg.SectionOrder) {
				t.Errorf("%s: got %v", tc.desc, cfg.SectionOrder)
			}
		})
	}
}

func TestWithoutSectionsRendering(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleRich), WithoutSections(SectionExamples, SectionSeeAlso)))

	if strings.Contains(result, "Examples") {
		t.Errorf("Examples section should be excluded:\n%s", result)
	}
	if strings.Contains(result, "See Also") {
		t.Errorf("See Also section should be excluded:\n%s", result)
	}
	// Other sections should still be present.
	if !strings.Contains(result, "myapp") {
		t.Errorf("Name section should be present:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Benchmark
// ---------------------------------------------------------------------------

func BenchmarkRender(b *testing.B) {
	ty := herald.New()
	cmd := fullCommand()

	for b.Loop() {
		_ = Render(ty, cmd)
	}
}
