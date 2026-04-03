// A real CLI with a "greet" command and compact-style themed help via herald.
// Run:
//
//	go run ./examples/compact/000_manual/ greet --name Alice
//	go run ./examples/compact/000_manual/ --help
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
			{Long: "--help", Short: "-h", Type: "bool", Desc: "Show this help"},
		},
		Commands: []heraldhelp.CommandRef{
			{Name: "greet", Desc: "Greet someone by name"},
		},
		Examples: []heraldhelp.Example{
			{Desc: "Greet Alice:", Command: "greeter greet --name Alice"},
			{Desc: "Show help:", Command: "greeter --help"},
		},
		Footer: heraldhelp.FormatVersion("greeter", "1.0.0"),
	}

	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(heraldhelp.Render(ty, cmd))
		return
	}

	if args[0] == "greet" {
		name := "World"
		for i, a := range args {
			if a == "--name" && i+1 < len(args) {
				name = args[i+1]
			}
		}
		fmt.Printf("Hello, %s!\n", name)
		return
	}

	os.Stderr.WriteString("unknown command: " + args[0] + "\n")
	os.Exit(1)
}
