// Package heraldcobra converts a cobra.Command into a heraldhelp.Command
// for styled help rendering with herald.
//
// Quick start:
//
//	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
//	    ty := herald.New()
//	    heraldhelp.RenderTo(cmd.OutOrStdout(), ty, heraldcobra.FromCobra(cmd))
//	})
package heraldcobra

import (
	"strings"

	heraldhelp "github.com/indaco/herald-help"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FromCobra converts a cobra.Command into a heraldhelp.Command. It extracts
// flags (local and inherited), subcommands, aliases, examples, and deprecated
// status from the cobra command.
func FromCobra(cmd *cobra.Command) heraldhelp.Command {
	hc := heraldhelp.Command{
		Name:        cmd.Name(),
		Synopsis:    cmd.UseLine(),
		Description: longOrShort(cmd),
		Aliases:     cmd.Aliases,
		Deprecated:  cmd.Deprecated,
	}

	// Local flags (non-inherited).
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		hc.Flags = append(hc.Flags, convertPflag(f, false))
	})

	// Inherited (persistent) flags.
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		hc.Flags = append(hc.Flags, convertPflag(f, true))
	})

	// Subcommands: include non-hidden children.
	// We avoid IsAvailableCommand() because it returns false for children of
	// deprecated commands, even though they are still valid.
	for _, sub := range cmd.Commands() {
		if !sub.Hidden {
			hc.Commands = append(hc.Commands, heraldhelp.CommandRef{
				Name:    sub.Name(),
				Aliases: sub.Aliases,
				Desc:    sub.Short,
			})
		}
	}

	// Examples.
	if cmd.Example != "" {
		hc.Examples = parseExamples(cmd.Example)
	}

	return hc
}

// longOrShort returns the long description if available, otherwise the short.
func longOrShort(cmd *cobra.Command) string {
	if cmd.Long != "" {
		return cmd.Long
	}
	return cmd.Short
}

// convertPflag converts a pflag.Flag to a heraldhelp.Flag.
func convertPflag(f *pflag.Flag, inherited bool) heraldhelp.Flag {
	hf := heraldhelp.Flag{
		Long:       "--" + f.Name,
		Type:       f.Value.Type(),
		Default:    f.DefValue,
		Desc:       f.Usage,
		Hidden:     f.Hidden,
		Deprecated: f.Deprecated,
		Inherited:  inherited,
	}

	if f.Shorthand != "" {
		hf.Short = "-" + f.Shorthand
	}

	// Check annotations for required.
	if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
		hf.Required = true
	}

	return hf
}

// parseExamples splits cobra's example string into heraldhelp.Example entries.
// Lines starting with "$" or indented lines are treated as commands; other
// non-empty lines are treated as descriptions.
func parseExamples(raw string) []heraldhelp.Example {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var examples []heraldhelp.Example
	var currentDesc string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if isCommandLine(trimmed) {
			cmd := strings.TrimPrefix(trimmed, "$ ")
			cmd = strings.TrimPrefix(cmd, "$")
			examples = append(examples, heraldhelp.Example{
				Desc:    strings.TrimSpace(currentDesc),
				Command: strings.TrimSpace(cmd),
			})
			currentDesc = ""
		} else {
			if currentDesc != "" {
				currentDesc += " "
			}
			currentDesc += trimmed
		}
	}

	// Trailing description with no command.
	if currentDesc != "" {
		examples = append(examples, heraldhelp.Example{
			Desc: strings.TrimSpace(currentDesc),
		})
	}

	return examples
}

// isCommandLine checks if a trimmed line looks like a command invocation.
func isCommandLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "$ ") || strings.HasPrefix(trimmed, "$\t")
}
