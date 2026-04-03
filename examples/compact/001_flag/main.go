// A real CLI using stdlib flag with compact-style themed help via herald.
// Run:
//
//	go run ./examples/compact/001_flag/ -name Alice
//	go run ./examples/compact/001_flag/ -help
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
)

func main() {
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	name := fs.String("name", "World", "Name of the person to greet")
	showHelp := fs.Bool("help", false, "Show themed help")

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if *showHelp {
		ty := herald.New()
		cmd := heraldhelp.FromFlagSet("greeter", fs)
		cmd.Description = "A friendly CLI that greets people by name."
		cmd.Examples = []heraldhelp.Example{
			{Desc: "Greet Alice:", Command: "greeter -name Alice"},
		}
		cmd.Footer = heraldhelp.FormatVersion("greeter", "1.0.0")

		fmt.Println(heraldhelp.Render(ty, cmd))
		return
	}

	fmt.Printf("Hello, %s!\n", *name)
}
