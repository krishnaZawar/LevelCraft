package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/handler"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/router"
)

func main() {
	app := fiber.New()

	fileName, ok := getFilenameForArgs()
	if fileName == "" {
		panic("game file should be provided to run the game")
	}
	if !ok {
		panic("game file provided should be json")
	}

	err := initGameScene(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to load game scene due to %+v", err))
	}

	router.CreateRoutes(app)

	go handler.GameLoop()

	app.Listen(":8000")
}

// reads the CLI args passed and enforces whether the file was passed and was it a json
// program panics on failure
func getFilenameForArgs() (string, bool) {
	fileName := flag.String("file-name", "", "name of the game file (should be json)")

	flag.Parse()

	if *fileName == "" {
		return "", false
	}

	if filepath.Ext(*fileName) != ".json" {
		return *fileName, false
	}
	return *fileName, true
}

// reads the file and initializes the game scene
// program panics on failure
func initGameScene(filepath string) error {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var sceneData map[string]interface{}
	if err := json.Unmarshal(fileData, &sceneData); err != nil {
		return err
	}

	if err := gamestatemanager.Get().BuildFromDetails(sceneData); err != nil {
		return err
	}

	return nil
}
