// A real cobra CLI with a "greet" command and compact-style themed help via herald.
// Run:
//
//	cd examples/compact/002_cobra && go run . greet --name Alice
//	cd examples/compact/002_cobra && go run . --help
//	cd examples/compact/002_cobra && go run . greet --help
package main

import (
	"fmt"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
	heraldcobra "github.com/indaco/herald-help/cobra"
	"github.com/spf13/cobra"
)

func main() {
	ty := herald.New()

	rootCmd := &cobra.Command{
		Use:   "greeter [flags] <command>",
		Short: "A friendly CLI that greets people",
		Long:  "A friendly CLI that greets people by name, powered by cobra.",
	}

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	greetCmd := &cobra.Command{
		Use:     "greet",
		Short:   "Greet someone by name",
		Aliases: []string{"g"},
		Example: "  greeter greet --name Alice",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Hello, %s!\n", name)
		},
	}
	greetCmd.Flags().StringP("name", "n", "World", "Name of the person to greet")

	rootCmd.AddCommand(greetCmd)

	// Replace the default help function with herald compact-style help.
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		hc := heraldcobra.FromCobra(cmd)
		hc.Footer = heraldhelp.FormatVersion("greeter", "1.0.0")
		_ = heraldhelp.RenderTo(cmd.OutOrStdout(), ty, hc)
	})

	_ = rootCmd.Execute()
}
