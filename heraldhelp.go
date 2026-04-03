// Package heraldhelp renders CLI help/usage output using herald typography.
//
// It provides a Command model that captures the full structure of a CLI command
// (flags, arguments, subcommands, examples) and renders it with herald's
// styled typography methods.
//
// Quick start (manual):
//
//	ty := herald.New()
//	cmd := heraldhelp.Command{
//	    Name:     "myapp",
//	    Synopsis: "myapp [flags] <subcommand>",
//	}
//	fmt.Println(heraldhelp.Render(ty, cmd))
//
// With the flag adapter:
//
//	cmd := heraldhelp.FromFlagSet("myapp", flag.CommandLine)
//	fmt.Println(heraldhelp.Render(herald.New(), cmd))
package heraldhelp

// Command describes a CLI command for help rendering. It captures all the
// information typically shown in a --help page: name, synopsis, flags,
// subcommands, examples, etc.
type Command struct {
	Name        string // command name, e.g. "myapp"
	Synopsis    string // one-line usage, e.g. "myapp [flags] <subcommand>"
	Description string // long description (may be multi-paragraph)
	Aliases     []string
	Deprecated  string // if non-empty, command is deprecated

	Flags      []Flag
	FlagGroups []FlagGroup // cobra flag groups, kong groups, urfave categories

	Args []Arg // positional arguments

	Commands      []CommandRef   // subcommands (summary only, not recursive)
	CommandGroups []CommandGroup // grouped subcommands

	Examples []Example
	SeeAlso  []string
	Footer   string // version string, bug URL, etc.
}

// Flag describes a single CLI flag.
type Flag struct {
	Long       string   // e.g. "--output"
	Short      string   // e.g. "-o" (empty if none)
	Type       string   // e.g. "string", "bool"
	Default    string   // default value
	Desc       string   // description
	Required   bool     // whether the flag is required
	Hidden     bool     // adapters skip hidden flags by default
	EnvVars    []string // environment variable bindings
	Enum       []string // enum values (kong)
	Deprecated string   // deprecation notice
	Inherited  bool     // inherited/persistent flag (rendered in separate section)
}

// FlagGroup groups related flags under a heading.
type FlagGroup struct {
	Name  string
	Flags []Flag
}

// Arg describes a positional argument.
type Arg struct {
	Name     string // e.g. "<file>"
	Desc     string
	Required bool
	Default  string
}

// CommandRef is a summary reference to a subcommand. It does not recurse
// into the subcommand's own flags or children.
type CommandRef struct {
	Name    string
	Aliases []string
	Desc    string
}

// CommandGroup groups related subcommands under a heading.
type CommandGroup struct {
	Name     string
	Commands []CommandRef
}

// Example pairs a description with a command invocation.
type Example struct {
	Desc    string
	Command string
}
