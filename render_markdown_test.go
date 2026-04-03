package heraldhelp

import (
	"bytes"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Markdown: full command
// ---------------------------------------------------------------------------

func TestMarkdownRenderFullCommand(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := Render(ty, cmd, WithStyle(StyleMarkdown))

	checks := []string{
		"# myapp",
		"Deprecated",
		"use newapp instead",
		"myapp [flags] <subcommand>",
		"A sample application",
		"## Arguments",
		"`<file>`",
		"Input file",
		"## Flags",
		"`-o, --output`",
		"--verbose",
		"--port",
		"`$PORT`",
		"### Server Options",
		"--host",
		"## Commands",
		"`serve`",
		"`build`",
		"### Management",
		"`config`",
		"## Examples",
		"myapp serve --port 9090",
		"## See Also",
		"newapp",
		"myapp version 1.2.3",
	}

	for _, want := range checks {
		if !strings.Contains(result, want) {
			t.Errorf("Markdown Render missing %q in:\n%s", want, result)
		}
	}
}

// ---------------------------------------------------------------------------
// Markdown: no ANSI
// ---------------------------------------------------------------------------

func TestMarkdownRenderNoANSI(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := Render(ty, cmd, WithStyle(StyleMarkdown))

	if strings.Contains(result, "\x1b[") {
		t.Errorf("Markdown output should not contain ANSI escapes:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Markdown: empty command
// ---------------------------------------------------------------------------

func TestMarkdownRenderEmptyCommand(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{}, WithStyle(StyleMarkdown))
	if result != "" {
		t.Errorf("Markdown Render(empty) = %q, want empty", result)
	}
}

// ---------------------------------------------------------------------------
// Markdown: individual sections
// ---------------------------------------------------------------------------

func TestMarkdownRenderNameOnly(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{Name: "cmd"}, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "# cmd") {
		t.Errorf("should contain # cmd:\n%s", result)
	}
}

func TestMarkdownRenderNameEmpty(t *testing.T) {
	result := mdName(Command{})
	if result != "" {
		t.Errorf("mdName(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderSynopsisOnly(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{Name: "cmd", Synopsis: "cmd [flags]"}, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "```\ncmd [flags]\n```") {
		t.Errorf("should contain fenced synopsis:\n%s", result)
	}
}

func TestMarkdownRenderSynopsisEmpty(t *testing.T) {
	result := mdSynopsis(Command{})
	if result != "" {
		t.Errorf("mdSynopsis(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderDescriptionOnly(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{Name: "cmd", Description: "A tool."}, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "A tool.") {
		t.Errorf("should contain description:\n%s", result)
	}
}

func TestMarkdownRenderDescriptionEmpty(t *testing.T) {
	result := mdDescription(Command{})
	if result != "" {
		t.Errorf("mdDescription(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderDeprecatedOnly(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, Command{Name: "cmd", Deprecated: "use other"}, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "> **Warning:**") {
		t.Errorf("should contain blockquote warning:\n%s", result)
	}
}

func TestMarkdownRenderDeprecatedEmpty(t *testing.T) {
	result := mdDeprecated(Command{})
	if result != "" {
		t.Errorf("mdDeprecated(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderArgsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		Args: []Arg{
			{Name: "<src>", Desc: "Source", Required: true},
			{Name: "<dst>", Desc: "Dest", Default: "/tmp"},
		},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "## Arguments") {
		t.Errorf("should contain Arguments heading:\n%s", result)
	}
	if !strings.Contains(result, "`<src>`") {
		t.Errorf("should contain arg in backticks:\n%s", result)
	}
	if !strings.Contains(result, "yes") {
		t.Errorf("should contain required:\n%s", result)
	}
}

func TestMarkdownRenderArgsEmpty(t *testing.T) {
	result := mdArgs(Command{})
	if result != "" {
		t.Errorf("mdArgs(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderFlagsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:  "cmd",
		Flags: []Flag{{Long: "--debug", Short: "-d", Type: "bool", Desc: "Debug"}},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "## Flags") {
		t.Errorf("should contain Flags heading:\n%s", result)
	}
	if !strings.Contains(result, "`-d, --debug`") {
		t.Errorf("should contain flag in backticks:\n%s", result)
	}
}

func TestMarkdownRenderFlagGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		FlagGroups: []FlagGroup{
			{Name: "Output", Flags: []Flag{{Long: "--json", Type: "bool", Desc: "JSON"}}},
		},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "### Output") {
		t.Errorf("should contain group heading:\n%s", result)
	}
}

func TestMarkdownRenderEmptyFlagGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", FlagGroups: []FlagGroup{{Name: "Empty"}}}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if strings.Contains(result, "## Flags") {
		t.Errorf("empty flag group should not produce Flags section:\n%s", result)
	}
}

func TestMarkdownRenderInheritedFlags(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "### Inherited Flags") {
		t.Errorf("should contain Inherited Flags:\n%s", result)
	}
}

func TestMarkdownRenderHiddenFlagsExcluded(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown))
	if strings.Contains(result, "--secret") {
		t.Errorf("hidden flag should not appear:\n%s", result)
	}
}

func TestMarkdownRenderHiddenFlagsIncluded(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown), WithShowHidden(true))
	if !strings.Contains(result, "--secret") {
		t.Errorf("hidden flag should appear with ShowHidden:\n%s", result)
	}
}

func TestMarkdownRenderEnvVarDisplayOff(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown), WithEnvVarDisplay(false))
	if strings.Contains(result, "$PORT") {
		t.Errorf("env var should not appear:\n%s", result)
	}
}

func TestMarkdownRenderEnumFlags(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown))
	if !strings.Contains(result, "enum:") {
		t.Errorf("enum values should appear:\n%s", result)
	}
}

func TestMarkdownRenderRequiredFlag(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown))
	if !strings.Contains(result, "**(required)**") {
		t.Errorf("required annotation should appear:\n%s", result)
	}
}

func TestMarkdownRenderDeprecatedFlag(t *testing.T) {
	ty := newTestTypography()
	result := Render(ty, fullCommand(), WithStyle(StyleMarkdown))
	if !strings.Contains(result, "DEPRECATED: use --format") {
		t.Errorf("deprecated flag notice should appear:\n%s", result)
	}
}

func TestMarkdownRenderCommandsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:     "cmd",
		Commands: []CommandRef{{Name: "init", Desc: "Initialize"}},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "## Commands") {
		t.Errorf("should contain Commands heading:\n%s", result)
	}
	if !strings.Contains(result, "`init`") {
		t.Errorf("should contain command in backticks:\n%s", result)
	}
}

func TestMarkdownRenderCommandGroupsOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name: "cmd",
		CommandGroups: []CommandGroup{
			{Name: "Advanced", Commands: []CommandRef{{Name: "migrate", Desc: "Run migrations"}}},
		},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "### Advanced") {
		t.Errorf("should contain group heading:\n%s", result)
	}
}

func TestMarkdownRenderEmptyCommandGroup(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", CommandGroups: []CommandGroup{{Name: "Empty"}}}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if strings.Contains(result, "## Commands") {
		t.Errorf("empty command group should not produce Commands section:\n%s", result)
	}
}

func TestMarkdownRenderExamplesOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{
		Name:     "cmd",
		Examples: []Example{{Desc: "Basic:", Command: "cmd run"}},
	}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "## Examples") {
		t.Errorf("should contain Examples heading:\n%s", result)
	}
	if !strings.Contains(result, "```\ncmd run\n```") {
		t.Errorf("should contain fenced command:\n%s", result)
	}
}

func TestMarkdownRenderExampleNoDesc(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Examples: []Example{{Command: "cmd --help"}}}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "cmd --help") {
		t.Errorf("should contain example command:\n%s", result)
	}
}

func TestMarkdownRenderExamplesEmpty(t *testing.T) {
	result := mdExamples(Command{})
	if result != "" {
		t.Errorf("mdExamples(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderSeeAlsoOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", SeeAlso: []string{"other-cmd"}}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "## See Also") {
		t.Errorf("should contain See Also:\n%s", result)
	}
	if !strings.Contains(result, "- other-cmd") {
		t.Errorf("should contain list item:\n%s", result)
	}
}

func TestMarkdownRenderSeeAlsoEmpty(t *testing.T) {
	result := mdSeeAlso(Command{})
	if result != "" {
		t.Errorf("mdSeeAlso(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderFooterOnly(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "cmd", Footer: "v1.0.0"}
	result := Render(ty, cmd, WithStyle(StyleMarkdown))
	if !strings.Contains(result, "*v1.0.0*") {
		t.Errorf("should contain italic footer:\n%s", result)
	}
	if !strings.Contains(result, "---") {
		t.Errorf("should contain horizontal rule:\n%s", result)
	}
}

func TestMarkdownRenderFooterEmpty(t *testing.T) {
	result := mdFooter(Command{})
	if result != "" {
		t.Errorf("mdFooter(empty) = %q, want empty", result)
	}
}

func TestMarkdownRenderTo(t *testing.T) {
	ty := newTestTypography()
	cmd := Command{Name: "testcmd", Synopsis: "testcmd [flags]"}
	var buf bytes.Buffer
	err := RenderTo(&buf, ty, cmd, WithStyle(StyleMarkdown))
	if err != nil {
		t.Fatalf("RenderTo error: %v", err)
	}
	result := buf.String()
	if !strings.Contains(result, "# testcmd") {
		t.Errorf("should contain Markdown heading:\n%s", result)
	}
}

func TestMdTableEmpty(t *testing.T) {
	result := mdTable(nil)
	if result != "" {
		t.Errorf("mdTable(nil) = %q, want empty", result)
	}
}

func TestMarkdownRenderSectionUnknown(t *testing.T) {
	result := renderSectionMarkdown(Command{}, Section(999), &RenderConfig{})
	if result != "" {
		t.Errorf("unknown section should return empty, got %q", result)
	}
}

func TestMarkdownRenderCustomSectionOrder(t *testing.T) {
	ty := newTestTypography()
	cmd := fullCommand()
	result := Render(ty, cmd, WithStyle(StyleMarkdown), WithSectionOrder(SectionName, SectionFooter))
	if !strings.Contains(result, "# myapp") {
		t.Errorf("should contain name:\n%s", result)
	}
	if strings.Contains(result, "## Flags") {
		t.Errorf("should not contain Flags:\n%s", result)
	}
}

// ---------------------------------------------------------------------------
// Benchmark
// ---------------------------------------------------------------------------

func BenchmarkMarkdownRender(b *testing.B) {
	ty := newTestTypography()
	cmd := fullCommand()
	for b.Loop() {
		_ = Render(ty, cmd, WithStyle(StyleMarkdown))
	}
}
