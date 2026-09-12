package executor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/layout"
)

func packagedLayout() *layout.Layout {
	return layout.NewPackaged(
		"/opt/levelcraft",
		filepath.Join("bin", "editor-backend"),
		filepath.Join("bin", "builder-backend"),
		filepath.Join("apps", "editor", "LevelCraft"),
		filepath.Join("apps", "builder", "Player"),
	)
}

// A packaged install runs prebuilt binaries; a checkout runs from source. The
// same process has to come up either way, so both configs are checked together.
func Test_CommandConfigs(t *testing.T) {
	tests := []struct {
		name     string
		packaged bool
		wantName string
		wantArgs []string
	}{
		{
			name:     "editor backend from source",
			wantName: "go",
			wantArgs: []string{"run", "cmd/main.go", "--port", "4000"},
		},
		{
			name:     "editor backend packaged",
			packaged: true,
			wantName: filepath.Join("/opt/levelcraft", "bin", "editor-backend"),
			wantArgs: []string{"--port", "4000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := layout.NewSource()
			if tt.packaged {
				l = packagedLayout()
			}

			comm := editorBackendConfig(l, "4000")

			if comm.Name != tt.wantName {
				t.Fatalf("Name = %q, want %q", comm.Name, tt.wantName)
			}
			if strings.Join(comm.Args, " ") != strings.Join(tt.wantArgs, " ") {
				t.Fatalf("Args = %v, want %v", comm.Args, tt.wantArgs)
			}
			if comm.Port != "4000" {
				t.Fatalf("Port = %q, want 4000", comm.Port)
			}
		})
	}
}

// The scene path is what makes a builder run a particular game, so it has to
// survive into the command in both layouts.
func Test_BuilderBackendConfig_CarriesTheScenePath(t *testing.T) {
	for _, l := range []*layout.Layout{layout.NewSource(), packagedLayout()} {
		comm := builderBackendConfig(l, "4000", "/tmp/scene.json")

		args := strings.Join(comm.Args, " ")
		if !strings.Contains(args, "--file-name /tmp/scene.json") {
			t.Fatalf("Args = %v, want them to carry the scene path", comm.Args)
		}
		if !strings.Contains(args, "--port 4000") {
			t.Fatalf("Args = %v, want them to carry the port", comm.Args)
		}
	}
}

// A packaged client is launched directly, with no npm or package scripts
// involved, so nothing may leak in from the source layout.
func Test_FrontendConfigs_Packaged(t *testing.T) {
	l := packagedLayout()

	for _, comm := range []struct {
		name     string
		got      string
		wantPath string
		args     []string
	}{
		{"editor", editorFrontendConfig(l, "5000").Name, filepath.Join("apps", "editor", "LevelCraft"), editorFrontendConfig(l, "5000").Args},
		{"builder", builderFrontendConfig(l, "5000").Name, filepath.Join("apps", "builder", "Player"), builderFrontendConfig(l, "5000").Args},
	} {
		want := filepath.Join("/opt/levelcraft", comm.wantPath)
		if comm.got != want {
			t.Fatalf("%s Name = %q, want %q", comm.name, comm.got, want)
		}
		if len(comm.args) != 0 {
			t.Fatalf("%s Args = %v, want none", comm.name, comm.args)
		}
	}
}

func Test_FrontendConfigs_FromSource(t *testing.T) {
	l := layout.NewSource()

	editor := editorFrontendConfig(l, "5000")
	if editor.Name != "npm" || editor.Pwd != "../editor/app" {
		t.Fatalf("editor config = %+v, want npm in ../editor/app", editor)
	}

	builder := builderFrontendConfig(l, "5000")
	if builder.Name != "npm" || builder.Pwd != "../builder/app" {
		t.Fatalf("builder config = %+v, want npm in ../builder/app", builder)
	}
}
