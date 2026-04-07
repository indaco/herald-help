package heraldhelp

import "slices"

// RenderConfig holds configuration for help rendering. Use RenderOption
// functional options to customise the defaults.
type RenderConfig struct {
	Style         Style     // visual style (StyleRich or StyleCompact)
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

// WithoutSections excludes the listed sections from the default section order.
// This is a convenience alternative to WithSectionOrder when you only want to
// hide a few sections rather than listing all the ones you want.
func WithoutSections(exclude ...Section) RenderOption {
	return func(cfg *RenderConfig) {
		wanted := DefaultSectionOrder()
		filtered := make([]Section, 0, len(wanted))
		for _, s := range wanted {
			if !slices.Contains(exclude, s) {
				filtered = append(filtered, s)
			}
		}
		cfg.SectionOrder = filtered
	}
}
