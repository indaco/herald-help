# herald-help/kong

<p>
  <a href="https://github.com/indaco/herald-help/actions/workflows/ci-kong.yml" target="_blank">
    <img src="https://github.com/indaco/herald-help/actions/workflows/ci-kong.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://pkg.go.dev/github.com/indaco/herald-help/kong" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/indaco/herald-help/kong.svg" alt="Go Reference" />
  </a>
  <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square" alt="License" />
  </a>
</p>

Kong adapter for [herald-help](../) - converts a `kong.Kong` application into a `heraldhelp.Command` for themed help rendering.

## Installation

```sh
go get github.com/indaco/herald-help/kong@latest
```

## Quick Start

```go
import (
    "os"
    "github.com/alecthomas/kong"
    "github.com/indaco/herald"
    heraldhelp "github.com/indaco/herald-help"
    heraldkong "github.com/indaco/herald-help/kong"
)

parser, _ := kong.New(&cli,
    kong.Help(func(_ kong.HelpOptions, ctx *kong.Context) error {
        ty := herald.New()
        return heraldhelp.RenderTo(os.Stdout, ty, heraldkong.FromKong(ctx.Kong))
    }),
)
// The default style is StyleCompact. Pass heraldhelp.WithStyle(heraldhelp.StyleRich),
// heraldhelp.WithStyle(heraldhelp.StyleGrouped), or heraldhelp.WithStyle(heraldhelp.StyleMarkdown)
// to RenderTo to use a different style.
```

## Functions

- `FromKong(app *kong.Kong) heraldhelp.Command` - converts the root application
- `FromNode(node *kong.Node) heraldhelp.Command` - converts a specific subcommand node

## What it extracts

| Kong field            | heraldhelp field             | Notes                                            |
| --------------------- | ---------------------------- | ------------------------------------------------ |
| `Node.Name`           | `Name`                       |                                                  |
| `Node.Help`           | `Description`                |                                                  |
| Flags                 | `Flags` / `FlagGroups`       | Grouped by `group` tag when set                  |
| Flag `env` tag        | `Flag.EnvVars`               | Environment variable bindings shown inline       |
| Flag `enum` tag       | `Flag.Enum`                  | Enum values shown in description                 |
| Hidden flags/commands | `Hidden`                     | Excluded by default, shown with `WithShowHidden` |
| Positional arguments  | `Args`                       | Name, help, required status, default value       |
| Subcommands           | `Commands` / `CommandGroups` | Grouped by `group` tag when set                  |
| `Node.Aliases`        | `CommandRef.Aliases`         |                                                  |

See the [main README](../) for full rendering options and configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
