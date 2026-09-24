package component

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_GetColor(t *testing.T) {
	r, g, b, a := 100, 100, 100, 100
	color := NewColor(r, g, b, a)

	actualR, actualG, actualB, actualA := color.Get()

	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_SetColor(t *testing.T) {
	color := newBaseColor()

	r, g, b, a := 100, 100, 100, 100

	color.Set(r, g, b, a)

	actualR, actualG, actualB, actualA := color.Get()

	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_MarshalAndUnmarshalColor(t *testing.T) {
	r, g, b, a := 100, 120, 140, 160
	expected := NewColor(r, g, b, a)

	found := newBaseColor()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_Color, details.Name)
	assert.Equal(t, expected, found)
}

func Test_UnmarshalColor(t *testing.T) {
	tests := []struct {
		name          string
		baseComp      *Color
		data          json.RawMessage
		expectedErr   bool
		expectedColor *Color
	}{
		{
			name:     "valid unmarshal",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 100,
					"g": 120,
					"b": 140,
					"a": 160
				}`,
			),
			expectedErr:   false,
			expectedColor: NewColor(100, 120, 140, 160),
		},
		{
			name:     "invalid unmarshal",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 100,
					"g": 120,
					"b": 140
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "r is below bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": -1,
					"g": 120,
					"b": 140,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "r is above bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 256,
					"g": 120,
					"b": 140,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "g is below bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 0,
					"g": -1,
					"b": 140,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "g is above bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 255,
					"g": 256,
					"b": 140,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "b is below bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 0,
					"g": 0,
					"b": -1,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "b is above bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 255,
					"g": 255,
					"b": 256,
					"a": 160
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "a is below bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 0,
					"g": 0,
					"b": 0,
					"a": -1
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
		},
		{
			name:     "a is above bounds",
			baseComp: newBaseColor(),
			data: json.RawMessage(
				`{
					"r": 255,
					"g": 255,
					"b": 255,
					"a": 256
				}`,
			),
			expectedErr:   true,
			expectedColor: newBaseColor(),
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
			assert.Equal(t, tt.expectedColor.r, tt.baseComp.r)
			assert.Equal(t, tt.expectedColor.g, tt.baseComp.g)
			assert.Equal(t, tt.expectedColor.b, tt.baseComp.b)
			assert.Equal(t, tt.expectedColor.a, tt.baseComp.a)
		})
	}
}
