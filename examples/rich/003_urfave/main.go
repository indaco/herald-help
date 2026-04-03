// A real urfave/cli app with a "greet" command and themed help via herald.
// Run:
//
//	cd examples/rich/003_urfave && go run . greet --name Alice
//	cd examples/rich/003_urfave && go run . --help
//	cd examples/rich/003_urfave && go run . greet --help
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
	heraldurfave "github.com/indaco/herald-help/urfave"
	"github.com/urfave/cli/v3"
)

func main() {
	ty := herald.New()

	// Override the global help printer with herald-themed output.
	cli.HelpPrinter = func(w io.Writer, _ string, data any) {
		cmd, ok := data.(*cli.Command)
		if !ok {
			return
		}
		hc := heraldurfave.FromUrfave(cmd)
		hc.Footer = heraldhelp.FormatVersion("greeter", "1.0.0")
		_ = heraldhelp.RenderTo(w, ty, hc, heraldhelp.WithStyle(heraldhelp.StyleRich))
	}

	app := &cli.Command{
		Name:        "greeter",
		Usage:       "A friendly CLI that greets people",
		Description: "A friendly CLI that greets people by name, powered by urfave/cli.",
		UsageText:   "greeter [flags] <command>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "greet",
				Usage:   "Greet someone by name",
				Aliases: []string{"g"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "World",
						Usage:   "Name of the person to greet",
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					fmt.Printf("Hello, %s!\n", cmd.String("name"))
					return nil
				},
			},
		},
	}

	_ = app.Run(context.Background(), os.Args)
}
