package heraldcobra

import "testing"

func FuzzParseExamples(f *testing.F) {
	f.Add("$ go run main.go\nsome description\n$ go test ./...")
	f.Add("Just a description with no command")
	f.Add("")
	f.Add("First line\nSecond line\n$ myapp run")
	f.Add("$\tmyapp run")

	f.Fuzz(func(t *testing.T, input string) {
		// parseExamples must not panic on any input.
		_ = parseExamples(input)
	})
}
