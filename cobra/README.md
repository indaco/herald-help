# herald-help/cobra

<p>
  <a href="https://github.com/indaco/herald-help/actions/workflows/ci-cobra.yml" target="_blank">
    <img src="https://github.com/indaco/herald-help/actions/workflows/ci-cobra.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://pkg.go.dev/github.com/indaco/herald-help/cobra" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/indaco/herald-help/cobra.svg" alt="Go Reference" />
  </a>
  <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square" alt="License" />
  </a>
</p>

Cobra adapter for [herald-help](../) - converts a `cobra.Command` into a `heraldhelp.Command` for themed help rendering.

## Installation

```sh
go get github.com/indaco/herald-help/cobra@latest
```

## Quick Start

```go
import (
    "github.com/indaco/herald"
    heraldhelp "github.com/indaco/herald-help"
    heraldcobra "github.com/indaco/herald-help/cobra"
    "github.com/spf13/cobra"
)

rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
    ty := herald.New()
    _ = heraldhelp.RenderTo(cmd.OutOrStdout(), ty, heraldcobra.FromCobra(cmd))
})
// The default style is StyleCompact. Pass heraldhelp.WithStyle(heraldhelp.StyleRich),
// heraldhelp.WithStyle(heraldhelp.StyleGrouped), or heraldhelp.WithStyle(heraldhelp.StyleMarkdown)
// to RenderTo to use a different style.
```

## What it extracts

| Cobra field      | heraldhelp field                 | Notes                                            |
| ---------------- | -------------------------------- | ------------------------------------------------ |
| `Name()`         | `Name`                           |                                                  |
| `UseLine()`      | `Synopsis`                       |                                                  |
| `Long` / `Short` | `Description`                    | Prefers `Long`, falls back to `Short`            |
| `Aliases`        | `Aliases`                        |                                                  |
| `Deprecated`     | `Deprecated`                     |                                                  |
| Local flags      | `Flags`                          | Excludes hidden flags by default                 |
| Inherited flags  | `Flags` (with `Inherited: true`) | Rendered in a separate "Inherited Flags" section |
| Required flags   | `Flag.Required`                  | Via `cobra.MarkFlagRequired`                     |
| `Commands()`     | `Commands`                       | Excludes hidden subcommands                      |
| `Example`        | `Examples`                       | Parsed into description + command pairs          |

See the [main README](../) for full rendering options and configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
