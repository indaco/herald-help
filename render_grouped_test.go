package heraldhelp

import (
	"bytes"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Grouped: full command
// ---------------------------------------------------------------------------

func TestGroupedRenderFullCommand(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))

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
			t.Errorf("Grouped Render missing %q in:\n%s", want, result)
		}
	}
}

// ---------------------------------------------------------------------------
// Grouped: empty command
// ---------------------------------------------------------------------------

func TestGroupedRenderEmptyCommand(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{}, WithStyle(StyleGrouped))
	if result != "" {
		t.Errorf("Grouped Render(empty) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// Grouped: individual sections
// ---------------------------------------------------------------------------

func TestGroupedRenderNameOnly(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, Command{Name: "cmd"}, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "cmd") {
		t.Errorf("should contain name:\n%s", result)
	}
}

func TestGroupedRenderNameEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedName(ty, Command{})
	if result != "" {
		t.Errorf("groupedName(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderSynopsisOnly(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, Command{Name: "cmd", Synopsis: "cmd [flags]"}, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Usage") {
		t.Errorf("should contain Usage legend:\n%s", result)
	}
	if !strings.Contains(result, "cmd [flags]") {
		t.Errorf("should contain synopsis:\n%s", result)
	}
}

func TestGroupedRenderSynopsisEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedSynopsis(ty, Command{})
	if result != "" {
		t.Errorf("groupedSynopsis(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderDescriptionOnly(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, Command{Name: "cmd", Description: "A tool."}, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "A tool.") {
		t.Errorf("should contain description:\n%s", result)
	}
}

func TestGroupedRenderDescriptionEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedDescription(ty, Command{})
	if result != "" {
		t.Errorf("groupedDescription(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderDeprecatedOnly(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, Command{Name: "cmd", Deprecated: "use other"}, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Deprecated") {
		t.Errorf("should contain deprecation:\n%s", result)
	}
}

func TestGroupedRenderDeprecatedEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedDeprecated(ty, Command{})
	if result != "" {
		t.Errorf("groupedDeprecated(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderArgsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Args: []Arg{
			{Name: "<src>", Desc: "Source", Required: true},
			{Name: "<dst>", Desc: "Dest", Default: "/tmp"},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Arguments") {
		t.Errorf("should contain Arguments legend:\n%s", result)
	}
	if !strings.Contains(result, "<src>") {
		t.Errorf("should contain arg:\n%s", result)
	}
	if !strings.Contains(result, "(required)") {
		t.Errorf("should contain required:\n%s", result)
	}
}

func TestGroupedRenderArgsEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedArgs(ty, Command{})
	if result != "" {
		t.Errorf("groupedArgs(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderFlagsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:  "cmd",
		Flags: []Flag{{Long: "--debug", Short: "-d", Type: "bool", Desc: "Debug"}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Flags") {
		t.Errorf("should contain Flags legend:\n%s", result)
	}
	if !strings.Contains(result, "-d, --debug") {
		t.Errorf("should contain flag:\n%s", result)
	}
}

func TestGroupedRenderFlagGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		FlagGroups: []FlagGroup{
			{Name: "Output", Flags: []Flag{{Long: "--json", Type: "bool", Desc: "JSON"}}},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Output") {
		t.Errorf("should contain group legend:\n%s", result)
	}
}

func TestGroupedRenderEmptyFlagGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", FlagGroups: []FlagGroup{{Name: "Empty"}}}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if strings.Contains(result, "Flags") {
		t.Errorf("empty flag group should not produce Flags:\n%s", result)
	}
}

func TestGroupedRenderInheritedFlags(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, fullCommand(), WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Inherited Flags") {
		t.Errorf("should contain Inherited Flags:\n%s", result)
	}
}

func TestGroupedRenderHiddenFlagsExcluded(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, fullCommand(), WithStyle(StyleGrouped)))
	if strings.Contains(result, "--secret") {
		t.Errorf("hidden flag should not appear:\n%s", result)
	}
}

func TestGroupedRenderHiddenFlagsIncluded(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, fullCommand(), WithStyle(StyleGrouped), WithShowHidden(true)))
	if !strings.Contains(result, "--secret") {
		t.Errorf("hidden flag should appear:\n%s", result)
	}
}

func TestGroupedRenderEnvVarDisplayOff(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, fullCommand(), WithStyle(StyleGrouped), WithEnvVarDisplay(false)))
	if strings.Contains(result, "$PORT") {
		t.Errorf("env var should not appear:\n%s", result)
	}
}

func TestGroupedRenderCommandsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:     "cmd",
		Commands: []CommandRef{{Name: "init", Desc: "Initialize"}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Commands") {
		t.Errorf("should contain Commands legend:\n%s", result)
	}
	if !strings.Contains(result, "init") {
		t.Errorf("should contain command:\n%s", result)
	}
}

func TestGroupedRenderCommandGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		CommandGroups: []CommandGroup{
			{Name: "Advanced", Commands: []CommandRef{{Name: "migrate", Desc: "Run migrations"}}},
		},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Advanced") {
		t.Errorf("should contain group legend:\n%s", result)
	}
}

func TestGroupedRenderEmptyCommandGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", CommandGroups: []CommandGroup{{Name: "Empty"}}}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if strings.Contains(result, "Commands") {
		t.Errorf("empty command group should not produce Commands:\n%s", result)
	}
}

func TestGroupedRenderExamplesOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:     "cmd",
		Examples: []Example{{Desc: "Basic:", Command: "cmd run"}},
	}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "Examples") {
		t.Errorf("should contain Examples legend:\n%s", result)
	}
	if !strings.Contains(result, "cmd run") {
		t.Errorf("should contain command:\n%s", result)
	}
}

func TestGroupedRenderExamplesEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedExamples(ty, Command{})
	if result != "" {
		t.Errorf("groupedExamples(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderSeeAlsoOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", SeeAlso: []string{"other-cmd"}}
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "See Also") {
		t.Errorf("should contain See Also legend:\n%s", result)
	}
	if !strings.Contains(result, "other-cmd") {
		t.Errorf("should contain entry:\n%s", result)
	}
}

func TestGroupedRenderSeeAlsoEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedSeeAlso(ty, Command{})
	if result != "" {
		t.Errorf("groupedSeeAlso(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderFooterOnly(t *testing.T) {
	ty := newTestTypography()
	result := stripANSI(Render(ty, Command{Name: "cmd", Footer: "v1.0.0"}, WithStyle(StyleGrouped)))
	if !strings.Contains(result, "v1.0.0") {
		t.Errorf("should contain footer:\n%s", result)
	}
}

func TestGroupedRenderFooterEmpty(t *testing.T) {
	ty := newTestTypography()
	result := groupedFooter(ty, Command{})
	if result != "" {
		t.Errorf("groupedFooter(empty) = %q, want empty", result)
	}
}

func TestGroupedRenderTo(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd", Synopsis: "testcmd [flags]"}
	var buf bytes.Buffer
	err := RenderTo(&buf, ty, cmd, WithStyle(StyleGrouped))
	if err != nil {
		t.Fatalf("RenderTo error: %v", err)
	}
	result := stripANSI(buf.String())
	if !strings.Contains(result, "testcmd") {
		t.Errorf("should contain name:\n%s", result)
	}
}

func TestGroupedRenderSectionUnknown(t *testing.T) {
	ty := newTestTypography()
	result := renderSectionGrouped(ty, Command{}, Section(999), &RenderConfig{})
	if result != "" {
		t.Errorf("unknown section should return empty, got %q", result)
	}
}

func TestGroupedRenderCustomSectionOrder(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := stripANSI(Render(ty, cmd, WithStyle(StyleGrouped), WithSectionOrder(SectionName, SectionFooter)))
	if !strings.Contains(result, "myapp") {
		t.Errorf("should contain name:\n%s", result)
	}
	if strings.Contains(result, "Flags") {
		t.Errorf("should not contain Flags:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Benchmark
// ---------------------------------------------------------------------------

func BenchmarkGroupedRender(b *testing.B) {
	ty := newTestTypography()
	cmd := fullCommand()
	for b.Loop() {
		_ = Render(ty, cmd, WithStyle(StyleGrouped))
	}
}
