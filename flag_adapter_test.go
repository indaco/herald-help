package heraldhelp

import (
	"flag"
	"testing"
	"time"
)

func TestFromFlagSet(t *testing.T) {
	fs := flag.NewFlagSet("testapp", flag.ContinueOnError)
	fs.String("output", "stdout", "Output destination")
	fs.Bool("verbose", false, "Enable verbose mode")
	fs.Int("port", 8080, "Server port")

	cmd := FromFlagSet("testapp", fs)

	if cmd.Name != "testapp" {
		t.Errorf("Name = %q, want %q", cmd.Name, "testapp")
	}
	if cmd.Synopsis != "testapp [flags]" {
		t.Errorf("Synopsis = %q, want %q", cmd.Synopsis, "testapp [flags]")
	}
	if len(cmd.Flags) != 3 {
		t.Fatalf("Flags count = %d, want 3", len(cmd.Flags))
	}

	// flag.FlagSet.VisitAll visits in lexical order.
	flagMap := make(map[string]Flag)
	for _, f := range cmd.Flags {
		flagMap[f.Long] = f
	}

	t.Run("string flag", func(t *testing.T) {
		f := flagMap["--output"]
		if f.Type != "string" {
			t.Errorf("output type = %q, want %q", f.Type, "string")
		}
		if f.Default != "stdout" {
			t.Errorf("output default = %q, want %q", f.Default, "stdout")
		}
		if f.Desc != "Output destination" {
			t.Errorf("output desc = %q, want %q", f.Desc, "Output destination")
		}
	})

	t.Run("bool flag", func(t *testing.T) {
		f := flagMap["--verbose"]
		if f.Type != "bool" {
			t.Errorf("verbose type = %q, want %q", f.Type, "bool")
		}
		if f.Default != "false" {
			t.Errorf("verbose default = %q, want %q", f.Default, "false")
		}
	})

	t.Run("int flag", func(t *testing.T) {
		f := flagMap["--port"]
		if f.Type != "int" {
			t.Errorf("port type = %q, want %q", f.Type, "int")
		}
		if f.Default != "8080" {
			t.Errorf("port default = %q, want %q", f.Default, "8080")
		}
	})
}

func TestFromFlagSetEmpty(t *testing.T) {
	fs := flag.NewFlagSet("empty", flag.ContinueOnError)
	cmd := FromFlagSet("empty", fs)

	if cmd.Name != "empty" {
		t.Errorf("Name = %q, want %q", cmd.Name, "empty")
	}
	if len(cmd.Flags) != 0 {
		t.Errorf("Flags count = %d, want 0", len(cmd.Flags))
	}
}

func TestFromFlagSetAllTypes(t *testing.T) {
	fs := flag.NewFlagSet("types", flag.ContinueOnError)
	fs.Bool("bool", false, "a bool")
	fs.Int("int", 0, "an int")
	fs.Int64("int64", 0, "an int64")
	fs.Uint("uint", 0, "a uint")
	fs.Uint64("uint64", 0, "a uint64")
	fs.Float64("float64", 0, "a float64")
	fs.String("string", "", "a string")
	fs.Duration("duration", 0, "a duration")

	cmd := FromFlagSet("types", fs)

	expected := map[string]string{
		"--bool":     "bool",
		"--int":      "int",
		"--int64":    "int64",
		"--uint":     "uint",
		"--uint64":   "uint64",
		"--float64":  "float64",
		"--string":   "string",
		"--duration": "duration",
	}

	flagMap := make(map[string]Flag)
	for _, f := range cmd.Flags {
		flagMap[f.Long] = f
	}

	for long, wantType := range expected {
		t.Run(long, func(t *testing.T) {
			f, ok := flagMap[long]
			if !ok {
				t.Fatalf("flag %s not found", long)
			}
			if f.Type != wantType {
				t.Errorf("%s type = %q, want %q", long, f.Type, wantType)
			}
		})
	}
}

func TestFlagTypeWithNonZeroDefaults(t *testing.T) {
	fs := flag.NewFlagSet("nonzero", flag.ContinueOnError)
	fs.Int("count", 42, "a non-zero int")
	fs.Float64("rate", 1.5, "a non-zero float")
	fs.Duration("timeout", 5*time.Second, "a non-zero duration")
	fs.Bool("active", true, "a true bool")

	cmd := FromFlagSet("nonzero", fs)

	expected := map[string]string{
		"--count":   "int",
		"--rate":    "float64",
		"--timeout": "duration",
		"--active":  "bool",
	}

	flagMap := make(map[string]Flag)
	for _, f := range cmd.Flags {
		flagMap[f.Long] = f
	}

	for long, wantType := range expected {
		t.Run(long, func(t *testing.T) {
			f, ok := flagMap[long]
			if !ok {
				t.Fatalf("flag %s not found", long)
			}
			if f.Type != wantType {
				t.Errorf("%s type = %q, want %q", long, f.Type, wantType)
			}
		})
	}
}

func TestFlagTypeDetection(t *testing.T) {
	// Test flagType directly via FromFlagSet with various combinations.
	fs := flag.NewFlagSet("detect", flag.ContinueOnError)
	fs.String("name", "default-val", "a string with non-empty default")
	fs.Uint64("big", 0, "a uint64 with zero default")

	cmd := FromFlagSet("detect", fs)

	flagMap := make(map[string]Flag)
	for _, f := range cmd.Flags {
		flagMap[f.Long] = f
	}

	if f := flagMap["--name"]; f.Type != "string" {
		t.Errorf("name type = %q, want string", f.Type)
	}
	if f := flagMap["--big"]; f.Type != "uint64" {
		t.Errorf("big type = %q, want uint64", f.Type)
	}
}

func BenchmarkFromFlagSet(b *testing.B) {
	fs := flag.NewFlagSet("bench", flag.ContinueOnError)
	fs.String("output", "stdout", "Output destination")
	fs.Bool("verbose", false, "Enable verbose mode")
	fs.Int("port", 8080, "Server port")
	fs.Float64("rate", 1.0, "Rate limit")
	fs.Duration("timeout", 0, "Request timeout")

	for b.Loop() {
		_ = FromFlagSet("bench", fs)
	}
}
