// StyleGrouped demo for screenshot capture.
// Run: go run ./examples/demos/grouped/
package main

import (
	"fmt"

	"github.com/indaco/herald"
	heraldhelp "github.com/indaco/herald-help"
	"github.com/indaco/herald-help/examples/demos"
)

func main() {
	ty := herald.New()
	fmt.Println(heraldhelp.Render(ty, demos.DemoCommand(), heraldhelp.WithStyle(heraldhelp.StyleGrouped)))
}
