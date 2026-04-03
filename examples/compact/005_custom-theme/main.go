// Themed help using a built-in herald theme (Dracula).
// The theme is set on the herald.Typography instance - herald-help inherits
// all colors automatically.
// Run:
//
//	go run ./examples/compact/005_custom-theme/ --help
package main

import (
	"fmt"
	"os"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
)

func main() {
	// Use Dracula theme instead of the default Rose Pine.
	ty := herald.New(herald.WithTheme(herald.DraculaTheme()))

	cmd := heraldhelp.Command{
		Name:        "myapp",
		Synopsis:    "myapp [flags] <command>",
		Description: "A sample application with Dracula-themed help output.",
		Flags: []heraldhelp.Flag{
			{Long: "--output", Short: "-o", Type: "string", Default: "stdout", Desc: "Output destination"},
			{Long: "--verbose", Short: "-v", Type: "bool", Desc: "Enable verbose output"},
			{Long: "--port", Type: "int", Default: "8080", Desc: "Port number", EnvVars: []string{"PORT"}},
		},
		Commands: []heraldhelp.CommandRef{
			{Name: "serve", Aliases: []string{"s"}, Desc: "Start the server"},
			{Name: "build", Desc: "Build the project"},
		},
		Examples: []heraldhelp.Example{
			{Desc: "Start on port 9090:", Command: "myapp serve --port 9090"},
			{Desc: "Build with verbose:", Command: "myapp build -v"},
		},
		Footer: heraldhelp.FormatVersion("myapp", "1.0.0"),
	}

	args := os.Args[1:]
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(heraldhelp.Render(ty, cmd))
		return
	}

	os.Stderr.WriteString("unknown command: " + args[0] + "\n")
	os.Exit(1)
}
