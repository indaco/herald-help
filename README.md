<h1 align="center">
  herald-help
</h1>

<h2 align="center" style="font-size: 1.5rem;">
    Themed CLI help pages powered by herald typography.
</h2>

<p align="center">
  <a href="https://github.com/indaco/herald-help/actions/workflows/ci-root.yml" target="_blank">
    <img src="https://github.com/indaco/herald-help/actions/workflows/ci-root.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://codecov.io/gh/indaco/herald-help" target="_blank">
    <img src="https://codecov.io/gh/indaco/herald-help/branch/main/graph/badge.svg" alt="Code coverage" />
  </a>
  <a href="https://goreportcard.com/report/github.com/indaco/herald-help" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/indaco/herald-help" alt="Go Report Card" />
  </a>
  <a href="https://github.com/indaco/herald-help/actions/workflows/security.yml" target="_blank">
    <img src="https://github.com/indaco/herald-help/actions/workflows/security.yml/badge.svg" alt="Security Scan" />
  </a>
  <a href="https://github.com/indaco/herald-help/releases" target="_blank">
    <img src="https://img.shields.io/github/v/tag/indaco/herald-help?label=version&sort=semver&color=4c1" alt="version">
  </a>
  <a href="https://pkg.go.dev/github.com/indaco/herald-help" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/indaco/herald-help.svg" alt="Go Reference" />
  </a>
  <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square" alt="License" />
  </a>
</p>

<p align="center">
  <b><a href="#when-to-use-herald-help">When to use</a></b> |
  <b><a href="#installation">Installation</a></b> |
  <b><a href="#quick-start">Quick Start</a></b> |
  <b><a href="#styles">Styles</a></b> |
  <b><a href="#types">Types</a></b> |
  <b><a href="#rendering">Rendering</a></b> |
  <b><a href="#options">Options</a></b> |
  <b><a href="#examples">Examples</a></b> |
  <b><a href="#repository-structure">Structure</a></b>
</p>

herald-help generates themed CLI help pages through [herald](https://github.com/indaco/herald)'s typography system. Works with [cobra](https://github.com/spf13/cobra), [urfave/cli](https://github.com/urfave/cli), [kong](https://github.com/alecthomas/kong), and stdlib `flag`.

<p align="center">
  <img src="https://raw.githubusercontent.com/indaco/gh-assets/main/herald-help/demo-hero.png" alt="herald-help style variants" width="700" />
</p>

## When to use herald-help

Cobra, Kong, and urfave/cli all generate plain-text help output with basic formatting. herald-help takes a different approach. Paired with [herald](https://github.com/indaco/herald), it gives you themed help pages that match the rest of your CLI output. Use it when:

- You want **one theme for your entire CLI** - pick a built-in theme (Dracula, Catppuccin, Base16, Charm) or define your own, and help pages use it automatically.
- You need **structured, composable help** - `Render` returns a plain string you can pass to herald's `Compose` alongside other output.
- You want **framework-agnostic help rendering** - define your help structure once with `Command`, then use adapters to convert from any CLI framework.

## Installation

Requires Go 1.25+.

```sh
go get github.com/indaco/herald-help@latest
```

## Quick Start

Framework adapters are sub-modules — each has its own `go.mod`, so importing one does not pull in the others.

| Framework                                      | Install                                              | Package        |
| ---------------------------------------------- | ---------------------------------------------------- | -------------- |
| [cobra](https://github.com/spf13/cobra)        | `go get github.com/indaco/herald-help/cobra@latest`  | `heraldcobra`  |
| [urfave/cli/v3](https://github.com/urfave/cli) | `go get github.com/indaco/herald-help/urfave@latest` | `heraldurfave` |
| [kong](https://github.com/alecthomas/kong)     | `go get github.com/indaco/herald-help/kong@latest`   | `heraldkong`   |

The stdlib `flag` adapter (`FromFlagSet`) is included in the core module — no extra install needed.

### Standalone

See [`examples/compact/000_manual`](examples/compact/000_manual/) for a full runnable example. For stdlib `flag` integration, see [`examples/compact/001_flag`](examples/compact/001_flag/).

> [!NOTE]
> The Go module path is `github.com/indaco/herald-help` but the package name is `heraldhelp`. Use an import alias: `heraldhelp "github.com/indaco/herald-help"`.

```go
import (
    "fmt"
    "github.com/indaco/herald"
    heraldhelp "github.com/indaco/herald-help"
)

ty := herald.New()

cmd := heraldhelp.Command{
    Name:        "myapp",
    Synopsis:    "myapp [flags] <command>",
    Description: "A sample application.",
    Flags: []heraldhelp.Flag{
        {Long: "--output", Short: "-o", Type: "string", Default: "stdout", Desc: "Output destination"},
        {Long: "--verbose", Short: "-v", Type: "bool", Default: "false", Desc: "Enable verbose output"},
    },
    Commands: []heraldhelp.CommandRef{
        {Name: "serve", Desc: "Start the server"},
        {Name: "build", Desc: "Build the project"},
    },
    Footer: heraldhelp.FormatVersion("myapp", "1.0.0"),
}

fmt.Println(heraldhelp.Render(ty, cmd))
```

### With cobra

See [`examples/compact/002_cobra`](examples/compact/002_cobra/) for a full runnable example.

```go
import (
    "github.com/indaco/herald"
    heraldhelp "github.com/indaco/herald-help"
    heraldcobra "github.com/indaco/herald-help/cobra"
    "github.com/spf13/cobra"
)

rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
    ty := herald.New()
    heraldhelp.RenderTo(cmd.OutOrStdout(), ty, heraldcobra.FromCobra(cmd))
})
```

### With urfave/cli

See [`examples/compact/003_urfave`](examples/compact/003_urfave/) for a full runnable example.

```go
import (
    "io"
    heraldhelp "github.com/indaco/herald-help"
    heraldurfave "github.com/indaco/herald-help/urfave"
    "github.com/urfave/cli/v3"
)

cli.HelpPrinter = func(w io.Writer, _ string, data any) {
    cmd, ok := data.(*cli.Command)
    if !ok {
        return
    }
    ty := herald.New()
    heraldhelp.RenderTo(w, ty, heraldurfave.FromUrfave(cmd))
}
```

### With kong

See [`examples/compact/004_kong`](examples/compact/004_kong/) for a full runnable example.

```go
import (
    heraldhelp "github.com/indaco/herald-help"
    heraldkong "github.com/indaco/herald-help/kong"
)

parser, _ := kong.New(&cli,
    kong.Help(func(_ kong.HelpOptions, ctx *kong.Context) error {
        ty := herald.New()
        return heraldhelp.RenderTo(os.Stdout, ty, heraldkong.FromKong(ctx.Kong))
    }),
)
```

## Styles

herald-help supports four rendering styles, controlled by the `WithStyle` option.

| Style           | Description                                                                                                      |
| --------------- | ---------------------------------------------------------------------------------------------------------------- |
| `StyleCompact`  | **Default.** Uppercase colored headings, indented two-column lists, no borders. Closest to traditional CLI help. |
| `StyleRich`     | Decorated headings (H1/H2/H3), bordered tables, code blocks, and alert panels.                                   |
| `StyleGrouped`  | Each section wrapped in a bordered fieldset with the section name as legend.                                     |
| `StyleMarkdown` | Valid Markdown text. Pipe to `glow`/`bat`, save to file, or render via herald-md.                                |

<table align="center">
  <tr>
    <td align="center" valign="top"><img src="https://raw.githubusercontent.com/indaco/gh-assets/main/herald-help/demo-compact.png" alt="Compact style demo" width="400" /><br/><sub><code>StyleCompact</code> (default)</sub></td>
    <td align="center" valign="top"><img src="https://raw.githubusercontent.com/indaco/gh-assets/main/herald-help/demo-grouped.png" alt="Grouped style demo" width="400" /><br/><sub><code>StyleGrouped</code></sub></td>
  </tr>
  <tr>
    <td align="center" valign="top"><img src="https://raw.githubusercontent.com/indaco/gh-assets/main/herald-help/demo-rich.png" alt="Rich style demo" width="400" /><br/><sub><code>StyleRich</code></sub></td>
    <td align="center" valign="top"><img src="https://raw.githubusercontent.com/indaco/gh-assets/main/herald-help/demo-markdown.png" alt="Markdown style demo" width="400" /><br/><sub><code>StyleMarkdown</code></sub></td>
  </tr>
</table>

To use a non-default style, pass `WithStyle` to `Render` or `RenderTo`:

```go
// Rich style with bordered tables
heraldhelp.RenderTo(os.Stdout, ty, cmd, heraldhelp.WithStyle(heraldhelp.StyleRich))

// Grouped style with fieldset borders
heraldhelp.RenderTo(os.Stdout, ty, cmd, heraldhelp.WithStyle(heraldhelp.StyleGrouped))

// Markdown output — pipe to glow or render via herald-md
md := heraldhelp.Render(ty, cmd, heraldhelp.WithStyle(heraldhelp.StyleMarkdown))
```

## Types

| Type           | Description                                                                               |
| -------------- | ----------------------------------------------------------------------------------------- |
| `Command`      | Top-level struct: name, synopsis, description, flags, args, subcommands, examples, footer |
| `Flag`         | A CLI flag with long/short names, type, default, description, env vars, enum values       |
| `FlagGroup`    | Named group of flags (cobra groups, urfave categories, kong groups)                       |
| `Arg`          | Positional argument with name, description, required/default                              |
| `CommandRef`   | Subcommand summary: name, aliases, description                                            |
| `CommandGroup` | Named group of subcommands                                                                |
| `Example`      | Usage example with description and command string                                         |

## Rendering

Each `Command` field maps to a section in the output. The exact herald method used depends on the active style (see [Styles](#styles) above). The table below shows the `StyleRich` mapping as reference:

| Section           | Herald method (StyleRich)            | Condition                  |
| ----------------- | ------------------------------------ | -------------------------- |
| Command name      | `H1`                                 | Always                     |
| Deprecated notice | `Alert(AlertWarning)`                | If `Deprecated != ""`      |
| Synopsis          | `CodeBlock`                          | If `Synopsis != ""`        |
| Description       | `P`                                  | If `Description != ""`     |
| Positional args   | `H2("Arguments")` + `Table`          | If args present            |
| Flags             | `H2("Flags")` + `Table`              | If flags present           |
| Inherited flags   | `H3("Inherited Flags")` + `Table`    | If inherited flags present |
| Subcommands       | `H2("Commands")` + `Table`           | If commands present        |
| Examples          | `H2("Examples")` + `P` + `CodeBlock` | If examples present        |
| See Also          | `H2("See Also")` + `UL`              | If see-also present        |
| Footer            | `Small`                              | If `Footer != ""`          |

Flags are displayed in GNU-style format (`-o, --output`). Environment variables are shown inline (`[$PORT]`).

**Functions:**

- `Render(ty, cmd, opts...) string` - render to string
- `RenderTo(w, ty, cmd, opts...) error` - render to `io.Writer`
- `FormatVersion(name, version) string` - format a version footer string

## Options

| Option                          | Description                                                                        |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `WithStyle(style)`              | Set rendering style (`StyleCompact`, `StyleRich`, `StyleGrouped`, `StyleMarkdown`) |
| `WithWidth(n)`                  | Set output width (default: auto-detect terminal width)                             |
| `WithSectionOrder(sections...)` | Customize section ordering; omitted sections are hidden                            |
| `WithShowHidden(bool)`          | Show hidden flags (default: false)                                                 |
| `WithEnvVarDisplay(bool)`       | Show environment variable bindings (default: true)                                 |

## Examples

Runnable examples are in the [`examples/`](examples/) directory, organized by style. The `0xx` examples run directly; the adapter examples (`002+`) are separate modules.

**Compact style** (default — uppercase headings, indented two-column lists, no borders):

| Example                                                | Description                              | Run                                                 |
| ------------------------------------------------------ | ---------------------------------------- | --------------------------------------------------- |
| [000_manual](examples/compact/000_manual/)             | Manual `Command` construction            | `go run ./examples/compact/000_manual/`             |
| [001_flag](examples/compact/001_flag/)                 | stdlib `flag.FlagSet` adapter            | `go run ./examples/compact/001_flag/`               |
| [002_cobra](examples/compact/002_cobra/)               | Real cobra CLI with themed `--help`      | `cd examples/compact/002_cobra && go run . --help`  |
| [003_urfave](examples/compact/003_urfave/)             | Real urfave/cli app with themed `--help` | `cd examples/compact/003_urfave && go run . --help` |
| [004_kong](examples/compact/004_kong/)                 | Real kong CLI with themed `--help`       | `cd examples/compact/004_kong && go run . --help`   |
| [005_custom-theme](examples/compact/005_custom-theme/) | Dracula theme applied via herald         | `go run ./examples/compact/005_custom-theme/`       |

**Rich style** (decorated headings, bordered tables, code blocks):

| Example                                 | Description                              | Run                                              |
| --------------------------------------- | ---------------------------------------- | ------------------------------------------------ |
| [000_manual](examples/rich/000_manual/) | Manual `Command` construction            | `go run ./examples/rich/000_manual/`             |
| [001_flag](examples/rich/001_flag/)     | stdlib `flag.FlagSet` adapter            | `go run ./examples/rich/001_flag/`               |
| [002_cobra](examples/rich/002_cobra/)   | Real cobra CLI with themed `--help`      | `cd examples/rich/002_cobra && go run . --help`  |
| [003_urfave](examples/rich/003_urfave/) | Real urfave/cli app with themed `--help` | `cd examples/rich/003_urfave && go run . --help` |
| [004_kong](examples/rich/004_kong/)     | Real kong CLI with themed `--help`       | `cd examples/rich/004_kong && go run . --help`   |

**Grouped style** (each section in a bordered fieldset):

| Example                                    | Description                   | Run                                     |
| ------------------------------------------ | ----------------------------- | --------------------------------------- |
| [000_manual](examples/grouped/000_manual/) | Manual `Command` construction | `go run ./examples/grouped/000_manual/` |

**Markdown style** (valid Markdown text, pipeable to glow/bat):

| Example                                                           | Description                         | Run                                                      |
| ----------------------------------------------------------------- | ----------------------------------- | -------------------------------------------------------- |
| [000_manual](examples/markdown/000_manual/)                       | Markdown output                     | `go run ./examples/markdown/000_manual/ --help`          |
| [001_heraldmd-pipeline](examples/markdown/001_heraldmd-pipeline/) | Markdown rendered through herald-md | `cd examples/markdown/001_heraldmd-pipeline && go run .` |

## Repository structure

This is a multi-module Go repository. The core module lives at the root; framework adapters are sub-modules in their own directories, each with a separate `go.mod`:

```text
herald-help/
  go.mod            # github.com/indaco/herald-help (core types + renderer + flag adapter)
  cobra/go.mod      # github.com/indaco/herald-help/cobra
  urfave/go.mod     # github.com/indaco/herald-help/urfave
  kong/go.mod       # github.com/indaco/herald-help/kong
  examples/
    compact/        # StyleCompact examples (default; 002-004 are separate modules)
    rich/           # StyleRich examples (002-004 are separate modules)
    grouped/        # StyleGrouped examples
    markdown/       # StyleMarkdown examples (001 is a separate module with herald-md)
```

Use `just test-all`, `just lint-all`, `just check-all` to run across all modules.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
