// Package heraldurfave converts a urfave/cli/v3 Command into a
// heraldhelp.Command for styled help rendering with herald.
//
// Quick start:
//
//	app := &cli.Command{
//	    Name: "myapp",
//	    Action: func(ctx context.Context, cmd *cli.Command) error {
//	        ty := herald.New()
//	        return heraldhelp.RenderTo(cmd.Writer, ty, heraldurfave.FromUrfave(cmd))
//	    },
//	}
package heraldurfave

import (
	"fmt"
	"sort"

	heraldhelp "github.com/indaco/herald-help"
	"github.com/urfave/cli/v3"
)

// FromUrfave converts a urfave/cli/v3 Command into a heraldhelp.Command.
// It extracts flags (with categories and env vars), subcommands (with
// categories), and other metadata.
func FromUrfave(cmd *cli.Command) heraldhelp.Command {
	hc := heraldhelp.Command{
		Name:        cmd.Name,
		Synopsis:    cmd.UsageText,
		Description: descOrUsage(cmd),
	}

	if len(cmd.Aliases) > 0 {
		hc.Aliases = cmd.Aliases
	}

	// Flags.
	convertUrfaveFlags(cmd, &hc)

	// Subcommands.
	convertUrfaveCommands(cmd, &hc)

	return hc
}

// descOrUsage returns the description if available, otherwise the usage.
func descOrUsage(cmd *cli.Command) string {
	if cmd.Description != "" {
		return cmd.Description
	}
	return cmd.Usage
}

// convertUrfaveFlags extracts flags from the command, grouping by category
// when categories are present.
func convertUrfaveFlags(cmd *cli.Command, hc *heraldhelp.Command) {
	categories := make(map[string][]heraldhelp.Flag)
	var uncategorized []heraldhelp.Flag

	for _, f := range cmd.Flags {
		hf := convertOneFlag(f)
		cat := flagCategory(f)
		if cat != "" {
			categories[cat] = append(categories[cat], hf)
		} else {
			uncategorized = append(uncategorized, hf)
		}
	}

	hc.Flags = uncategorized
	catKeys := make([]string, 0, len(categories))
	for k := range categories {
		catKeys = append(catKeys, k)
	}
	sort.Strings(catKeys)
	for _, cat := range catKeys {
		flags := categories[cat]
		hc.FlagGroups = append(hc.FlagGroups, heraldhelp.FlagGroup{
			Name:  cat,
			Flags: flags,
		})
	}
}

// convertOneFlag converts a single urfave flag to a heraldhelp.Flag.
func convertOneFlag(f cli.Flag) heraldhelp.Flag {
	names := f.Names()
	hf := heraldhelp.Flag{
		Hidden: isHidden(f),
	}

	for _, n := range names {
		if len(n) == 1 {
			hf.Short = "-" + n
		} else if hf.Long == "" {
			hf.Long = "--" + n
		}
	}

	// Extract type-specific details using type assertions.
	populateFlagDetails(f, &hf)

	return hf
}

// isHidden checks if a flag is hidden via type assertions on concrete types.
func isHidden(f cli.Flag) bool {
	switch tf := f.(type) {
	case *cli.StringFlag:
		return tf.Hidden
	case *cli.BoolFlag:
		return tf.Hidden
	case *cli.IntFlag:
		return tf.Hidden
	case *cli.FloatFlag:
		return tf.Hidden
	case *cli.DurationFlag:
		return tf.Hidden
	case *cli.StringSliceFlag:
		return tf.Hidden
	default:
		return false
	}
}

// populateFlagDetails fills in type, default, description, env vars, and
// required status from concrete flag types.
func populateFlagDetails(f cli.Flag, hf *heraldhelp.Flag) {
	switch tf := f.(type) {
	case *cli.StringFlag:
		hf.Type = "string"
		hf.Default = tf.Value
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	case *cli.BoolFlag:
		hf.Type = "bool"
		if tf.Value {
			hf.Default = "true"
		} else {
			hf.Default = "false"
		}
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	case *cli.IntFlag:
		hf.Type = "int"
		hf.Default = fmt.Sprintf("%d", tf.Value)
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	case *cli.FloatFlag:
		hf.Type = "float64"
		hf.Default = fmt.Sprintf("%g", tf.Value)
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	case *cli.DurationFlag:
		hf.Type = "duration"
		hf.Default = tf.Value.String()
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	case *cli.StringSliceFlag:
		hf.Type = "[]string"
		hf.Desc = tf.Usage
		hf.EnvVars = tf.Sources.EnvKeys()
		hf.Required = tf.Required
	default:
		hf.Type = "string"
	}
}

// flagCategory extracts the category from a urfave flag.
func flagCategory(f cli.Flag) string {
	switch tf := f.(type) {
	case *cli.StringFlag:
		return tf.Category
	case *cli.BoolFlag:
		return tf.Category
	case *cli.IntFlag:
		return tf.Category
	case *cli.FloatFlag:
		return tf.Category
	case *cli.DurationFlag:
		return tf.Category
	case *cli.StringSliceFlag:
		return tf.Category
	default:
		return ""
	}
}

// convertUrfaveCommands extracts subcommands, grouping by category when
// categories are present.
func convertUrfaveCommands(cmd *cli.Command, hc *heraldhelp.Command) {
	categories := make(map[string][]heraldhelp.CommandRef)
	var uncategorized []heraldhelp.CommandRef

	for _, sub := range cmd.Commands {
		if sub.Hidden {
			continue
		}
		ref := heraldhelp.CommandRef{
			Name:    sub.Name,
			Aliases: sub.Aliases,
			Desc:    sub.Usage,
		}
		if sub.Category != "" {
			categories[sub.Category] = append(categories[sub.Category], ref)
		} else {
			uncategorized = append(uncategorized, ref)
		}
	}

	hc.Commands = uncategorized
	catKeys := make([]string, 0, len(categories))
	for k := range categories {
		catKeys = append(catKeys, k)
	}
	sort.Strings(catKeys)
	for _, cat := range catKeys {
		cmds := categories[cat]
		hc.CommandGroups = append(hc.CommandGroups, heraldhelp.CommandGroup{
			Name:     cat,
			Commands: cmds,
		})
	}
}
