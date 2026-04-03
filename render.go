package heraldhelp

import (
	"fmt"
	"io"
	"strings"

	"github.com/indaco/herald"
)

// Render renders the help output for a Command using the given Typography
// instance and returns it as a string. If ty is nil, an empty string is
// returned.
func Render(ty *herald.Typography, cmd Command, opts ...RenderOption) string {
	if ty == nil {
		return ""
	}

	cfg := buildConfig(opts)
	order := cfg.SectionOrder
	if len(order) == 0 {
		order = DefaultSectionOrder
	}

	blocks := make([]string, 0, len(order))
	for _, sec := range order {
		if s := renderSection(ty, cmd, sec, cfg); s != "" {
			blocks = append(blocks, s)
		}
	}

	if cfg.Style == StyleMarkdown {
		return strings.Join(blocks, "\n\n")
	}
	return ty.Compose(blocks...)
}

// RenderTo renders the help output and writes it to w.
func RenderTo(w io.Writer, ty *herald.Typography, cmd Command, opts ...RenderOption) error {
	s := Render(ty, cmd, opts...)
	_, err := io.WriteString(w, s+"\n")
	return err
}

// buildConfig applies functional options to a default RenderConfig.
func buildConfig(opts []RenderOption) *RenderConfig {
	cfg := &RenderConfig{
		EnvVarDisplay: true,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// renderSection dispatches a single section to its renderer.
func renderSection(ty *herald.Typography, cmd Command, sec Section, cfg *RenderConfig) string {
	switch cfg.Style {
	case StyleRich:
		return renderSectionRich(ty, cmd, sec, cfg)
	case StyleGrouped:
		return renderSectionGrouped(ty, cmd, sec, cfg)
	case StyleMarkdown:
		return renderSectionMarkdown(cmd, sec, cfg)
	default:
		return renderSectionCompact(ty, cmd, sec, cfg)
	}
}

// renderSectionRich dispatches a section using the rich (default) style.
func renderSectionRich(ty *herald.Typography, cmd Command, sec Section, cfg *RenderConfig) string {
	switch sec {
	case SectionName:
		return renderName(ty, cmd)
	case SectionDeprecated:
		return renderDeprecated(ty, cmd)
	case SectionSynopsis:
		return renderSynopsis(ty, cmd)
	case SectionDescription:
		return renderDescription(ty, cmd)
	case SectionArgs:
		return renderArgs(ty, cmd)
	case SectionFlags:
		return renderFlags(ty, cmd, cfg)
	case SectionInheritedFlags:
		return renderInheritedFlags(ty, cmd, cfg)
	case SectionCommands:
		return renderCommands(ty, cmd)
	case SectionExamples:
		return renderExamples(ty, cmd)
	case SectionSeeAlso:
		return renderSeeAlso(ty, cmd)
	case SectionFooter:
		return renderFooter(ty, cmd)
	default:
		return ""
	}
}

// renderName renders the command name as H1.
func renderName(ty *herald.Typography, cmd Command) string {
	if cmd.Name == "" {
		return ""
	}
	return ty.H1(cmd.Name)
}

// renderDeprecated renders a deprecation warning alert.
func renderDeprecated(ty *herald.Typography, cmd Command) string {
	if cmd.Deprecated == "" {
		return ""
	}
	return ty.Alert(herald.AlertWarning, "Deprecated: "+cmd.Deprecated)
}

// renderSynopsis renders the usage synopsis as a code block.
func renderSynopsis(ty *herald.Typography, cmd Command) string {
	if cmd.Synopsis == "" {
		return ""
	}
	return ty.CodeBlock(cmd.Synopsis)
}

// renderDescription renders the long description as a paragraph.
func renderDescription(ty *herald.Typography, cmd Command) string {
	if cmd.Description == "" {
		return ""
	}
	return ty.P(cmd.Description)
}

// renderArgs renders positional arguments as a table under an H2 heading.
func renderArgs(ty *herald.Typography, cmd Command) string {
	if len(cmd.Args) == 0 {
		return ""
	}

	rows := make([][]string, 0, len(cmd.Args)+1)
	rows = append(rows, []string{"Argument", "Description", "Required", "Default"})

	for _, arg := range cmd.Args {
		req := ""
		if arg.Required {
			req = "yes"
		}
		rows = append(rows, []string{arg.Name, arg.Desc, req, arg.Default})
	}

	return ty.H2("Arguments") + "\n" + ty.Table(rows)
}

// renderFlags renders flags (flat or grouped) as tables.
func renderFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	// Collect non-inherited flat flags.
	flatFlags := filterFlags(cmd.Flags, cfg.ShowHidden, false)

	// Collect grouped flags (non-inherited).
	var groupBlocks []string
	for _, g := range cmd.FlagGroups {
		filtered := filterFlags(g.Flags, cfg.ShowHidden, false)
		if len(filtered) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, ty.H3(g.Name)+"\n"+buildFlagTable(ty, filtered, cfg))
	}

	if len(flatFlags) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := []string{ty.H2("Flags")}
	if len(flatFlags) > 0 {
		parts = append(parts, buildFlagTable(ty, flatFlags, cfg))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n")
}

// renderInheritedFlags renders inherited/persistent flags in a separate section.
func renderInheritedFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	inherited := filterFlags(cmd.Flags, cfg.ShowHidden, true)
	// Also check flag groups for inherited flags.
	for _, g := range cmd.FlagGroups {
		inherited = append(inherited, filterFlags(g.Flags, cfg.ShowHidden, true)...)
	}

	if len(inherited) == 0 {
		return ""
	}

	return ty.H3("Inherited Flags") + "\n" + buildFlagTable(ty, inherited, cfg)
}

// filterFlags returns flags matching the inherited status, excluding hidden
// flags unless showHidden is true.
func filterFlags(flags []Flag, showHidden, inherited bool) []Flag {
	var out []Flag
	for i := range flags {
		if flags[i].Inherited != inherited {
			continue
		}
		if flags[i].Hidden && !showHidden {
			continue
		}
		out = append(out, flags[i])
	}
	return out
}

// buildFlagTable builds a table of flags.
func buildFlagTable(ty *herald.Typography, flags []Flag, cfg *RenderConfig) string {
	rows := make([][]string, 0, len(flags)+1)
	rows = append(rows, []string{"Flag", "Type", "Default", "Description"})

	for i := range flags {
		f := &flags[i]
		name := formatFlagName(*f)
		desc := f.Desc
		if cfg.EnvVarDisplay && len(f.EnvVars) > 0 {
			desc += " [" + formatEnvVars(f.EnvVars) + "]"
		}
		if f.Required {
			desc += " (required)"
		}
		if f.Deprecated != "" {
			desc += " (DEPRECATED: " + f.Deprecated + ")"
		}
		if len(f.Enum) > 0 {
			desc += " [enum: " + strings.Join(f.Enum, ", ") + "]"
		}
		rows = append(rows, []string{name, f.Type, f.Default, desc})
	}

	return ty.Table(rows)
}

// formatFlagName formats a flag name in GNU-style: -o, --output
func formatFlagName(f Flag) string {
	switch {
	case f.Short != "" && f.Long != "":
		return f.Short + ", " + f.Long
	case f.Long != "":
		return "    " + f.Long
	case f.Short != "":
		return f.Short
	default:
		return ""
	}
}

// formatEnvVars joins environment variable names with a dollar prefix.
func formatEnvVars(vars []string) string {
	parts := make([]string, len(vars))
	for i, v := range vars {
		parts[i] = "$" + v
	}
	return strings.Join(parts, ", ")
}

// renderCommands renders subcommands (flat or grouped) as tables.
func renderCommands(ty *herald.Typography, cmd Command) string {
	var groupBlocks []string
	for _, g := range cmd.CommandGroups {
		if len(g.Commands) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, ty.H3(g.Name)+"\n"+buildCommandTable(ty, g.Commands))
	}

	if len(cmd.Commands) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := []string{ty.H2("Commands")}
	if len(cmd.Commands) > 0 {
		parts = append(parts, buildCommandTable(ty, cmd.Commands))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n")
}

// buildCommandTable builds a table of subcommand references.
func buildCommandTable(ty *herald.Typography, cmds []CommandRef) string {
	rows := make([][]string, 0, len(cmds)+1)
	rows = append(rows, []string{"Command", "Aliases", "Description"})

	for _, c := range cmds {
		aliases := strings.Join(c.Aliases, ", ")
		rows = append(rows, []string{c.Name, aliases, c.Desc})
	}

	return ty.Table(rows)
}

// renderExamples renders usage examples with descriptions and code blocks.
func renderExamples(ty *herald.Typography, cmd Command) string {
	if len(cmd.Examples) == 0 {
		return ""
	}

	parts := []string{ty.H2("Examples")}
	for _, ex := range cmd.Examples {
		if ex.Desc != "" {
			parts = append(parts, ty.P(ex.Desc))
		}
		if ex.Command != "" {
			parts = append(parts, ty.CodeBlock(ex.Command))
		}
	}

	return strings.Join(parts, "\n")
}

// renderSeeAlso renders related commands/resources as an unordered list.
func renderSeeAlso(ty *herald.Typography, cmd Command) string {
	if len(cmd.SeeAlso) == 0 {
		return ""
	}
	return ty.H2("See Also") + "\n" + ty.UL(cmd.SeeAlso...)
}

// renderFooter renders footer text as small/faint text.
func renderFooter(ty *herald.Typography, cmd Command) string {
	if cmd.Footer == "" {
		return ""
	}
	return ty.Small(cmd.Footer)
}

// FormatVersion returns a formatted version string suitable for the Footer field.
func FormatVersion(name, version string) string {
	return fmt.Sprintf("%s version %s", name, version)
}
