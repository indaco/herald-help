// A real kong CLI with a "greet" command and themed help via herald.
// Run:
//
//	cd examples/rich/004_kong && go run . greet --name Alice
//	cd examples/rich/004_kong && go run . --help
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
	heraldkong "github.com/indaco/herald-help/kong"
)

type CLI struct {
	Verbose bool     `help:"Enable verbose output" short:"v"`
	Greet   GreetCmd `cmd:"" help:"Greet someone by name" aliases:"g"`
}

type GreetCmd struct {
	Name string `help:"Name of the person to greet" short:"n" default:"World"`
}

func (g *GreetCmd) Run() error {
	fmt.Printf("Hello, %s!\n", g.Name)
	return nil
}

func main() {
	ty := herald.New()

	var cli CLI
	parser, err := kong.New(&cli,
		kong.Name("greeter"),
		kong.Description("A friendly CLI that greets people by name, powered by kong."),
		kong.Help(func(_ kong.HelpOptions, ctx *kong.Context) error {
			hc := heraldkong.FromKong(ctx.Kong)
			hc.Footer = heraldhelp.FormatVersion("greeter", "1.0.0")
			return heraldhelp.RenderTo(os.Stdout, ty, hc, heraldhelp.WithStyle(heraldhelp.StyleRich))
		}),
	)
	if err != nil {
		os.Stderr.WriteString("error: " + err.Error() + "\n")
		os.Exit(1)
	}

	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)
	err = ctx.Run()
	parser.FatalIfErrorf(err)
}
