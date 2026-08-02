package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()

	app.Post("/save", HandleSaveGame)

	return app
}

func Test_SaveGame(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	obj := gameobject.NewGameobject()
	gsm.AddGameobject(obj)

	/*
		expected state:
		{
			"<obj-id>" : {
				"id": "<obj-id>",
				"name": "",
				"group": "",
				"components": {}
			}
		}

	*/
	expectedGameState := map[string]interface{}{
		obj.GetID(): obj.GetGameobjectDetails(),
	}

	fp := filepath.Join(t.TempDir(), "game.json")

	err := saveGame(gsm, fp)
	if err != nil {
		t.Errorf("%+v", err)
	}

	data, err := os.ReadFile(fp)
	if err != nil {
		t.Errorf("Failed to read file: %+v", err)
	}

	var gameState map[string]interface{}
	err = json.Unmarshal(data, &gameState)
	if err != nil {
		t.Errorf("Failed to unmarshal game state: %+v", err)
	}

	assert.Equal(t, expectedGameState, gameState)
}

func Test_SaveGameHandler(t *testing.T) {
	app := setupApp()

	t.Run("json file not passed", func(t *testing.T) {
		payload := map[string]interface{}{
			"filepath": "game.txt",
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Errorf("failed to marshal payload")
		}

		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Nil(t, err)
	})

	t.Run("json file passed", func(t *testing.T) {
		fp := filepath.Join(t.TempDir(), "game.json")
		payload := map[string]interface{}{
			"filepath": fp,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Errorf("failed to marshal payload")
		}

		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
