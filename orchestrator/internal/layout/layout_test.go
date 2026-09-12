package layout

import (
	"os"
	"path/filepath"
	"testing"
)

func writeManifest(t *testing.T, dir string, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, manifestFile), []byte(contents), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
}

func Test_Resolve_NoManifest(t *testing.T) {
	// The test binary's own directory has no manifest, which is exactly the
	// repo-checkout case: the orchestrator should run from source.
	l := resolve()

	if l.IsPackaged() {
		t.Fatal("IsPackaged() = true without a manifest, want false")
	}
}

// A half-written manifest can't describe a working install, so it must not be
// trusted into spawning binaries that aren't there.
func Test_Resolve_IncompleteManifestIsNotPackaged(t *testing.T) {
	tests := []struct {
		name     string
		contents string
	}{
		{"missing entries", `{"editorBackend":"bin/editor-backend"}`},
		{"empty entry", `{"editorBackend":"a","builderBackend":"b","editorApp":"","builderApp":"d"}`},
		{"malformed json", `{`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeManifest(t, dir, tt.contents)

			if l := readManifest(dir); l.IsPackaged() {
				t.Fatalf("IsPackaged() = true for %s, want false", tt.name)
			}
		})
	}
}

func Test_Resolve_CompleteManifest(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"editorBackend": "bin/editor-backend",
		"builderBackend": "bin/builder-backend",
		"editorApp": "apps/editor/LevelCraft.app/Contents/MacOS/LevelCraft",
		"builderApp": "apps/builder/Player.app/Contents/MacOS/Player"
	}`)

	l := readManifest(dir)

	if !l.IsPackaged() {
		t.Fatal("IsPackaged() = false for a complete manifest, want true")
	}

	want := filepath.Join(dir, "bin", "editor-backend")
	if got := l.Path(l.EditorBackend); got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
}

// An install can be moved anywhere, so relative entries resolve against the
// manifest, while an absolute entry is taken as given.
func Test_Path_AbsoluteEntryIsUntouched(t *testing.T) {
	dir := t.TempDir()
	absolute := filepath.Join(dir, "elsewhere", "editor-backend")
	l := &Layout{root: "/install/root"}

	if got := l.Path(absolute); got != absolute {
		t.Fatalf("Path() = %q, want %q", got, absolute)
	}
}
