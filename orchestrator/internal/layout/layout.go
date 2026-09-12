package layout

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Written by the packaging script next to the orchestrator executable.
// Its presence is what distinguishes a packaged install from a repo checkout.
const manifestFile = "layout.json"

// Layout is where the packaged application keeps the things it spawns.
// Paths are relative to the manifest, so an install can be moved anywhere.
type Layout struct {
	EditorBackend  string `json:"editorBackend"`
	BuilderBackend string `json:"builderBackend"`
	EditorApp      string `json:"editorApp"`
	BuilderApp     string `json:"builderApp"`

	root     string
	packaged bool
}

// IsPackaged reports whether the orchestrator is running from a packaged
// install. A repo checkout has no manifest, and is run from source instead.
func (l *Layout) IsPackaged() bool {
	return l.packaged
}

// Path resolves one of the manifest's entries against the install root
func (l *Layout) Path(relative string) string {
	if filepath.IsAbs(relative) {
		return relative
	}
	return filepath.Join(l.root, relative)
}

// reads the manifest sitting next to the running executable, if there is one
func resolve() *Layout {
	exe, err := os.Executable()
	if err != nil {
		return &Layout{}
	}
	return readManifest(filepath.Dir(exe))
}

// reads and validates the manifest in root, if there is one
func readManifest(root string) *Layout {
	data, err := os.ReadFile(filepath.Join(root, manifestFile))
	if err != nil {
		return &Layout{}
	}

	var l Layout
	if err := json.Unmarshal(data, &l); err != nil {
		return &Layout{}
	}

	// A manifest missing any of its entries cannot describe a working install,
	// so it is treated as absent rather than half-trusted.
	if l.EditorBackend == "" || l.BuilderBackend == "" || l.EditorApp == "" || l.BuilderApp == "" {
		return &Layout{}
	}

	l.root = root
	l.packaged = true
	return &l
}

var (
	once   sync.Once
	layout *Layout
)

// NewPackaged builds a layout for an install rooted at root
func NewPackaged(root string, editorBackend, builderBackend, editorApp, builderApp string) *Layout {
	return &Layout{
		EditorBackend:  editorBackend,
		BuilderBackend: builderBackend,
		EditorApp:      editorApp,
		BuilderApp:     builderApp,
		root:           root,
		packaged:       true,
	}
}

// NewSource builds a layout for a repo checkout, where nothing is prebuilt
func NewSource() *Layout {
	return &Layout{}
}

// Get returns the layout of the current install, resolved once
func Get() *Layout {
	once.Do(func() {
		layout = resolve()
	})
	return layout
}
