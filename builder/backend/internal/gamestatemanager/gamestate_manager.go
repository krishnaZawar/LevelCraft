package gamestatemanager

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
)

var (
	ErrExpectedMapStringInterface = errors.New("GamesceneError: Expected Map string interface")
)

// GameStateManager is a thread safe manager
// used to hold and update the state of the game
type GameStateManager struct {
	mu sync.RWMutex
	// gameobject.id -> gameobject mapping
	gameobjects map[string]*gameobject.Gameobject
}

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		gameobjects: map[string]*gameobject.Gameobject{},
	}
}

// Adds a new gameobject to the scene
func (gsm *GameStateManager) AddGameobject(obj *gameobject.Gameobject) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.gameobjects[obj.GetID()] = obj
}

// Deletes a gameobject from the scene
func (gsm *GameStateManager) DeleteGameobject(id string) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	delete(gsm.gameobjects, id)
}

// fetches a gameobject from the scene
func (gsm *GameStateManager) GetGameobject(id string) (*gameobject.Gameobject, bool) {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()
	obj, ok := gsm.gameobjects[id]
	return obj, ok
}

// returns the entire gamestate of the scene
func (gsm *GameStateManager) GetGameState() []gameobject.GameobjectDetails {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()
	state := []gameobject.GameobjectDetails{}

	for _, obj := range gsm.gameobjects {
		state = append(state, obj.GetGameobjectDetails())
	}

	return state
}

// Builds the game scene from the details passed
func (gsm *GameStateManager) BuildFromDetails(data []byte) error {
	var scene []gameobject.GameobjectDetails
	err := json.Unmarshal(data, &scene)
	if err != nil {
		return err
	}
	for _, objData := range scene {
		obj := gameobject.NewGameobjectWithID(objData.Id)
		err := obj.BuildFromDetails(objData)
		if err != nil {
			return err
		}
		gsm.AddGameobject(obj)
	}
	return nil
}

var gsm = NewGameStateManager()

func Get() *GameStateManager {
	return gsm
}
