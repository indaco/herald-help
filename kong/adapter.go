// Package heraldkong converts a kong.Kong application into a
// heraldhelp.Command for styled help rendering with herald.
//
// Quick start:
//
//	parser, _ := kong.New(&cli)
//	kong.Help(func(options kong.HelpOptions, ctx *kong.Context) error {
//	    ty := herald.New()
//	    return heraldhelp.RenderTo(os.Stdout, ty, heraldkong.FromKong(parser))
//	})
package heraldkong

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	heraldhelp "github.com/indaco/herald-help"
)

// FromKong converts a kong.Kong application into a heraldhelp.Command.
// It extracts flags (with groups and env vars), positional arguments,
// subcommands, and other metadata from the kong model.
func FromKong(app *kong.Kong) heraldhelp.Command {
	return fromNode(app.Model.Node)
}

// FromNode converts a kong.Node (command or application) into a
// heraldhelp.Command. This is useful for rendering help for subcommands.
func FromNode(node *kong.Node) heraldhelp.Command {
	return fromNode(node)
}

// fromNode does the actual conversion from a kong.Node.
func fromNode(node *kong.Node) heraldhelp.Command {
	hc := heraldhelp.Command{
		Name:        node.Name,
		Description: detailOrHelp(node),
	}

	if len(node.Aliases) > 0 {
		hc.Aliases = node.Aliases
	}

	// Synopsis: build from name + flags + positional args + children.
	hc.Synopsis = buildSynopsis(node)

	// Flags.
	convertKongFlags(node, &hc)

	// Positional arguments.
	for _, p := range node.Positional {
		hc.Args = append(hc.Args, heraldhelp.Arg{
			Name:     "<" + p.Name + ">",
			Desc:     p.Help,
			Required: p.Required,
			Default:  p.Default,
		})
	}

	// Subcommands.
	convertKongChildren(node, &hc)

	return hc
}

// detailOrHelp returns the detailed help if available, otherwise the short help.
func detailOrHelp(node *kong.Node) string {
	if node.Detail != "" {
		return node.Detail
	}
	return node.Help
}

// buildSynopsis constructs a usage line from the node structure.
func buildSynopsis(node *kong.Node) string {
	parts := []string{node.Name}
	if len(node.Flags) > 0 {
		parts = append(parts, "[flags]")
	}
	for _, p := range node.Positional {
		if p.Required {
			parts = append(parts, "<"+p.Name+">")
		} else {
			parts = append(parts, "["+p.Name+"]")
		}
	}
	if len(node.Children) > 0 {
		parts = append(parts, "<command>")
	}
	return strings.Join(parts, " ")
}

// convertKongFlags extracts flags from the node, grouping by kong Group
// when groups are present.
func convertKongFlags(node *kong.Node, hc *heraldhelp.Command) {
	groups := make(map[string][]heraldhelp.Flag)
	var ungrouped []heraldhelp.Flag

	for _, f := range node.Flags {
		// Skip the help flag (auto-added by kong).
		if f.Name == "help" {
			continue
		}

		hf := convertKongFlag(f)

		if f.Group != nil && f.Group.Title != "" {
			groups[f.Group.Title] = append(groups[f.Group.Title], hf)
		} else {
			ungrouped = append(ungrouped, hf)
		}
	}

	hc.Flags = ungrouped
	for title, flags := range groups {
		hc.FlagGroups = append(hc.FlagGroups, heraldhelp.FlagGroup{
			Name:  title,
			Flags: flags,
		})
	}
}

// convertKongFlag converts a single kong.Flag to a heraldhelp.Flag.
func convertKongFlag(f *kong.Flag) heraldhelp.Flag {
	hf := heraldhelp.Flag{
		Long:     "--" + f.Name,
		Type:     kongFlagType(f),
		Default:  f.Default,
		Desc:     f.Help,
		Required: f.Required,
		Hidden:   f.Hidden,
		EnvVars:  f.Envs,
	}

	if f.Short != 0 {
		hf.Short = fmt.Sprintf("-%c", f.Short)
	}

	if f.Enum != "" {
		hf.Enum = strings.Split(f.Enum, ",")
	}

	return hf
}

// kongFlagType returns a human-readable type name for a kong flag.
func kongFlagType(f *kong.Flag) string {
	if f.IsBool() {
		return "bool"
	}
	if f.IsCounter() {
		return "counter"
	}
	// Use the reflect type of the target field.
	return f.Target.Type().String()
}

// convertKongChildren extracts child commands/arguments, grouping by kong
// Group when groups are present.
func convertKongChildren(node *kong.Node, hc *heraldhelp.Command) {
	groups := make(map[string][]heraldhelp.CommandRef)
	var ungrouped []heraldhelp.CommandRef

	for _, child := range node.Children {
		if child == nil || child.Hidden {
			continue
		}
		// Skip argument nodes (they're positional, not subcommands).
		if child.Type == kong.ArgumentNode {
			continue
		}

		ref := heraldhelp.CommandRef{
			Name:    child.Name,
			Aliases: child.Aliases,
			Desc:    child.Help,
		}

		if child.Group != nil && child.Group.Title != "" {
			groups[child.Group.Title] = append(groups[child.Group.Title], ref)
		} else {
			ungrouped = append(ungrouped, ref)
		}
	}

	hc.Commands = ungrouped
	for title, cmds := range groups {
		hc.CommandGroups = append(hc.CommandGroups, heraldhelp.CommandGroup{
			Name:     title,
			Commands: cmds,
		})
	}
}
