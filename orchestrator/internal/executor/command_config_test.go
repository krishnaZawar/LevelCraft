package executor

import (
	"strings"
	"testing"
)

// Every process is run from source: the orchestrator drives a repo checkout,
// so the backends go through `go run` and the clients through npm.
func Test_BackendConfigs_RunFromSource(t *testing.T) {
	editor := editorBackendConfig("4000")
	if editor.Name != "go" || editor.Pwd != "../editor/backend" {
		t.Fatalf("editor config = %+v, want go in ../editor/backend", editor)
	}
	if got := strings.Join(editor.Args, " "); got != "run cmd/main.go --port 4000" {
		t.Fatalf("editor Args = %q, want the port to reach the binary", got)
	}
	if editor.Port != "4000" {
		t.Fatalf("editor Port = %q, want 4000", editor.Port)
	}

	builder := builderBackendConfig("4000", "/tmp/scene.json")
	if builder.Name != "go" || builder.Pwd != "../builder/backend" {
		t.Fatalf("builder config = %+v, want go in ../builder/backend", builder)
	}
}

// The scene path is what makes a builder run a particular game, so it has to
// survive into the command.
func Test_BuilderBackendConfig_CarriesTheScenePath(t *testing.T) {
	comm := builderBackendConfig("4000", "/tmp/scene.json")

	args := strings.Join(comm.Args, " ")
	if !strings.Contains(args, "--file-name /tmp/scene.json") {
		t.Fatalf("Args = %v, want them to carry the scene path", comm.Args)
	}
	if !strings.Contains(args, "--port 4000") {
		t.Fatalf("Args = %v, want them to carry the port", comm.Args)
	}
}

func Test_FrontendConfigs_RunFromSource(t *testing.T) {
	editor := editorFrontendConfig("5000")
	if editor.Name != "npm" || editor.Pwd != "../editor/app" {
		t.Fatalf("editor config = %+v, want npm in ../editor/app", editor)
	}
	if editor.Port != "5000" {
		t.Fatalf("editor Port = %q, want 5000", editor.Port)
	}

	builder := builderFrontendConfig("5000")
	if builder.Name != "npm" || builder.Pwd != "../builder/app" {
		t.Fatalf("builder config = %+v, want npm in ../builder/app", builder)
	}
}
