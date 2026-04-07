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
		order = DefaultSectionOrder()
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

// collectInheritedFlags gathers inherited flags from both flat flags and flag
// groups, filtering by hidden status. This centralises the logic used by all
// four render styles.
func collectInheritedFlags(cmd Command, showHidden bool) []Flag {
	inherited := filterFlags(cmd.Flags, showHidden, true)
	for _, g := range cmd.FlagGroups {
		inherited = append(inherited, filterFlags(g.Flags, showHidden, true)...)
	}
	return inherited
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

// FormatVersion returns a formatted version string suitable for the Footer field.
func FormatVersion(name, version string) string {
	return fmt.Sprintf("%s version %s", name, version)
}
