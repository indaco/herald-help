package heraldkong

import (
	"testing"

	"github.com/alecthomas/kong"
)

type testCLI struct {
	Debug   bool   `help:"Enable debug mode" short:"d" env:"DEBUG"`
	Output  string `help:"Output file" short:"o" default:"stdout" enum:"stdout,file"`
	Verbose int    `help:"Verbosity level" default:"0"`
	Hidden  string `help:"Hidden flag" hidden:""`

	Serve serveCmd `cmd:"" help:"Start server"`
	Build buildCmd `cmd:"" help:"Build the project"`
}

type serveCmd struct {
	Port int `help:"Port number" default:"8080" env:"PORT"`
}

type buildCmd struct {
	Target string `arg:"" help:"Build target" default:"all"`
}

func newTestApp(t *testing.T) *kong.Kong {
	t.Helper()
	var cli testCLI
	parser, err := kong.New(&cli, kong.Name("myapp"), kong.Description("A sample app"))
	if err != nil {
		t.Fatalf("kong.New: %v", err)
	}
	return parser
}

func TestFromKong(t *testing.T) {
	app := newTestApp(t)
	hc := FromKong(app)

	if hc.Name != "myapp" {
		t.Errorf("Name = %q, want %q", hc.Name, "myapp")
	}
	if hc.Description != "A sample app" {
		t.Errorf("Description = %q", hc.Description)
	}
	if hc.Synopsis == "" {
		t.Error("Synopsis should not be empty")
	}
}

func TestFromKongFlags(t *testing.T) {
	app := newTestApp(t)
	hc := FromKong(app)

	flagMap := make(map[string]struct {
		Short   string
		Type    string
		Default string
		Hidden  bool
		EnvVars []string
		Enum    []string
	})
	for _, f := range hc.Flags {
		flagMap[f.Long] = struct {
			Short   string
			Type    string
			Default string
			Hidden  bool
			EnvVars []string
			Enum    []string
		}{f.Short, f.Type, f.Default, f.Hidden, f.EnvVars, f.Enum}
	}

	t.Run("debug flag", func(t *testing.T) {
		f, ok := flagMap["--debug"]
		if !ok {
			t.Fatal("--debug not found")
		}
		if f.Short != "-d" {
			t.Errorf("debug short = %q", f.Short)
		}
		if f.Type != "bool" {
			t.Errorf("debug type = %q", f.Type)
		}
		if len(f.EnvVars) != 1 || f.EnvVars[0] != "DEBUG" {
			t.Errorf("debug envvars = %v", f.EnvVars)
		}
	})

	t.Run("output flag", func(t *testing.T) {
		f, ok := flagMap["--output"]
		if !ok {
			t.Fatal("--output not found")
		}
		if f.Short != "-o" {
			t.Errorf("output short = %q", f.Short)
		}
		if f.Default != "stdout" {
			t.Errorf("output default = %q", f.Default)
		}
		if len(f.Enum) != 2 {
			t.Errorf("output enum = %v", f.Enum)
		}
	})

	t.Run("verbose flag", func(t *testing.T) {
		f, ok := flagMap["--verbose"]
		if !ok {
			t.Fatal("--verbose not found")
		}
		if f.Type != "int" {
			t.Errorf("verbose type = %q", f.Type)
		}
	})

	t.Run("hidden flag", func(t *testing.T) {
		f, ok := flagMap["--hidden"]
		if !ok {
			t.Fatal("--hidden not found")
		}
		if !f.Hidden {
			t.Error("hidden flag should be hidden")
		}
	})
}

func TestFromKongSubcommands(t *testing.T) {
	app := newTestApp(t)
	hc := FromKong(app)

	names := make(map[string]bool)
	for _, c := range hc.Commands {
		names[c.Name] = true
	}

	if !names["serve"] {
		t.Error("missing serve command")
	}
	if !names["build"] {
		t.Error("missing build command")
	}
}

func TestFromKongSubcommandHelp(t *testing.T) {
	app := newTestApp(t)

	var serveNode *kong.Node
	for _, child := range app.Model.Children {
		if child != nil && child.Name == "serve" {
			serveNode = child
			break
		}
	}

	if serveNode == nil {
		t.Fatal("serve node not found")
	}

	hc := FromNode(serveNode)

	if hc.Name != "serve" {
		t.Errorf("Name = %q", hc.Name)
	}

	var found bool
	for _, f := range hc.Flags {
		if f.Long == "--port" {
			found = true
			if f.Default != "8080" {
				t.Errorf("port default = %q", f.Default)
			}
			if len(f.EnvVars) != 1 || f.EnvVars[0] != "PORT" {
				t.Errorf("port envvars = %v", f.EnvVars)
			}
		}
	}
	if !found {
		t.Error("--port flag not found in serve command")
	}
}

func TestFromKongPositionalArgs(t *testing.T) {
	app := newTestApp(t)

	var buildNode *kong.Node
	for _, child := range app.Model.Children {
		if child != nil && child.Name == "build" {
			buildNode = child
			break
		}
	}

	if buildNode == nil {
		t.Fatal("build node not found")
	}

	hc := FromNode(buildNode)

	if len(hc.Args) != 1 {
		t.Fatalf("Args count = %d, want 1", len(hc.Args))
	}
	if hc.Args[0].Name != "<target>" {
		t.Errorf("arg name = %q", hc.Args[0].Name)
	}
	if hc.Args[0].Desc != "Build target" {
		t.Errorf("arg desc = %q", hc.Args[0].Desc)
	}
}

func TestFromKongSynopsis(t *testing.T) {
	app := newTestApp(t)
	hc := FromKong(app)

	if hc.Synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if hc.Synopsis[:5] != "myapp" {
		t.Errorf("Synopsis should start with app name: %q", hc.Synopsis)
	}
}

func TestFromKongNoDetail(t *testing.T) {
	var cli struct{}
	parser, err := kong.New(&cli, kong.Name("bare"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	if hc.Name != "bare" {
		t.Errorf("Name = %q", hc.Name)
	}
}

func TestFromKongWithDetail(t *testing.T) {
	var cli struct{}
	parser, err := kong.New(&cli,
		kong.Name("detailed"),
		kong.Description("Short help."),
	)
	if err != nil {
		t.Fatal(err)
	}
	parser.Model.Detail = "Detailed description."

	hc := FromKong(parser)
	if hc.Description != "Detailed description." {
		t.Errorf("Description = %q, want detailed", hc.Description)
	}
}

type groupedCLI struct {
	Host string `help:"Host" group:"server" default:"localhost"`
	Port int    `help:"Port" group:"server" default:"8080"`
	Name string `help:"Name" default:"app"`

	Run  runCmd  `cmd:"" help:"Run the app" group:"execution"`
	Test testCmd `cmd:"" help:"Test the app"`
}

type runCmd struct{}
type testCmd struct{}

func TestFromKongFlagGroups(t *testing.T) {
	var cli groupedCLI
	parser, err := kong.New(&cli, kong.Name("grouped"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)

	var nameFound bool
	for _, f := range hc.Flags {
		if f.Long == "--name" {
			nameFound = true
		}
	}
	if !nameFound {
		t.Error("ungrouped --name flag not found in flat flags")
	}

	if len(hc.FlagGroups) == 0 {
		t.Fatal("expected flag groups")
	}

	var serverGroup bool
	for _, g := range hc.FlagGroups {
		if g.Name == "server" {
			serverGroup = true
			if len(g.Flags) != 2 {
				t.Errorf("server group should have 2 flags, got %d", len(g.Flags))
			}
		}
	}
	if !serverGroup {
		t.Error("server flag group not found")
	}
}

func TestFromKongCommandGroups(t *testing.T) {
	var cli groupedCLI
	parser, err := kong.New(&cli, kong.Name("grouped"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)

	var testFound bool
	for _, c := range hc.Commands {
		if c.Name == "test" {
			testFound = true
		}
	}
	if !testFound {
		t.Error("ungrouped test command not found")
	}

	var execGroup bool
	for _, g := range hc.CommandGroups {
		if g.Name == "execution" {
			execGroup = true
		}
	}
	if !execGroup {
		t.Error("execution command group not found")
	}
}

type hiddenCLI struct {
	Visible  visibleCmd  `cmd:"" help:"Visible command"`
	Internal internalCmd `cmd:"" help:"Internal" hidden:""`
}

type visibleCmd struct{}
type internalCmd struct{}

func TestFromKongHiddenCommands(t *testing.T) {
	var cli hiddenCLI
	parser, err := kong.New(&cli, kong.Name("hidden"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)

	names := make(map[string]bool)
	for _, c := range hc.Commands {
		names[c.Name] = true
	}

	if !names["visible"] {
		t.Error("visible command should appear")
	}
	if names["internal"] {
		t.Error("hidden command should not appear")
	}
}

func TestFromKongNoAliases(t *testing.T) {
	var cli struct{}
	parser, err := kong.New(&cli, kong.Name("noalias"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	if len(hc.Aliases) != 0 {
		t.Errorf("Aliases should be empty, got %v", hc.Aliases)
	}
}

type counterCLI struct {
	Verbose int `help:"Verbosity" short:"v" type:"counter"`
}

func TestFromKongCounterFlag(t *testing.T) {
	var cli counterCLI
	parser, err := kong.New(&cli, kong.Name("counter"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)

	var found bool
	for _, f := range hc.Flags {
		if f.Long == "--verbose" {
			found = true
			if f.Type != "counter" {
				t.Errorf("verbose type = %q, want counter", f.Type)
			}
		}
	}
	if !found {
		t.Error("--verbose flag not found")
	}
}

type aliasCLI struct {
	Serve aliasServeCmd `cmd:"" aliases:"s,sv" help:"Start"`
}

type aliasServeCmd struct{}

func TestFromKongCommandAliases(t *testing.T) {
	var cli aliasCLI
	parser, err := kong.New(&cli, kong.Name("alias"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)

	for _, c := range hc.Commands {
		if c.Name == "serve" {
			if len(c.Aliases) == 0 {
				t.Error("serve command should have aliases")
			}
			break
		}
	}

	for _, child := range parser.Model.Children {
		if child != nil && child.Name == "serve" {
			hc2 := FromNode(child)
			if len(hc2.Aliases) == 0 {
				t.Error("serve node should have aliases via FromNode")
			}
			return
		}
	}
	t.Error("serve node not found")
}

func TestFromKongNoFlagsSynopsis(t *testing.T) {
	app := newTestApp(t)

	var buildNode *kong.Node
	for _, child := range app.Model.Children {
		if child != nil && child.Name == "build" {
			buildNode = child
			break
		}
	}

	if buildNode == nil {
		t.Fatal("build node not found")
	}

	hc := FromNode(buildNode)

	if hc.Synopsis == "" {
		t.Error("synopsis should not be empty")
	}
}

func TestFromKongEmptyNode(t *testing.T) {
	var cli struct{}
	parser, err := kong.New(&cli, kong.Name("empty"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	if hc.Name != "empty" {
		t.Errorf("Name = %q", hc.Name)
	}
	if hc.Synopsis == "" {
		t.Error("synopsis should not be empty")
	}
}

func TestFromKongLeafCommand(t *testing.T) {
	type leafCLI struct {
		Run leafRunCmd `cmd:"" help:"Run something"`
	}
	var cli leafCLI
	parser, err := kong.New(&cli, kong.Name("leaf"), kong.NoDefaultHelp())
	if err != nil {
		t.Fatal(err)
	}

	for _, child := range parser.Model.Children {
		if child != nil && child.Name == "run" {
			hc := FromNode(child)
			if hc.Synopsis != "run" {
				t.Errorf("leaf synopsis = %q, want just 'run'", hc.Synopsis)
			}
			return
		}
	}
	t.Error("run node not found")
}

type leafRunCmd struct{}

func TestFromKongOptionalPositional(t *testing.T) {
	type optCLI struct {
		File string `arg:"" help:"File" optional:""`
	}

	var cli optCLI
	parser, err := kong.New(&cli, kong.Name("opt"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	if hc.Synopsis == "" {
		t.Error("synopsis should not be empty")
	}
}

func TestFromKongRequiredPositional(t *testing.T) {
	type reqCLI struct {
		File string `arg:"" help:"Input file" required:""`
	}

	var cli reqCLI
	parser, err := kong.New(&cli, kong.Name("req"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	if hc.Synopsis == "" {
		t.Error("synopsis should not be empty")
	}
	if hc.Args[0].Required != true {
		t.Error("file arg should be required")
	}
}

func TestFromKongArgumentNode(t *testing.T) {
	type argBranchCLI struct {
		Target struct {
			Target string   `arg:"" help:"Target name"`
			Run    struct{} `cmd:"" help:"Run it"`
		} `arg:"" help:"Build target"`
	}

	var cli argBranchCLI
	parser, err := kong.New(&cli, kong.Name("argbranch"))
	if err != nil {
		t.Fatal(err)
	}

	hc := FromKong(parser)
	for _, c := range hc.Commands {
		if c.Name == "target" {
			t.Error("argument node should not appear as a command")
		}
	}
}
