package heraldhelp

import (
	"strings"

	"github.com/indaco/herald"
)

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
	inherited := collectInheritedFlags(cmd, cfg.ShowHidden)
	if len(inherited) == 0 {
		return ""
	}

	return ty.H3("Inherited Flags") + "\n" + buildFlagTable(ty, inherited, cfg)
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
