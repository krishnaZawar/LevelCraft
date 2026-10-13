package eventgenerator

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

var inputMapping = &input.InputMapping{
	Keyboard: []input.InputData{
		{
			Name:  "KEY_A",
			Code:  1,
			Label: "A",
		},
		{
			Name:  "KEY_D",
			Code:  2,
			Label: "D",
		},
		{
			Name:  "KEY_S",
			Code:  3,
			Label: "S",
		},
		{
			Name:  "KEY_W",
			Code:  4,
			Label: "W",
		},
	},
}

var compRegistry = component.NewComponentRegistry(inputMapping)

func Test_InputEventGeneratorHandle(t *testing.T) {
	tests := []struct {
		name           string
		objID          string
		ctxObjID       string
		input          string
		comp           component.Component
		compData       json.RawMessage
		expectedEvents []models.Event
		expectedErr    error
	}{
		{
			name:           "component passed is nil",
			objID:          "123",
			ctxObjID:       "123",
			input:          "KEY_A",
			comp:           nil,
			compData:       json.RawMessage(`{}`),
			expectedEvents: nil,
			expectedErr:    base.ErrNilComponentPassed,
		},
		{
			name:           "incompatible component type passed",
			objID:          "123",
			ctxObjID:       "123",
			input:          "KEY_A",
			comp:           component.NewTransform(100.0, 100.0, 100.0, 100.0),
			compData:       json.RawMessage(`{}`),
			expectedEvents: nil,
			expectedErr:    base.NewErrComponentTypeMisMatch(compbase.ComponentName_InputMovement, compbase.ComponentName_Transform),
		},
		{
			name:           "gameobject not found for context",
			objID:          "123",
			ctxObjID:       "321",
			input:          "KEY_A",
			comp:           component.NewInputMovement(inputMapping, 100.0),
			compData:       json.RawMessage(`{}`),
			expectedEvents: nil,
			expectedErr:    base.NewErrGameobjectNotFound("321"),
		},
		{
			name:     "valid gameobject state and eventCtx passed",
			objID:    "123",
			ctxObjID: "123",
			input:    "KEY_D",
			comp:     component.NewInputMovement(inputMapping, 100.0),
			compData: json.RawMessage(`{
				"speed": 100.0,
				"allowX": true,
				"leftKey": "KEY_A",
				"rightKey": "KEY_D"
			}`),
			expectedEvents: []models.Event{
				&event.MoveEvent{
					GameobjectID: "123",
					MoveX:        100.0,
					MoveY:        0.0,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gsm := gamestatemanager.NewGameStateManager()

			obj := gameobject.NewGameobjectWithID(tt.objID, compRegistry)

			if tt.comp != nil {
				err := tt.comp.BuildFromDetails(tt.compData)
				assert.Nil(t, err)
				ok := obj.AddComponent(tt.comp)
				assert.Equal(t, true, ok)
			}

			gsm.AddGameobject(obj)

			generator := NewInputMovementEventGenerator(gsm)

			eventCtx := EventGeneratorCtx{
				GameobjectID: tt.ctxObjID,
				Component:    tt.comp,
				InputKey:     tt.input,
			}

			events, err := generator.GenerateEvents(eventCtx)
			assert.ElementsMatch(t, tt.expectedEvents, events)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
