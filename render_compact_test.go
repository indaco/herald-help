package heraldhelp

import (
	"bytes"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Compact: full command
// ---------------------------------------------------------------------------

func TestCompactRenderFullCommand(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	checks := []string{
		"myapp",
		"Deprecated",
		"use newapp instead",
		"myapp [flags] <subcommand>",
		"A sample application",
		"ARGUMENTS",
		"<file>",
		"Input file",
		"FLAGS",
		"--output",
		"-o",
		"--verbose",
		"--port",
		"$PORT",
		"SERVER OPTIONS",
		"--host",
		"COMMANDS",
		"serve",
		"build",
		"MANAGEMENT",
		"config",
		"EXAMPLES",
		"myapp serve --port 9090",
		"SEE ALSO",
		"newapp",
		"myapp version 1.2.3",
	}

	for _, want := range checks {
		if !strings.Contains(result, want) {
			t.Errorf("Compact Render missing %q in:\n%s", want, result)
		}
	}
}

// ---------------------------------------------------------------------------
// Compact: nil typography
// ---------------------------------------------------------------------------

func TestCompactRenderNilTypography(t *testing.T) {
	result := Render(nil, Command{Name: "test"}, WithStyle(StyleCompact))
	if result != "" {
		t.Errorf("Render(nil, ..., Compact) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// Compact: empty command
// ---------------------------------------------------------------------------

func TestCompactRenderEmptyCommand(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{}, WithStyle(StyleCompact))
	if result != "" {
		t.Errorf("Compact Render(empty) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// Compact: individual sections
// ---------------------------------------------------------------------------

func TestCompactRenderNameOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "testcmd") {
		t.Errorf("should contain command name:\n%s", result)
	}
}

func TestCompactRenderNameEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactName(ty, Command{})
	if result != "" {
		t.Errorf("compactName(empty) = %q, want empty", result)
	}
}

func TestCompactRenderSynopsisOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Synopsis: "cmd [flags]"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "USAGE") {
		t.Errorf("should contain USAGE heading:\n%s", result)
	}
	if !strings.Contains(result, "cmd [flags]") {
		t.Errorf("should contain synopsis:\n%s", result)
	}
}

func TestCompactRenderSynopsisEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactSynopsis(ty, Command{})
	if result != "" {
		t.Errorf("compactSynopsis(empty) = %q, want empty", result)
	}
}

func TestCompactRenderDescriptionOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Description: "A great tool."}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "A great tool.") {
		t.Errorf("should contain description:\n%s", result)
	}
}

func TestCompactRenderDescriptionEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactDescription(ty, Command{})
	if result != "" {
		t.Errorf("compactDescription(empty) = %q, want empty", result)
	}
}

func TestCompactRenderDeprecatedOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Deprecated: "use other"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "Deprecated") {
		t.Errorf("should contain deprecation warning:\n%s", result)
	}
}

func TestCompactRenderDeprecatedEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactDeprecated(ty, Command{})
	if result != "" {
		t.Errorf("compactDeprecated(empty) = %q, want empty", result)
	}
}

func TestCompactRenderArgsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Args: []Arg{
			{Name: "<src>", Desc: "Source file", Required: true},
			{Name: "<dst>", Desc: "Destination", Default: "/tmp"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "ARGUMENTS") {
		t.Errorf("should contain ARGUMENTS heading:\n%s", result)
	}
	if !strings.Contains(result, "<src>") {
		t.Errorf("should contain arg name:\n%s", result)
	}
	if !strings.Contains(result, "(required)") {
		t.Errorf("should contain required marker:\n%s", result)
	}
	if !strings.Contains(result, "(default: /tmp)") {
		t.Errorf("should contain default value:\n%s", result)
	}
}

func TestCompactRenderArgsEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactArgs(ty, Command{})
	if result != "" {
		t.Errorf("compactArgs(empty) = %q, want empty", result)
	}
}

func TestCompactRenderFlagsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Flags: []Flag{
			{Long: "--debug", Short: "-d", Type: "bool", Default: "false", Desc: "Debug mode"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "FLAGS") {
		t.Errorf("should contain FLAGS heading:\n%s", result)
	}
	if !strings.Contains(result, "-d, --debug") {
		t.Errorf("should contain GNU-style flag format:\n%s", result)
	}
}

func TestCompactRenderFlagTypeInline(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Flags: []Flag{
			{Long: "--port", Type: "int", Default: "8080", Desc: "Port number"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "--port int") {
		t.Errorf("should contain type inline with flag name:\n%s", result)
	}
	if !strings.Contains(result, "(default: 8080)") {
		t.Errorf("should contain default inline:\n%s", result)
	}
}

func TestCompactRenderHiddenFlagsExcluded(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if strings.Contains(result, "--secret") {
		t.Errorf("hidden flag --secret should not appear:\n%s", result)
	}
}

func TestCompactRenderHiddenFlagsIncluded(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact), WithShowHidden(true)))

	if !strings.Contains(result, "--secret") {
		t.Errorf("hidden flag --secret should appear with WithShowHidden:\n%s", result)
	}
}

func TestCompactRenderEnvVarDisplayOff(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact), WithEnvVarDisplay(false)))

	if strings.Contains(result, "$PORT") {
		t.Errorf("env var $PORT should not appear:\n%s", result)
	}
}

func TestCompactRenderEnumFlags(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "enum: json, yaml, toml") {
		t.Errorf("enum values should appear:\n%s", result)
	}
}

func TestCompactRenderDeprecatedFlag(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "DEPRECATED: use --format") {
		t.Errorf("deprecated flag notice should appear:\n%s", result)
	}
}

func TestCompactRenderRequiredFlag(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "(required)") {
		t.Errorf("required flag annotation should appear:\n%s", result)
	}
}

func TestCompactRenderInheritedFlags(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "INHERITED FLAGS") {
		t.Errorf("should contain INHERITED FLAGS section:\n%s", result)
	}
	if !strings.Contains(result, "--config") {
		t.Errorf("inherited flag --config should appear:\n%s", result)
	}
	if !strings.Contains(result, "--inherited-group") {
		t.Errorf("inherited flag --inherited-group should appear:\n%s", result)
	}
}

func TestCompactRenderFlagGroupsOnly(t *testing.T) {
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
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "OUTPUT") {
		t.Errorf("should contain flag group name uppercase:\n%s", result)
	}
	if !strings.Contains(result, "--json") {
		t.Errorf("should contain grouped flag:\n%s", result)
	}
}

func TestCompactRenderEmptyFlagGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:       "cmd",
		FlagGroups: []FlagGroup{{Name: "Empty", Flags: nil}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if strings.Contains(result, "FLAGS") {
		t.Errorf("empty flag group should not produce FLAGS section:\n%s", result)
	}
}

func TestCompactRenderCommandsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Commands: []CommandRef{
			{Name: "init", Desc: "Initialize project"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "COMMANDS") {
		t.Errorf("should contain COMMANDS heading:\n%s", result)
	}
	if !strings.Contains(result, "init") {
		t.Errorf("should contain subcommand:\n%s", result)
	}
}

func TestCompactRenderCommandAliases(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Commands: []CommandRef{
			{Name: "serve", Aliases: []string{"s"}, Desc: "Start server"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "serve, s") {
		t.Errorf("should contain command with aliases:\n%s", result)
	}
}

func TestCompactRenderCommandGroupsOnly(t *testing.T) {
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
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "COMMANDS") {
		t.Errorf("should contain COMMANDS heading:\n%s", result)
	}
	if !strings.Contains(result, "ADVANCED") {
		t.Errorf("should contain group name uppercase:\n%s", result)
	}
	if !strings.Contains(result, "migrate") {
		t.Errorf("should contain grouped command:\n%s", result)
	}
}

func TestCompactRenderEmptyCommandGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:          "cmd",
		CommandGroups: []CommandGroup{{Name: "Empty", Commands: nil}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if strings.Contains(result, "COMMANDS") {
		t.Errorf("empty command group should not produce COMMANDS section:\n%s", result)
	}
}

func TestCompactRenderExamplesOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Examples: []Example{
			{Desc: "Basic usage:", Command: "cmd run"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "EXAMPLES") {
		t.Errorf("should contain EXAMPLES heading:\n%s", result)
	}
	if !strings.Contains(result, "Basic usage:") {
		t.Errorf("should contain example description:\n%s", result)
	}
	if !strings.Contains(result, "cmd run") {
		t.Errorf("should contain example command:\n%s", result)
	}
}

func TestCompactRenderExampleNoDesc(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:     "cmd",
		Examples: []Example{{Command: "cmd --help"}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "cmd --help") {
		t.Errorf("should contain example command:\n%s", result)
	}
}

func TestCompactRenderExamplesEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactExamples(ty, Command{})
	if result != "" {
		t.Errorf("compactExamples(empty) = %q, want empty", result)
	}
}

func TestCompactRenderSeeAlsoOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:    "cmd",
		SeeAlso: []string{"other-cmd", "docs"},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "SEE ALSO") {
		t.Errorf("should contain SEE ALSO heading:\n%s", result)
	}
	if !strings.Contains(result, "other-cmd") {
		t.Errorf("should contain see-also entry:\n%s", result)
	}
}

func TestCompactRenderSeeAlsoEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactSeeAlso(ty, Command{})
	if result != "" {
		t.Errorf("compactSeeAlso(empty) = %q, want empty", result)
	}
}

func TestCompactRenderFooterOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Footer: "v1.0.0"}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact)))

	if !strings.Contains(result, "v1.0.0") {
		t.Errorf("should contain footer:\n%s", result)
	}
}

func TestCompactRenderFooterEmpty(t *testing.T) {
	ty := newTestTypography()
	result := compactFooter(ty, Command{})
	if result != "" {
		t.Errorf("compactFooter(empty) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// Compact: RenderTo
// ---------------------------------------------------------------------------

func TestCompactRenderTo(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd", Synopsis: "testcmd [flags]"}
	var buf bytes.Buffer
	err := RenderTo(&buf, ty, cmd, WithStyle(StyleCompact))
	if err != nil {
		t.Fatalf("RenderTo returned error: %v", err)
	}
	result := stripANSI(buf.String())
	if !strings.Contains(result, "testcmd") {
		t.Errorf("RenderTo compact output missing command name:\n%s", result)
	}
	if !strings.Contains(result, "USAGE") {
		t.Errorf("RenderTo compact output missing USAGE heading:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Compact: section ordering
// ---------------------------------------------------------------------------

func TestCompactRenderCustomSectionOrder(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleCompact), WithSectionOrder(SectionName, SectionFooter)))

	if !strings.Contains(result, "myapp") {
		t.Errorf("should contain name:\n%s", result)
	}
	if !strings.Contains(result, "version 1.2.3") {
		t.Errorf("should contain footer:\n%s", result)
	}
	if strings.Contains(result, "ARGUMENTS") {
		t.Errorf("should not contain Arguments section:\n%s", result)
	}
	if strings.Contains(result, "EXAMPLES") {
		t.Errorf("should not contain Examples section:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Compact: unknown section
// ---------------------------------------------------------------------------

func TestCompactRenderSectionUnknown(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd"}
	result := renderSectionCompact(ty, cmd, Section(999), &RenderConfig{})
	if result != "" {
		t.Errorf("unknown section should return empty, got %q", result)
	}
}

// ---------------------------------------------------------------------------
// WithStyle option
// ---------------------------------------------------------------------------

func TestWithStyle(t *testing.T) {
	cfg := &RenderConfig{}
	WithStyle(StyleCompact)(cfg)
	if cfg.Style != StyleCompact {
		t.Errorf("WithStyle(StyleCompact) = %d, want %d", cfg.Style, StyleCompact)
	}
}

func TestWithStyleCompactIsDefault(t *testing.T) {
	cfg := buildConfig(nil)
	if cfg.Style != StyleCompact {
		t.Errorf("default style = %d, want StyleCompact (%d)", cfg.Style, StyleCompact)
	}
}

// ---------------------------------------------------------------------------
// Compact: benchmark
// ---------------------------------------------------------------------------

func BenchmarkCompactRender(b *testing.B) {
	ty := newTestTypography()
	cmd := fullCommand()

	for b.Loop() {
		_ = Render(ty, cmd, WithStyle(StyleCompact))
	}
}
