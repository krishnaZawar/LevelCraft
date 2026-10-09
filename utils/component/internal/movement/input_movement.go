package movement

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/input"
)

const (
	// default speed value the components hold on init
	defaultSpeed = 100
)

// Allows the player to control the object's movement using input keys.
//
// The user can also restrict movement to particular axes and control the inputKeys used
type InputMovement struct {
	speed int // the speed with which the object should move

	// X direction related config
	// Input Keys for X direction can only be mapped if movement is allowed on that axis
	allowX   bool   // allow movement on X axis
	leftKey  string // key used for moving in left direction
	rightKey string // key used for moving in right direction

	// Y direction related config
	// Input Keys for Y direction can only be mapped if movement is allowed on that axis
	allowY  bool   // allow movement on Y axis
	upKey   string // key used	for moving in upward direction
	downKey string // key used for moving in downward direction

	inputMapping *input.InputMapping // store input mapping internally to for validating the keys passed
}

// intermediary structure of InputMovement used for marshalling and unmarshalling component details
type inputMovementJSON struct {
	Speed int `json:"speed"`

	AllowX   bool   `json:"allowX"`
	LeftKey  string `json:"leftKey,omitempty"`
	RightKey string `json:"rightKey,omitempty"`

	AllowY  bool   `json:"allowY"`
	UpKey   string `json:"upKey,omitempty"`
	DownKey string `json:"downKey,omitempty"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBaseInputMovement(inputMapping *input.InputMapping) *InputMovement {
	return &InputMovement{
		speed:        defaultSpeed,
		inputMapping: inputMapping,
	}
}

func NewInputMovement(inputMapping *input.InputMapping, speed int) *InputMovement {
	inputMovement := NewBaseInputMovement(inputMapping)
	inputMovement.speed = speed
	return inputMovement
}

// Returns the name of the component
func (im *InputMovement) GetComponentName() string {
	return base.ComponentName_InputMovement
}

// Returns a snapshot of the complete data stored in the component
func (im *InputMovement) GetComponentDetails() models.ComponentDetails {
	data := inputMovementJSON{
		Speed: im.speed,

		AllowX:   im.allowX,
		LeftKey:  im.leftKey,
		RightKey: im.rightKey,

		AllowY:  im.allowY,
		UpKey:   im.upKey,
		DownKey: im.downKey,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
		Name: im.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (im *InputMovement) BuildFromDetails(data json.RawMessage) error {
	var componentData inputMovementJSON
	err := json.Unmarshal(data, &componentData)
	if err != nil {
		return err
	}

	if !componentData.AllowX {
		componentData.LeftKey = ""
		componentData.RightKey = ""
	} else {
		if err := im.inputMapping.KeyExists(componentData.LeftKey); err != nil {
			return err
		}
		if err := im.inputMapping.KeyExists(componentData.RightKey); err != nil {
			return err
		}
	}

	if !componentData.AllowY {
		componentData.UpKey = ""
		componentData.DownKey = ""
	} else {
		if err := im.inputMapping.KeyExists(componentData.UpKey); err != nil {
			return err
		}
		if err := im.inputMapping.KeyExists(componentData.DownKey); err != nil {
			return err
		}
	}

	im.speed = componentData.Speed

	im.allowX = componentData.AllowX
	im.leftKey = componentData.LeftKey
	im.rightKey = componentData.RightKey

	im.allowY = componentData.AllowY
	im.upKey = componentData.UpKey
	im.downKey = componentData.DownKey

	return nil
}

// can interact only with the keys
func (im *InputMovement) HandlesInput(input string) bool {
	if im.leftKey == input || im.rightKey == input {
		return im.allowX
	} else if im.upKey == input || im.downKey == input {
		return im.allowY
	}
	return false
}

var _ models.Component = &InputMovement{}
