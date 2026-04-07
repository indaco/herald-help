package heraldhelp

import (
	"strings"

	"github.com/indaco/herald"
)

// buildKVFlagList renders a slice of flags as a KV list with colored flag
// names, faint types/defaults, and highlighted env vars. Flag names and types
// are independently padded for clean column alignment. The kvOpts parameter
// controls style-specific options (indent, separator, etc.).
func buildKVFlagList(ty *herald.Typography, flags []Flag, cfg *RenderConfig, kvOpts []herald.KVOption) string {
	// First pass: compute max widths for name and type columns.
	maxNameW := 0
	maxTypeW := 0
	for i := range flags {
		if w := len(formatFlagName(flags[i])); w > maxNameW {
			maxNameW = w
		}
		ft := flags[i].Type
		if ft == "bool" {
			ft = ""
		}
		if len(ft) > maxTypeW {
			maxTypeW = len(ft)
		}
	}

	// Second pass: build aligned keys.
	pairs := make([][2]string, 0, len(flags))
	for i := range flags {
		f := &flags[i]
		name := formatFlagName(*f)
		namePad := strings.Repeat(" ", maxNameW-len(name))

		ft := f.Type
		if ft == "bool" {
			ft = ""
		}
		typePad := strings.Repeat(" ", maxTypeW-len(ft))

		key := ty.Var(name) + namePad
		if ft != "" {
			key += " " + ty.Small(ft) + typePad + "  "
		} else if maxTypeW > 0 {
			key += " " + strings.Repeat(" ", maxTypeW) + "  "
		}

		// Build description with colored metadata.
		desc := f.Desc
		if f.Default != "" {
			desc += " " + ty.Small("(default: "+f.Default+")")
		}
		if cfg.EnvVarDisplay && len(f.EnvVars) > 0 {
			desc += " " + ty.Kbd(formatEnvVars(f.EnvVars))
		}
		if f.Required {
			desc += " " + ty.Bold("(required)")
		}
		if f.Deprecated != "" {
			desc += " " + ty.Bold("(DEPRECATED: "+f.Deprecated+")")
		}
		if len(f.Enum) > 0 {
			desc += " " + ty.Small("[enum: "+strings.Join(f.Enum, ", ")+"]")
		}
		pairs = append(pairs, [2]string{key, desc})
	}

	return ty.KVGroupWithOpts(pairs, kvOpts...)
}

// buildKVCommandList renders a slice of command refs as a KV list with colored
// command names. Name and alias columns are independently aligned. The kvOpts
// parameter controls style-specific options.
func buildKVCommandList(ty *herald.Typography, cmds []CommandRef, kvOpts []herald.KVOption) string {
	// Compute max widths for name and alias columns.
	maxNameW := 0
	maxAliasW := 0
	hasAliases := false
	for _, c := range cmds {
		if len(c.Name) > maxNameW {
			maxNameW = len(c.Name)
		}
		a := strings.Join(c.Aliases, ", ")
		if a != "" {
			hasAliases = true
		}
		if len(a) > maxAliasW {
			maxAliasW = len(a)
		}
	}

	pairs := make([][2]string, len(cmds))
	for i, c := range cmds {
		a := strings.Join(c.Aliases, ", ")

		// Name with comma attached when aliases exist.
		key := ty.Var(c.Name)
		if a != "" {
			key += ty.Small(",")
		} else if hasAliases {
			key += " " // space in place of comma for alignment
		}
		key += strings.Repeat(" ", maxNameW-len(c.Name))

		// Alias column, padded independently.
		if hasAliases {
			key += " " + ty.Small(a) + strings.Repeat(" ", maxAliasW-len(a)) + "  "
		}

		pairs[i] = [2]string{key, c.Desc}
	}
	return ty.KVGroupWithOpts(pairs, kvOpts...)
}
