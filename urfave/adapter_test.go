package heraldurfave

import (
	"testing"
	"time"

	heraldhelp "github.com/indaco/herald-help"
	"github.com/urfave/cli/v3"
)

// findFlag searches flat flags for the given long name, failing the test if
// not found. It also searches flag groups.
func findFlag(t *testing.T, flags []heraldhelp.Flag, long string) heraldhelp.Flag {
	t.Helper()
	for i := range flags {
		if flags[i].Long == long {
			return flags[i]
		}
	}
	t.Fatalf("flag %q not found", long)
	return heraldhelp.Flag{}
}

func newTestCommand() *cli.Command {
	return &cli.Command{
		Name:        "myapp",
		Usage:       "A sample application",
		Description: "A longer description of the sample application.",
		UsageText:   "myapp [flags] <command>",
		Aliases:     []string{"ma"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "stdout",
				Usage:   "Output destination",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
			},
			&cli.IntFlag{
				Name:     "port",
				Value:    8080,
				Usage:    "Port number",
				Sources:  cli.EnvVars("PORT"),
				Required: true,
			},
			&cli.FloatFlag{
				Name:  "rate",
				Value: 1.5,
				Usage: "Rate limit",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Value: 30 * time.Second,
				Usage: "Request timeout",
			},
			&cli.StringSliceFlag{
				Name:  "tags",
				Usage: "Tags to apply",
			},
			&cli.StringFlag{
				Name:     "host",
				Usage:    "Host address",
				Category: "Server",
			},
			&cli.BoolFlag{
				Name:     "tls",
				Usage:    "Enable TLS",
				Category: "Server",
			},
			&cli.IntFlag{
				Name:     "workers",
				Usage:    "Worker count",
				Category: "Server",
			},
			&cli.FloatFlag{
				Name:     "threshold",
				Usage:    "Error threshold",
				Category: "Server",
			},
			&cli.DurationFlag{
				Name:     "grace",
				Usage:    "Grace period",
				Category: "Server",
			},
			&cli.StringSliceFlag{
				Name:     "origins",
				Usage:    "Allowed origins",
				Category: "Server",
			},
			&cli.BoolFlag{
				Name:   "debug",
				Usage:  "Debug mode",
				Hidden: true,
			},
			&cli.BoolFlag{
				Name:  "active",
				Value: true,
				Usage: "Active mode",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Usage:   "Start the server",
				Aliases: []string{"s"},
			},
			{
				Name:     "build",
				Usage:    "Build the project",
				Category: "Development",
			},
			{
				Name:   "internal",
				Usage:  "Internal command",
				Hidden: true,
			},
		},
	}
}

func TestFromUrfave(t *testing.T) {
	cmd := newTestCommand()
	hc := FromUrfave(cmd)

	if hc.Name != "myapp" {
		t.Errorf("Name = %q, want %q", hc.Name, "myapp")
	}
	if hc.Synopsis != "myapp [flags] <command>" {
		t.Errorf("Synopsis = %q", hc.Synopsis)
	}
	if hc.Description != "A longer description of the sample application." {
		t.Errorf("Description = %q", hc.Description)
	}
	if len(hc.Aliases) != 1 || hc.Aliases[0] != "ma" {
		t.Errorf("Aliases = %v", hc.Aliases)
	}
}

func TestFromUrfaveStringFlag(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--output")

	if f.Short != "-o" {
		t.Errorf("output short = %q", f.Short)
	}
	if f.Type != "string" {
		t.Errorf("output type = %q", f.Type)
	}
	if f.Default != "stdout" {
		t.Errorf("output default = %q", f.Default)
	}
}

func TestFromUrfaveBoolFlag(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--verbose")

	if f.Type != "bool" {
		t.Errorf("verbose type = %q", f.Type)
	}
	if f.Default != "false" {
		t.Errorf("verbose default = %q", f.Default)
	}
}

func TestFromUrfaveBoolFlagTrueDefault(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--active")

	if f.Type != "bool" {
		t.Errorf("active type = %q", f.Type)
	}
	if f.Default != "true" {
		t.Errorf("active default = %q", f.Default)
	}
}

func TestFromUrfaveIntFlagWithEnv(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--port")

	if f.Type != "int" {
		t.Errorf("port type = %q", f.Type)
	}
	if !f.Required {
		t.Error("port should be required")
	}
	if len(f.EnvVars) != 1 || f.EnvVars[0] != "PORT" {
		t.Errorf("port envvars = %v", f.EnvVars)
	}
}

func TestFromUrfaveFloatFlag(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--rate")

	if f.Type != "float64" {
		t.Errorf("rate type = %q", f.Type)
	}
	if f.Default != "1.5" {
		t.Errorf("rate default = %q", f.Default)
	}
}

func TestFromUrfaveDurationFlag(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--timeout")

	if f.Type != "duration" {
		t.Errorf("timeout type = %q", f.Type)
	}
	if f.Default != "30s" {
		t.Errorf("timeout default = %q", f.Default)
	}
}

func TestFromUrfaveStringSliceFlag(t *testing.T) {
	hc := FromUrfave(newTestCommand())
	f := findFlag(t, hc.Flags, "--tags")

	if f.Type != "[]string" {
		t.Errorf("tags type = %q", f.Type)
	}
}

func TestFromUrfaveHiddenFlags(t *testing.T) {
	cmd := newTestCommand()
	hc := FromUrfave(cmd)

	for _, f := range hc.Flags {
		if f.Long == "--debug" {
			if !f.Hidden {
				t.Error("debug flag should be hidden")
			}
			return
		}
	}
	// debug is hidden so might be in flat or groups.
	for _, g := range hc.FlagGroups {
		for _, f := range g.Flags {
			if f.Long == "--debug" {
				if !f.Hidden {
					t.Error("debug flag should be hidden")
				}
				return
			}
		}
	}
	t.Error("debug flag not found")
}

func TestFromUrfaveFlagCategories(t *testing.T) {
	cmd := newTestCommand()
	hc := FromUrfave(cmd)

	if len(hc.FlagGroups) == 0 {
		t.Fatal("expected flag groups from categories")
	}

	var found bool
	for _, g := range hc.FlagGroups {
		if g.Name == "Server" {
			found = true
			if len(g.Flags) < 2 {
				t.Errorf("Server group should have at least 2 flags, got %d", len(g.Flags))
			}
		}
	}
	if !found {
		t.Error("Server flag group not found")
	}
}

func TestFromUrfaveSubcommands(t *testing.T) {
	cmd := newTestCommand()
	hc := FromUrfave(cmd)

	names := make(map[string]bool)
	for _, c := range hc.Commands {
		names[c.Name] = true
	}
	for _, g := range hc.CommandGroups {
		for _, c := range g.Commands {
			names[c.Name] = true
		}
	}

	if !names["serve"] {
		t.Error("missing serve command")
	}
	if !names["build"] {
		t.Error("missing build command")
	}
	if names["internal"] {
		t.Error("hidden command should not appear")
	}
}

func TestFromUrfaveCommandCategories(t *testing.T) {
	cmd := newTestCommand()
	hc := FromUrfave(cmd)

	var found bool
	for _, g := range hc.CommandGroups {
		if g.Name == "Development" {
			found = true
		}
	}
	if !found {
		t.Error("Development command group not found")
	}
}

func TestFromUrfaveNoDescription(t *testing.T) {
	cmd := &cli.Command{
		Name:  "simple",
		Usage: "A simple command",
	}
	hc := FromUrfave(cmd)

	if hc.Description != "A simple command" {
		t.Errorf("Description should fall back to Usage, got %q", hc.Description)
	}
}

func TestFromUrfaveNoAliases(t *testing.T) {
	cmd := &cli.Command{Name: "bare"}
	hc := FromUrfave(cmd)

	if len(hc.Aliases) != 0 {
		t.Errorf("Aliases should be empty, got %v", hc.Aliases)
	}
}

func TestFromUrfaveEmptyCommand(t *testing.T) {
	cmd := &cli.Command{Name: "empty"}
	hc := FromUrfave(cmd)

	if hc.Name != "empty" {
		t.Errorf("Name = %q", hc.Name)
	}
	if len(hc.Flags) != 0 {
		t.Errorf("Flags should be empty, got %d", len(hc.Flags))
	}
	if len(hc.Commands) != 0 {
		t.Errorf("Commands should be empty, got %d", len(hc.Commands))
	}
}

func TestFlagCategoryAllTypes(t *testing.T) {
	tests := []struct {
		name string
		flag cli.Flag
		want string
	}{
		{"string", &cli.StringFlag{Name: "a", Category: "cat"}, "cat"},
		{"bool", &cli.BoolFlag{Name: "b", Category: "cat"}, "cat"},
		{"int", &cli.IntFlag{Name: "c", Category: "cat"}, "cat"},
		{"float", &cli.FloatFlag{Name: "d", Category: "cat"}, "cat"},
		{"duration", &cli.DurationFlag{Name: "e", Category: "cat"}, "cat"},
		{"stringslice", &cli.StringSliceFlag{Name: "f", Category: "cat"}, "cat"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := flagCategory(tc.flag)
			if got != tc.want {
				t.Errorf("flagCategory = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFlagCategoryUnknown(t *testing.T) {
	got := flagCategory(&cli.UintFlag{Name: "x"})
	if got != "" {
		t.Errorf("unknown flag type category should be empty, got %q", got)
	}
}

func TestConvertOneFlagUnknownType(t *testing.T) {
	f := &cli.UintFlag{Name: "x"}
	hf := convertOneFlag(f)
	if hf.Type != "string" {
		t.Errorf("unknown flag type should default to string, got %q", hf.Type)
	}
}

func TestIsHiddenUnknownType(t *testing.T) {
	f := &cli.UintFlag{Name: "x"}
	if isHidden(f) {
		t.Error("unknown flag type should not be hidden")
	}
}
