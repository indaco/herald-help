# herald-help/urfave

<p>
  <a href="https://github.com/indaco/herald-help/actions/workflows/ci-urfave.yml" target="_blank">
    <img src="https://github.com/indaco/herald-help/actions/workflows/ci-urfave.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://pkg.go.dev/github.com/indaco/herald-help/urfave" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/indaco/herald-help/urfave.svg" alt="Go Reference" />
  </a>
  <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square" alt="License" />
  </a>
</p>

urfave/cli adapter for [herald-help](../README.md) - converts a `cli.Command` into a `heraldhelp.Command` for themed help rendering.

## Installation

```sh
go get github.com/indaco/herald-help/urfave@latest
```

## Quick Start

```go
import (
    "io"
    "github.com/indaco/herald"
    heraldhelp "github.com/indaco/herald-help"
    heraldurfave "github.com/indaco/herald-help/urfave"
    "github.com/urfave/cli/v3"
)

// Override the global help printer with herald-themed output.
cli.HelpPrinter = func(w io.Writer, _ string, data any) {
    cmd, ok := data.(*cli.Command)
    if !ok {
        return
    }
    ty := herald.New()
    heraldhelp.RenderTo(w, ty, heraldurfave.FromUrfave(cmd))
}
// The default style is StyleCompact. Pass heraldhelp.WithStyle(heraldhelp.StyleRich),
// heraldhelp.WithStyle(heraldhelp.StyleGrouped), or heraldhelp.WithStyle(heraldhelp.StyleMarkdown)
// to RenderTo to use a different style.
```

## What it extracts

| urfave/cli field        | heraldhelp field             | Notes                                            |
| ----------------------- | ---------------------------- | ------------------------------------------------ |
| `Name`                  | `Name`                       |                                                  |
| `UsageText`             | `Synopsis`                   |                                                  |
| `Description` / `Usage` | `Description`                | Prefers `Description`, falls back to `Usage`     |
| Flags                   | `Flags` / `FlagGroups`       | Grouped by category when categories are set      |
| Flag `Sources` (env)    | `Flag.EnvVars`               | Environment variable bindings shown inline       |
| Hidden flags            | `Flag.Hidden`                | Excluded by default, shown with `WithShowHidden` |
| `Commands`              | `Commands` / `CommandGroups` | Grouped by category when categories are set      |
| Hidden commands         | Excluded                     |                                                  |

See the [main README](../README.md) for full rendering options and configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
