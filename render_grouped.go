package heraldhelp

import (
	"strings"

	"github.com/indaco/herald"
)

// renderSectionGrouped dispatches a section using the grouped style.
// Each section is wrapped in a herald Fieldset with the section name as legend.
func renderSectionGrouped(ty *herald.Typography, cmd Command, sec Section, cfg *RenderConfig) string {
	switch sec {
	case SectionName:
		return groupedName(ty, cmd)
	case SectionDeprecated:
		return groupedDeprecated(ty, cmd)
	case SectionSynopsis:
		return groupedSynopsis(ty, cmd)
	case SectionDescription:
		return groupedDescription(ty, cmd)
	case SectionArgs:
		return groupedArgs(ty, cmd)
	case SectionFlags:
		return groupedFlags(ty, cmd, cfg)
	case SectionInheritedFlags:
		return groupedInheritedFlags(ty, cmd, cfg)
	case SectionCommands:
		return groupedCommands(ty, cmd)
	case SectionExamples:
		return groupedExamples(ty, cmd)
	case SectionSeeAlso:
		return groupedSeeAlso(ty, cmd)
	case SectionFooter:
		return groupedFooter(ty, cmd)
	default:
		return ""
	}
}

// groupedName renders the command name using H1 style without decoration.
func groupedName(ty *herald.Typography, cmd Command) string {
	if cmd.Name == "" {
		return ""
	}
	return ty.Theme().H1.UnsetMarginBottom().Render(cmd.Name)
}

// groupedDeprecated renders a deprecation warning alert.
func groupedDeprecated(ty *herald.Typography, cmd Command) string {
	if cmd.Deprecated == "" {
		return ""
	}
	return ty.Alert(herald.AlertWarning, "Deprecated: "+cmd.Deprecated)
}

// groupedSynopsis renders the synopsis in a fieldset.
func groupedSynopsis(ty *herald.Typography, cmd Command) string {
	if cmd.Synopsis == "" {
		return ""
	}
	return ty.Fieldset("Usage", ty.Code(cmd.Synopsis))
}

// groupedDescription renders the description as a paragraph.
func groupedDescription(ty *herald.Typography, cmd Command) string {
	if cmd.Description == "" {
		return ""
	}
	return cmd.Description
}

// groupedArgs renders positional arguments in a fieldset with KV layout.
func groupedArgs(ty *herald.Typography, cmd Command) string {
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

	return ty.Fieldset("Arguments", ty.KVGroupWithOpts(pairs, groupedKVOpts()...))
}

// groupedFlags renders flags in a fieldset with KV layout.
func groupedFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	flatFlags := filterFlags(cmd.Flags, cfg.ShowHidden, false)

	var groupBlocks []string
	for _, g := range cmd.FlagGroups {
		filtered := filterFlags(g.Flags, cfg.ShowHidden, false)
		if len(filtered) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, ty.Fieldset(g.Name, groupedFlagList(ty, filtered, cfg)))
	}

	if len(flatFlags) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := make([]string, 0, 1+len(groupBlocks))
	if len(flatFlags) > 0 {
		parts = append(parts, ty.Fieldset("Flags", groupedFlagList(ty, flatFlags, cfg)))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n")
}

// groupedInheritedFlags renders inherited flags in a fieldset.
func groupedInheritedFlags(ty *herald.Typography, cmd Command, cfg *RenderConfig) string {
	inherited := collectInheritedFlags(cmd, cfg.ShowHidden)
	if len(inherited) == 0 {
		return ""
	}

	return ty.Fieldset("Inherited Flags", groupedFlagList(ty, inherited, cfg))
}

// groupedFlagList builds the KV content for a flag fieldset using grouped-style
// options.
func groupedFlagList(ty *herald.Typography, flags []Flag, cfg *RenderConfig) string {
	return buildKVFlagList(ty, flags, cfg, groupedKVOpts())
}

// groupedCommands renders subcommands in a fieldset with KV layout.
func groupedCommands(ty *herald.Typography, cmd Command) string {
	var groupBlocks []string
	for _, g := range cmd.CommandGroups {
		if len(g.Commands) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, ty.Fieldset(g.Name, groupedCommandList(ty, g.Commands)))
	}

	if len(cmd.Commands) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := make([]string, 0, 1+len(groupBlocks))
	if len(cmd.Commands) > 0 {
		parts = append(parts, ty.Fieldset("Commands", groupedCommandList(ty, cmd.Commands)))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n")
}

// groupedCommandList builds the KV content for a command fieldset using
// grouped-style options.
func groupedCommandList(ty *herald.Typography, cmds []CommandRef) string {
	return buildKVCommandList(ty, cmds, groupedKVOpts())
}

// groupedExamples renders examples in a fieldset.
func groupedExamples(ty *herald.Typography, cmd Command) string {
	if len(cmd.Examples) == 0 {
		return ""
	}

	pairs := make([][2]string, 0, len(cmd.Examples))
	for _, ex := range cmd.Examples {
		pairs = append(pairs, [2]string{ty.Small(ex.Desc), ty.Code(ex.Command)})
	}

	return ty.Fieldset("Examples", ty.KVGroupWithOpts(pairs, groupedKVOpts()...))
}

// groupedSeeAlso renders see-also items in a fieldset.
func groupedSeeAlso(ty *herald.Typography, cmd Command) string {
	if len(cmd.SeeAlso) == 0 {
		return ""
	}

	items := make([]string, len(cmd.SeeAlso))
	for i, s := range cmd.SeeAlso {
		items[i] = ty.Italic(s)
	}

	return ty.Fieldset("See Also", strings.Join(items, "\n"))
}

// groupedFooter renders footer text as small/faint text.
func groupedFooter(ty *herald.Typography, cmd Command) string {
	if cmd.Footer == "" {
		return ""
	}
	return ty.Small(cmd.Footer)
}

// groupedKVOpts returns the standard KV options for grouped-style rendering.
func groupedKVOpts() []herald.KVOption {
	return []herald.KVOption{
		herald.WithKVGroupSeparator(""),
		herald.WithKVRawKeys(true),
		herald.WithKVRawValues(true),
	}
}
