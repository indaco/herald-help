package heraldhelp

import (
	"strings"

	"github.com/indaco/herald"
)

// compactIndent is the left indentation for content lines in compact style.
const compactIndent = "  "

// renderSectionCompact dispatches a section using the compact style.
func renderSectionCompact(ty *herald.Typography, cmd Command, sec Section, cfg *RenderConfig) string {
	switch sec {
	case SectionName:
		return compactName(ty, cmd)
	case SectionDeprecated:
		return compactDeprecated(ty, cmd)
	case SectionSynopsis:
		return compactSynopsis(ty, cmd)
	case SectionDescription:
		return compactDescription(ty, cmd)
	case SectionArgs:
		return compactArgs(ty, cmd)
	case SectionFlags:
		return compactFlags(ty, cmd, cfg)
	case SectionInheritedFlags:
		return compactInheritedFlags(ty, cmd, cfg)
	case SectionCommands:
		return compactCommands(ty, cmd)
	case SectionExamples:
		return compactExamples(ty, cmd)
	case SectionSeeAlso:
		return compactSeeAlso(ty, cmd)
	case SectionFooter:
		return compactFooter(ty, cmd)
	default:
		return ""
	}
}

// compactHeading renders a section heading as uppercase colored text using
// the theme's H2 style (bold + secondary color) applied without decoration.
func compactHeading(ty *herald.Typography, text string) string {
	return ty.Theme().H2.UnsetMarginBottom().Render(strings.ToUpper(text))
}

// compactName renders the command name using the theme's H1 style (bold +
// primary color) without the underline decoration.
func compactName(ty *herald.Typography, cmd Command) string {
	if cmd.Name == "" {
		return ""
	}
	return ty.Theme().H1.UnsetMarginBottom().Render(cmd.Name)
}

// compactDeprecated renders a deprecation warning.
func compactDeprecated(ty *herald.Typography, cmd Command) string {
	if cmd.Deprecated == "" {
		return ""
	}
	return ty.Alert(herald.AlertWarning, "Deprecated: "+cmd.Deprecated)
}

// compactSynopsis renders the usage synopsis indented under a heading.
func compactSynopsis(ty *herald.Typography, cmd Command) string {
	if cmd.Synopsis == "" {
		return ""
	}
	return compactHeading(ty, "Usage") + "\n" + compactIndent + ty.Code(cmd.Synopsis)
}

// compactDescription renders the description as a paragraph.
func compactDescription(ty *herald.Typography, cmd Command) string {
	if cmd.Description == "" {
		return ""
	}
	return ty.P(cmd.Description)
}

// compactArgs renders positional arguments as an indented KV list.
func compactArgs(ty *herald.Typography, cmd Command) string {
	if len(cmd.Args) == 0 {
		return ""
	}

	pairs := make([][2]string, len(cmd.Args))
	for i, a := range cmd.Args {
		desc := a.Desc
		if a.Required {
			desc += " " + ty.Bold("(required)")
		}
		if a.Default != "" {
			desc += " " + ty.Small("(default: "+a.Default+")")
		}
		pairs[i] = [2]string{ty.Var(a.Name), desc}
	}

	return compactHeading(ty, "Arguments") + "\n" +
		ty.KVGroupWithOpts(pairs, compactKVOpts()...)
}

// compactFlags renders non-inherited flags as an indented two-column list.
func compactFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	flatFlags := filterFlags(cmd.Flags, cfg.ShowHidden, false)

	var groupBlocks []string
	for _, g := range cmd.FlagGroups {
		filtered := filterFlags(g.Flags, cfg.ShowHidden, false)
		if len(filtered) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, compactHeading(ty, g.Name)+"\n"+compactFlagList(ty, filtered, cfg))
	}

	if len(flatFlags) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := make([]string, 0, 1+len(groupBlocks))
	heading := compactHeading(ty, "Flags")
	if len(flatFlags) > 0 {
		heading += "\n" + compactFlagList(ty, flatFlags, cfg)
	}
	parts = append(parts, heading)
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n\n")
}

// compactInheritedFlags renders inherited flags in a compact section.
func compactInheritedFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	inherited := collectInheritedFlags(cmd, cfg.ShowHidden)
	if len(inherited) == 0 {
		return ""
	}

	return compactHeading(ty, "Inherited Flags") + "\n" + compactFlagList(ty, inherited, cfg)
}

// compactFlagList renders a slice of flags as an indented KV list using
// compact-style options.
func compactFlagList(ty *herald.Typography, flags []Flag, cfg *RenderConfig) string {
	return buildKVFlagList(ty, flags, cfg, compactKVOpts())
}

// compactCommands renders subcommands as an indented two-column list.
func compactCommands(ty *herald.Typography, cmd Command) string {
	var groupBlocks []string
	for _, g := range cmd.CommandGroups {
		if len(g.Commands) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, compactHeading(ty, g.Name)+"\n"+compactCommandList(ty, g.Commands))
	}

	if len(cmd.Commands) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := make([]string, 0, 1+len(groupBlocks))
	heading := compactHeading(ty, "Commands")
	if len(cmd.Commands) > 0 {
		heading += "\n" + compactCommandList(ty, cmd.Commands)
	}
	parts = append(parts, heading)
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n\n")
}

// compactCommandList renders a slice of command refs as an indented KV list
// using compact-style options.
func compactCommandList(ty *herald.Typography, cmds []CommandRef) string {
	return buildKVCommandList(ty, cmds, compactKVOpts())
}

// compactExamples renders examples as an indented KV list with descriptions
// as keys and commands as values.
func compactExamples(ty *herald.Typography, cmd Command) string {
	if len(cmd.Examples) == 0 {
		return ""
	}

	pairs := make([][2]string, 0, len(cmd.Examples))
	for _, ex := range cmd.Examples {
		key := ty.Small(ex.Desc)
		val := ty.Code(ex.Command)
		pairs = append(pairs, [2]string{key, val})
	}

	return compactHeading(ty, "Examples") + "\n" +
		ty.KVGroupWithOpts(pairs, compactKVOpts()...)
}

// compactSeeAlso renders related references as an indented list.
func compactSeeAlso(ty *herald.Typography, cmd Command) string {
	if len(cmd.SeeAlso) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(compactHeading(ty, "See Also"))

	for _, s := range cmd.SeeAlso {
		sb.WriteByte('\n')
		sb.WriteString(compactIndent + ty.Italic(s))
	}

	return sb.String()
}

// compactFooter renders footer text as small/faint text.
func compactFooter(ty *herald.Typography, cmd Command) string {
	if cmd.Footer == "" {
		return ""
	}
	return ty.Small(cmd.Footer)
}

// compactKVOpts returns the standard KV options for compact-style rendering:
// no separator, pre-styled keys and values, 2-space indent.
func compactKVOpts() []herald.KVOption {
	return []herald.KVOption{
		herald.WithKVGroupSeparator(""),
		herald.WithKVRawKeys(true),
		herald.WithKVRawValues(true),
		herald.WithKVIndent(2),
	}
}
