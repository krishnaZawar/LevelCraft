package gamestatemanager

import (
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
func (gsm *GameStateManager) GetGameState() map[string]interface{} {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()
	state := map[string]interface{}{}

	for id, obj := range gsm.gameobjects {
		state[id] = obj.GetGameobjectDetails()
	}

	return state
}

// Builds the game scene from the details passed
func (gsm *GameStateManager) BuildFromDetails(scene map[string]interface{}) error {
	for id, objData := range scene {
		obj := gameobject.NewGameobjectWithID(id)
		data, ok := objData.(map[string]interface{})
		if !ok {
			return ErrExpectedMapStringInterface
		}
		err := obj.BuildFromDetails(data)
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
