package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/handler"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/router"
)

func main() {
	app := fiber.New()

	port := flag.Int("port", 3000, "the value of the port on which the service starts")
	fileName := flag.String("file-name", "", "name of the game file (should be json)")

	flag.Parse()

	// The only consumer of this API is our own player client running locally
	// (dev server origin or a file:// packaged build), never a public caller,
	// so a permissive CORS policy is appropriate here. The websocket needed
	// none of this; the input-mapping request is the first plain HTTP call the
	// player makes, and a browser holds that to the same-origin policy.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))

	err := inputmanager.LoadMapping()
	if err != nil {
		panic(err)
	}

	err = verifyFilename(*fileName)
	if err != nil {
		panic(err)
	}

	err = initGameScene(*fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to load game scene due to %+v", err))
	}

	router.CreateRoutes(app)

	go handler.GameLoop()

	err = app.Listen(":" + strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}
}

// reads the CLI args passed and enforces whether the file was passed and was it a json
// program panics on failure
func verifyFilename(fileName string) error {
	if fileName == "" {
		return errors.New("game file should be provided to run the game")
	}

	if filepath.Ext(fileName) != ".json" {
		return errors.New("game file provided should be json")
	}
	return nil
}

// reads the file and initializes the game scene
// program panics on failure
func initGameScene(filepath string) error {
	sceneData, err := os.ReadFile(filepath) //nolint:gosec // filepath is user-controlled and required to load the scene but enough validation exists for safety
	if err != nil {
		return err
	}

	if err := gamestatemanager.Get().BuildFromDetails(sceneData); err != nil {
		return err
	}

	return nil
}
