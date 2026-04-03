package heraldhelp

// RenderConfig holds configuration for help rendering. Use RenderOption
// functional options to customise the defaults.
type RenderConfig struct {
	Style         Style     // visual style (StyleRich or StyleCompact)
	Width         int       // explicit terminal width override (0 = auto-detect)
	SectionOrder  []Section // section rendering order (nil = DefaultSectionOrder)
	ShowHidden    bool      // include hidden flags in output
	EnvVarDisplay bool      // show environment variable bindings inline
}

// RenderOption is a functional option for Render/RenderTo.
type RenderOption func(*RenderConfig)

// WithStyle sets the visual style for help rendering. The default is
// StyleRich (decorated headings, bordered tables). Use StyleCompact for a
// minimal, terminal-native layout with colored text and indented columns.
func WithStyle(s Style) RenderOption {
	return func(cfg *RenderConfig) {
		cfg.Style = s
	}
}

// WithWidth sets an explicit terminal width, overriding auto-detection.
func WithWidth(w int) RenderOption {
	return func(cfg *RenderConfig) {
		cfg.Width = w
	}
}

// WithSectionOrder sets the order of help sections. Sections not listed are
// omitted from the output.
func WithSectionOrder(order ...Section) RenderOption {
	return func(cfg *RenderConfig) {
		cfg.SectionOrder = order
	}
}

// WithShowHidden includes hidden flags in the help output.
func WithShowHidden(show bool) RenderOption {
	return func(cfg *RenderConfig) {
		cfg.ShowHidden = show
	}
}

// WithEnvVarDisplay enables inline display of environment variable bindings
// in flag descriptions (e.g. "Port number [$PORT]").
func WithEnvVarDisplay(show bool) RenderOption {
	return func(cfg *RenderConfig) {
		cfg.EnvVarDisplay = show
	}
}
