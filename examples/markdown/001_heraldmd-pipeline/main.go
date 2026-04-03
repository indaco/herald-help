// Render help as Markdown, then pass through herald-md for themed terminal output.
// This demonstrates the herald-help → Markdown → herald-md pipeline.
// Run:
//
//	cd examples/markdown/001_heraldmd-pipeline && go run .
package main

import (
	"fmt"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
	heraldmd "github.com/indaco/herald-md"
)

func main() {
	ty := herald.New()

	cmd := heraldhelp.Command{
		Name:        "greeter",
		Synopsis:    "greeter [flags] <command>",
		Description: "A friendly CLI that greets people by name.",
		Flags: []heraldhelp.Flag{
			{Long: "--name", Short: "-n", Type: "string", Default: "World", Desc: "Name to greet"},
			{Long: "--loud", Type: "bool", Desc: "Shout the greeting"},
		},
		Commands: []heraldhelp.CommandRef{
			{Name: "greet", Aliases: []string{"g"}, Desc: "Greet someone by name"},
			{Name: "farewell", Desc: "Say goodbye"},
		},
		Examples: []heraldhelp.Example{
			{Desc: "Greet Alice:", Command: "greeter greet --name Alice"},
		},
		Footer: heraldhelp.FormatVersion("greeter", "1.0.0"),
	}

	// Step 1: Render to Markdown
	md := heraldhelp.Render(ty, cmd, heraldhelp.WithStyle(heraldhelp.StyleMarkdown))

	// Step 2: Render the Markdown through herald-md for themed terminal output
	fmt.Println(heraldmd.Render(ty, []byte(md)))
}
