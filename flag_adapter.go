package heraldhelp

import (
	"flag"
	"fmt"
	"strings"
)

// FromFlagSet creates a Command from a standard library *flag.FlagSet.
// The name parameter sets the command name. All defined flags are extracted
// with their types and defaults.
func FromFlagSet(name string, fs *flag.FlagSet) Command {
	cmd := Command{
		Name:     name,
		Synopsis: name + " [flags]",
	}

	fs.VisitAll(func(f *flag.Flag) {
		fl := Flag{
			Long:    "--" + f.Name,
			Desc:    f.Usage,
			Default: f.DefValue,
			Type:    flagType(f),
		}
		cmd.Flags = append(cmd.Flags, fl)
	})

	return cmd
}

// flagType returns a human-readable type name for a flag.Flag value.
func flagType(f *flag.Flag) string {
	typeName := strings.ToLower(fmt.Sprintf("%T", f.Value))
	switch {
	case strings.Contains(typeName, "bool"):
		return "bool"
	case strings.Contains(typeName, "float"):
		return "float64"
	case strings.Contains(typeName, "duration"):
		return "duration"
	case strings.Contains(typeName, "uint64"):
		return "uint64"
	case strings.Contains(typeName, "uint"):
		return "uint"
	case strings.Contains(typeName, "int64"):
		return "int64"
	case strings.Contains(typeName, "int"):
		return "int"
	default:
		return "string"
	}
}
