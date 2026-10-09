package movement

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

var inputMapping = &input.InputMapping{
	Keyboard: []input.InputData{
		{
			Name:  "Key_A",
			Code:  65,
			Label: "A",
		},
		{
			Name:  "Digit_0",
			Code:  0,
			Label: "0",
		},
	},
	Mouse: []input.InputData{
		{
			Name:  "LEFT_MOUSE_BUTTON",
			Code:  1,
			Label: "left mouse button",
		},
		{
			Name:  "RIGHT_MOUSE_BUTTON",
			Code:  2,
			Label: "right mouse button",
		},
	},
}

func Test_MarshalAndUnmarshalInputMovement(t *testing.T) {
	expected := NewInputMovement(inputMapping, defaultSpeed-100)

	found := NewBaseInputMovement(inputMapping)

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_InputMovement, details.Name)
	assert.Equal(t, expected, found)
}

func Test_UnmarshalInputMovement(t *testing.T) {
	tests := []struct {
		name         string
		baseComp     *InputMovement
		data         json.RawMessage
		expectedErr  bool
		expectedComp *InputMovement
	}{
		{
			name:     "valid unmarshal with minimal fields",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": false,
					"allowY": false
				}`,
			),
			expectedErr:  false,
			expectedComp: NewInputMovement(inputMapping, 100),
		},
		{
			name:     "valid unmarshal with full payload",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": true,
					"leftKey": "Key_A",
					"rightKey": "Digit_0",
					"allowY": true,
					"upKey": "LEFT_MOUSE_BUTTON",
					"downKey": "RIGHT_MOUSE_BUTTON"
				}`,
			),
			expectedErr: false,
			expectedComp: &InputMovement{
				inputMapping: inputMapping,
				speed:        100,
				allowX:       true,
				leftKey:      "Key_A",
				rightKey:     "Digit_0",
				allowY:       true,
				upKey:        "LEFT_MOUSE_BUTTON",
				downKey:      "RIGHT_MOUSE_BUTTON",
			},
		},
		{
			name:     "invalid unmarshal with faulty payload",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": false
					"allowY": false
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowX true and valid keys",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": true,
					"leftKey": "Key_A",
					"rightKey": "Digit_0",
					"allowY": false
				}`,
			),
			expectedErr: false,
			expectedComp: &InputMovement{
				inputMapping: inputMapping,
				speed:        100,
				allowX:       true,
				leftKey:      "Key_A",
				rightKey:     "Digit_0",
			},
		},
		{
			name:     "unmarshal with allowX true and invalid left key",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": true,
					"leftKey": "A1",
					"rightKey": "Digit_0",
					"allowY": false
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowX true and invalid right key",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": true,
					"leftKey": "Key_A",
					"rightKey": "01",
					"allowY": false
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowX false and valid keys for x direction",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": false,
					"leftKey": "Key_A",
					"rightKey": "Digit_0"
				}`,
			),
			expectedErr:  false,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowX false and oversee invalid keys for x direction",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowX": false,
					"leftKey": "A1",
					"rightKey": "01"
				}`,
			),
			expectedErr:  false,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowY true and valid keys",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowY": true,
					"upKey": "Key_A",
					"downKey": "Digit_0"
				}`,
			),
			expectedErr: false,
			expectedComp: &InputMovement{
				inputMapping: inputMapping,
				speed:        100,
				allowY:       true,
				upKey:        "Key_A",
				downKey:      "Digit_0",
			},
		},
		{
			name:     "unmarshal with allowY true and invalid up key",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowY": true,
					"upKey": "A1",
					"downKey": "Digit_0"
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowY true and invalid down key",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowY": true,
					"upKey": "Key_A",
					"downKey": "01"
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowY false and valid keys for y direction",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowY": false,
					"upKey": "Key_A",
					"downKey": "Digit_0"
				}`,
			),
			expectedErr:  false,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
		{
			name:     "unmarshal with allowY false and oversee invalid keys for y direction",
			baseComp: NewBaseInputMovement(inputMapping),
			data: json.RawMessage(
				`{
					"speed": 100,
					"allowY": false,
					"upKey": "A1",
					"downKey": "01"
				}`,
			),
			expectedErr:  false,
			expectedComp: NewBaseInputMovement(inputMapping),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.baseComp.BuildFromDetails(tt.data)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expectedComp.speed, tt.baseComp.speed)

			assert.Equal(t, tt.expectedComp.allowX, tt.baseComp.allowX)
			assert.Equal(t, tt.expectedComp.leftKey, tt.baseComp.leftKey)
			assert.Equal(t, tt.expectedComp.rightKey, tt.baseComp.rightKey)

			assert.Equal(t, tt.expectedComp.allowY, tt.baseComp.allowY)
			assert.Equal(t, tt.expectedComp.upKey, tt.baseComp.upKey)
			assert.Equal(t, tt.expectedComp.downKey, tt.baseComp.downKey)
		})
	}
}
