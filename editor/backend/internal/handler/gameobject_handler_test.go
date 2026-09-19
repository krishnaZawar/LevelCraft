package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/entity"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
)

func newTestApp() *fiber.App {
	app := fiber.New()
	app.Put("/gameobjects/:objectID", UpdateGameobject)
	return app
}

func seedGameobject(t *testing.T, name string, group string) *gameobject.Gameobject {
	t.Helper()
	obj := gameobject.NewGameobject()
	obj.SetName(name)
	obj.SetGroup(group)
	gamestatemanager.Get().AddGameobject(obj)
	t.Cleanup(func() { gamestatemanager.Get().DeleteGameobject(obj.GetID()) })
	return obj
}

func updateRequest(t *testing.T, app *fiber.App, id string, body string) (*http.Response, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, "/gameobjects/"+id, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return resp, data
}

func Test_UpdateGameobject(t *testing.T) {
	app := newTestApp()
	obj := seedGameobject(t, "old name", "old group")

	resp, body := updateRequest(t, app, obj.GetID(), `{"name":"new name","group":"new group"}`)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", resp.StatusCode, http.StatusOK, body)
	}

	if obj.GetName() != "new name" {
		t.Fatalf("name = %q, want %q", obj.GetName(), "new name")
	}
	if obj.GetGroup() != "new group" {
		t.Fatalf("group = %q, want %q", obj.GetGroup(), "new group")
	}

	var parsed entity.UpdateGameobjectResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !parsed.Success {
		t.Fatal("success = false, want true")
	}
	if parsed.ObjectDetails[gameobject.Gameobject_CurLabelName] != "new name" {
		t.Fatalf("response details = %v, want the updated name", parsed.ObjectDetails)
	}
}

// An omitted field leaves that value alone: the editor sends only what the
// user actually edited, and the other field must not be blanked as a result.
func Test_UpdateGameobject_OmittedFieldsAreUntouched(t *testing.T) {
	app := newTestApp()
	obj := seedGameobject(t, "keep me", "keep me too")

	if _, body := updateRequest(t, app, obj.GetID(), `{"name":"renamed"}`); body == nil {
		t.Fatal("no response body")
	}

	if obj.GetName() != "renamed" {
		t.Fatalf("name = %q, want %q", obj.GetName(), "renamed")
	}
	if obj.GetGroup() != "keep me too" {
		t.Fatalf("group = %q, want it untouched", obj.GetGroup())
	}
}

// Clearing a field is a real edit, and has to be distinguishable from omitting it.
func Test_UpdateGameobject_EmptyStringClearsTheField(t *testing.T) {
	app := newTestApp()
	obj := seedGameobject(t, "named", "grouped")

	updateRequest(t, app, obj.GetID(), `{"group":""}`)

	if obj.GetGroup() != "" {
		t.Fatalf("group = %q, want it cleared", obj.GetGroup())
	}
	if obj.GetName() != "named" {
		t.Fatalf("name = %q, want it untouched", obj.GetName())
	}
}

func Test_UpdateGameobject_Errors(t *testing.T) {
	app := newTestApp()
	obj := seedGameobject(t, "name", "group")

	tests := []struct {
		name       string
		id         string
		body       string
		wantStatus int
	}{
		{"unknown gameobject", "does-not-exist", `{"name":"x"}`, http.StatusNotFound},
		{"malformed body", obj.GetID(), `not json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := updateRequest(t, app, tt.id, tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", resp.StatusCode, tt.wantStatus, body)
			}
		})
	}
}

func Test_DuplicateName(t *testing.T) {
	tests := []struct {
		name   string
		source string
		taken  []string
		want   string
	}{
		{"first copy", "Player", []string{"Player"}, "Player (1)"},
		{"next free index", "Player", []string{"Player", "Player (1)"}, "Player (2)"},
		// Copying a copy continues the series rather than nesting suffixes.
		{"copy of a copy", "Player (2)", []string{"Player", "Player (2)"}, "Player (1)"},
		{"unnamed stays unnamed", "", []string{""}, ""},
		{"spacing is normalised", "Player  (3)", []string{"Player"}, "Player (1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taken := map[string]struct{}{}
			for _, name := range tt.taken {
				taken[name] = struct{}{}
			}

			if got := duplicateName(tt.source, taken); got != tt.want {
				t.Fatalf("duplicateName(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

func Test_DuplicateGameobject(t *testing.T) {
	app := fiber.New()
	app.Post("/gameobjects/:objectID/duplicate", DuplicateGameobject)

	source := seedGameobject(t, "Player", "actors")
	if !source.AddComponent(component.NewTransform(10, 20, 30, 40)) {
		t.Fatal("failed to attach a Transform to the source")
	}

	req := httptest.NewRequest(http.MethodPost, "/gameobjects/"+source.GetID()+"/duplicate", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", resp.StatusCode, http.StatusOK, body)
	}

	var parsed entity.DuplicateGameobjectResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	cloneID, _ := parsed.ObjectDetails[gameobject.Gameobject_CurLabelID].(string)
	if cloneID == "" || cloneID == source.GetID() {
		t.Fatalf("clone id = %q, want a new id", cloneID)
	}
	t.Cleanup(func() { gamestatemanager.Get().DeleteGameobject(cloneID) })

	if parsed.ObjectDetails[gameobject.Gameobject_CurLabelName] != "Player (1)" {
		t.Fatalf("clone name = %v, want %q", parsed.ObjectDetails[gameobject.Gameobject_CurLabelName], "Player (1)")
	}
	if parsed.ObjectDetails[gameobject.Gameobject_CurLabelGroup] != "actors" {
		t.Fatalf("clone group = %v, want the source's group", parsed.ObjectDetails[gameobject.Gameobject_CurLabelGroup])
	}

	clone, found := gamestatemanager.Get().GetGameobject(cloneID)
	if !found {
		t.Fatal("the clone was not added to the scene")
	}

	// The copy has to own its components: editing one object's Transform must
	// not move the other.
	cloneTransform, found := clone.GetComponent(base.ComponentName_Transform)
	if !found {
		t.Fatal("the clone has no Transform")
	}
	cloneTransform.(*component.Transform).UpdatePosition(999, 999)

	sourceTransform, _ := source.GetComponent(base.ComponentName_Transform)
	if x, y := sourceTransform.(*component.Transform).GetPosition(); x != 10 || y != 20 {
		t.Fatalf("source moved with the clone: got (%d, %d), want (10, 20)", x, y)
	}
}

func Test_DuplicateGameobject_UnknownID(t *testing.T) {
	app := fiber.New()
	app.Post("/gameobjects/:objectID/duplicate", DuplicateGameobject)

	resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/gameobjects/nope/duplicate", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
