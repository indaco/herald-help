package heraldhelp

// Style controls the visual style of the rendered help output.
type Style int

const (
	// StyleCompact renders a minimal, terminal-native layout: uppercase
	// colored section headings, indented two-column lists for flags and
	// commands, and no table borders. Closer to traditional CLI help output
	// but with themed colors applied. This is the default style.
	StyleCompact Style = iota

	// StyleRich uses herald's full typography: bordered tables, decorated
	// headings (H1/H2/H3), code blocks, and alert panels.
	StyleRich

	// StyleGrouped wraps each section in a herald Fieldset with the section
	// name as the legend. Content inside uses compact-style KV layout.
	StyleGrouped

	// StyleMarkdown renders help as valid Markdown text. The output can be
	// piped to tools like glow or bat, saved to a file, or rendered back
	// through herald-md for themed terminal output.
	StyleMarkdown
)
