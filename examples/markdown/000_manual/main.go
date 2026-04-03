// Render help as valid Markdown text. Pipe to glow, bat, or save to file.
// Run:
//
//	go run ./examples/markdown/000_manual/ --help
//	go run ./examples/markdown/000_manual/ --help | glow
//	go run ./examples/markdown/000_manual/ --help | bat -l md
package main

import (
	"fmt"
	"os"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
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
			{Desc: "Shout:", Command: "greeter greet --name Bob --loud"},
		},
		Footer: heraldhelp.FormatVersion("greeter", "1.0.0"),
	}

	args := os.Args[1:]
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(heraldhelp.Render(ty, cmd, heraldhelp.WithStyle(heraldhelp.StyleMarkdown)))
		return
	}

	fmt.Println("Hello, World!")
}
