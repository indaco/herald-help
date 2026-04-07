package heraldhelp

// Section identifies a renderable section of the help output.
type Section int

const (
	// SectionName renders the command name as H1.
	SectionName Section = iota
	// SectionDeprecated renders a deprecation warning alert.
	SectionDeprecated
	// SectionSynopsis renders the usage synopsis as a code block.
	SectionSynopsis
	// SectionDescription renders the long description as a paragraph.
	SectionDescription
	// SectionArgs renders positional arguments as a table.
	SectionArgs
	// SectionFlags renders flags (flat or grouped) as tables.
	SectionFlags
	// SectionInheritedFlags renders inherited/persistent flags as a table.
	SectionInheritedFlags
	// SectionCommands renders subcommands (flat or grouped) as tables.
	SectionCommands
	// SectionExamples renders usage examples with descriptions and code blocks.
	SectionExamples
	// SectionSeeAlso renders related commands/resources as a list.
	SectionSeeAlso
	// SectionFooter renders footer text (version, bug URL, etc.).
	SectionFooter
)

// DefaultSectionOrder returns the default ordering of help sections.
// Each call returns a fresh copy, so callers cannot corrupt the shared default.
func DefaultSectionOrder() []Section {
	return []Section{
		SectionName,
		SectionDeprecated,
		SectionSynopsis,
		SectionDescription,
		SectionArgs,
		SectionFlags,
		SectionInheritedFlags,
		SectionCommands,
		SectionExamples,
		SectionSeeAlso,
		SectionFooter,
	}
}
