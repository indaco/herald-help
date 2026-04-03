package heraldhelp

import (
	"strings"
)

// renderSectionMarkdown dispatches a section using the Markdown style.
// The output is valid Markdown text with no ANSI escape sequences.
func renderSectionMarkdown(cmd Command, sec Section, cfg *RenderConfig) string {
	switch sec {
	case SectionName:
		return mdName(cmd)
	case SectionDeprecated:
		return mdDeprecated(cmd)
	case SectionSynopsis:
		return mdSynopsis(cmd)
	case SectionDescription:
		return mdDescription(cmd)
	case SectionArgs:
		return mdArgs(cmd)
	case SectionFlags:
		return mdFlags(cmd, cfg)
	case SectionInheritedFlags:
		return mdInheritedFlags(cmd, cfg)
	case SectionCommands:
		return mdCommands(cmd)
	case SectionExamples:
		return mdExamples(cmd)
	case SectionSeeAlso:
		return mdSeeAlso(cmd)
	case SectionFooter:
		return mdFooter(cmd)
	default:
		return ""
	}
}

func mdName(cmd Command) string {
	if cmd.Name == "" {
		return ""
	}
	return "# " + cmd.Name
}

func mdDeprecated(cmd Command) string {
	if cmd.Deprecated == "" {
		return ""
	}
	return "> **Warning:** Deprecated: " + cmd.Deprecated
}

func mdSynopsis(cmd Command) string {
	if cmd.Synopsis == "" {
		return ""
	}
	return "```\n" + cmd.Synopsis + "\n```"
}

func mdDescription(cmd Command) string {
	if cmd.Description == "" {
		return ""
	}
	return cmd.Description
}

func mdArgs(cmd Command) string {
	if len(cmd.Args) == 0 {
		return ""
	}

	rows := make([][]string, 0, 1+len(cmd.Args))
	rows = append(rows, []string{"Argument", "Description", "Required", "Default"})
	for _, a := range cmd.Args {
		req := ""
		if a.Required {
			req = "yes"
		}
		rows = append(rows, []string{"`" + a.Name + "`", a.Desc, req, a.Default})
	}

	return "## Arguments\n\n" + mdTable(rows)
}

func mdFlags(cmd Command, cfg *RenderConfig) string {
	flatFlags := filterFlags(cmd.Flags, cfg.ShowHidden, false)

	var groupBlocks []string
	for _, g := range cmd.FlagGroups {
		filtered := filterFlags(g.Flags, cfg.ShowHidden, false)
		if len(filtered) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, "### "+g.Name+"\n\n"+mdFlagTable(filtered, cfg))
	}

	if len(flatFlags) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := []string{"## Flags"}
	if len(flatFlags) > 0 {
		parts = append(parts, mdFlagTable(flatFlags, cfg))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n\n")
}

func mdInheritedFlags(cmd Command, cfg *RenderConfig) string {
	inherited := filterFlags(cmd.Flags, cfg.ShowHidden, true)
	for _, g := range cmd.FlagGroups {
		inherited = append(inherited, filterFlags(g.Flags, cfg.ShowHidden, true)...)
	}

	if len(inherited) == 0 {
		return ""
	}

	return "### Inherited Flags\n\n" + mdFlagTable(inherited, cfg)
}

func mdFlagTable(flags []Flag, cfg *RenderConfig) string {
	rows := make([][]string, 0, 1+len(flags))
	rows = append(rows, []string{"Flag", "Type", "Default", "Description"})

	for i := range flags {
		f := &flags[i]
		name := "`" + strings.TrimLeft(formatFlagName(*f), " ") + "`"
		desc := f.Desc
		if cfg.EnvVarDisplay && len(f.EnvVars) > 0 {
			desc += " [`" + formatEnvVars(f.EnvVars) + "`]"
		}
		if f.Required {
			desc += " **(required)**"
		}
		if f.Deprecated != "" {
			desc += " *(DEPRECATED: " + f.Deprecated + ")*"
		}
		if len(f.Enum) > 0 {
			desc += " [enum: `" + strings.Join(f.Enum, "`, `") + "`]"
		}
		rows = append(rows, []string{name, f.Type, f.Default, desc})
	}

	return mdTable(rows)
}

func mdCommands(cmd Command) string {
	var groupBlocks []string
	for _, g := range cmd.CommandGroups {
		if len(g.Commands) == 0 {
			continue
		}
		groupBlocks = append(groupBlocks, "### "+g.Name+"\n\n"+mdCommandTable(g.Commands))
	}

	if len(cmd.Commands) == 0 && len(groupBlocks) == 0 {
		return ""
	}

	parts := []string{"## Commands"}
	if len(cmd.Commands) > 0 {
		parts = append(parts, mdCommandTable(cmd.Commands))
	}
	parts = append(parts, groupBlocks...)

	return strings.Join(parts, "\n\n")
}

func mdCommandTable(cmds []CommandRef) string {
	rows := make([][]string, 0, 1+len(cmds))
	rows = append(rows, []string{"Command", "Aliases", "Description"})

	for _, c := range cmds {
		aliases := strings.Join(c.Aliases, ", ")
		rows = append(rows, []string{"`" + c.Name + "`", aliases, c.Desc})
	}

	return mdTable(rows)
}

func mdExamples(cmd Command) string {
	if len(cmd.Examples) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## Examples")

	for _, ex := range cmd.Examples {
		sb.WriteString("\n\n")
		if ex.Desc != "" {
			sb.WriteString(ex.Desc + "\n\n")
		}
		if ex.Command != "" {
			sb.WriteString("```\n" + ex.Command + "\n```")
		}
	}

	return sb.String()
}

func mdSeeAlso(cmd Command) string {
	if len(cmd.SeeAlso) == 0 {
		return ""
	}

	items := make([]string, len(cmd.SeeAlso))
	for i, s := range cmd.SeeAlso {
		items[i] = "- " + s
	}

	return "## See Also\n\n" + strings.Join(items, "\n")
}

func mdFooter(cmd Command) string {
	if cmd.Footer == "" {
		return ""
	}
	return "---\n\n*" + cmd.Footer + "*"
}

// mdTable renders a well-formatted Markdown table with aligned columns.
// The first row is the header. All cells are padded to uniform column widths.
func mdTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	// Determine column count from header.
	cols := len(rows[0])

	// Find max width for each column.
	widths := make([]int, cols)
	for _, row := range rows {
		for c := 0; c < cols && c < len(row); c++ {
			if len(row[c]) > widths[c] {
				widths[c] = len(row[c])
			}
		}
	}

	// Build header row.
	lines := make([]string, 0, len(rows)+1)
	lines = append(lines, mdTableRow(rows[0], widths))

	// Build separator row.
	sep := make([]string, cols)
	for c, w := range widths {
		sep[c] = strings.Repeat("-", w)
	}
	lines = append(lines, "| "+strings.Join(sep, " | ")+" |")

	// Build data rows.
	for _, row := range rows[1:] {
		lines = append(lines, mdTableRow(row, widths))
	}

	return strings.Join(lines, "\n")
}

// mdTableRow formats a single row with padded cells.
func mdTableRow(row []string, widths []int) string {
	cells := make([]string, len(widths))
	for c, w := range widths {
		val := ""
		if c < len(row) {
			val = row[c]
		}
		cells[c] = val + strings.Repeat(" ", w-len(val))
	}
	return "| " + strings.Join(cells, " | ") + " |"
}
