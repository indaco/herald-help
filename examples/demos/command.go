// Package demos provides a shared Command for demo screenshots.
package demos

import heraldhelp "github.com/indaco/herald-help"

// DemoCommand returns a representative Command for style demos.
func DemoCommand() heraldhelp.Command {
	return heraldhelp.Command{
		Name:        "myapp",
		Synopsis:    "myapp [flags] <command>",
		Description: "A sample CLI application for demonstrating herald-help styles.",
		Flags: []heraldhelp.Flag{
			{Long: "--output", Short: "-o", Type: "string", Default: "stdout", Desc: "Output destination"},
			{Long: "--verbose", Short: "-v", Type: "bool", Desc: "Enable verbose output"},
			{Long: "--port", Type: "int", Default: "8080", Desc: "Port number", EnvVars: []string{"PORT"}},
			{Long: "--format", Type: "string", Desc: "Output format", Enum: []string{"json", "yaml", "toml"}},
		},
		Commands: []heraldhelp.CommandRef{
			{Name: "serve", Aliases: []string{"s"}, Desc: "Start the server"},
			{Name: "build", Desc: "Build the project"},
			{Name: "config", Aliases: []string{"c"}, Desc: "Manage configuration"},
		},
		Examples: []heraldhelp.Example{
			{Desc: "Start on port 9090:", Command: "myapp serve --port 9090"},
			{Desc: "Build with verbose:", Command: "myapp build -v"},
		},
		Footer: heraldhelp.FormatVersion("myapp", "1.0.0"),
	}
}
