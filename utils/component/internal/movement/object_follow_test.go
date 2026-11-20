package movement

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_MarshalAndUnmarshalObjectFollow(t *testing.T) {
	expected := NewObjectFollow(defaultSpeed-10, "obj123", 10)
	expected.SetStopDistance(12)

	found := NewBaseObjectFollow()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_ObjectFollow, details.Name)
	assert.Equal(t, expected, found)
}

func Test_ObjectFollowGetAndSetStopDistance(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		stopDistVal := 12.0
		comp := NewBaseObjectFollow()
		comp.SetStopDistance(stopDistVal)
		assert.Equal(t, stopDistVal, comp.GetStopDistance())
	})
	t.Run("negative value", func(t *testing.T) {
		comp := NewBaseObjectFollow()
		comp.SetStopDistance(-1)
		assert.Equal(t, 0.0, comp.GetStopDistance())
	})
}

func Test_ObjectFollowGetSpeed(t *testing.T) {
	expectedSpeed := 123.0
	comp := NewObjectFollow(expectedSpeed, "123", 10)
	foundSpeed := comp.GetSpeed()
	assert.Equal(t, expectedSpeed, foundSpeed)
}

func Test_ObjectFollowGetObjectID(t *testing.T) {
	expectedObjID := "obj123"
	comp := NewObjectFollow(100.0, expectedObjID, 10)
	foundObjID := comp.GetObjectID()
	assert.Equal(t, expectedObjID, foundObjID)
}

func Test_UnmarshalObjectFollow(t *testing.T) {
	tests := []struct {
		name         string
		baseComp     *ObjectFollow
		data         json.RawMessage
		expectedErr  bool
		expectedComp *ObjectFollow
	}{
		{
			name:     "valid unmarshal with minimal fields",
			baseComp: NewBaseObjectFollow(),
			data: json.RawMessage(
				`{
					"speed": 100
				}`,
			),
			expectedErr: false,
			expectedComp: &ObjectFollow{
				speed: 100.0,
			},
		},
		{
			name:     "valid unmarshal with full payload",
			baseComp: NewBaseObjectFollow(),
			data: json.RawMessage(
				`{
					"speed": 100,
					"objectID": "obj123",
					"stopDistance": 1.0
				}`,
			),
			expectedErr: false,
			expectedComp: &ObjectFollow{
				speed:        100.0,
				objectID:     "obj123",
				stopDistance: 1.0,
			},
		},
		{
			name:     "invalid unmarshal",
			baseComp: NewBaseObjectFollow(),
			data: json.RawMessage(
				`{
					"speed": 100,
					"objectID": "obj123"
					"stopDistance": -1.0
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseObjectFollow(),
		},
		{
			name:     "stop distance negative",
			baseComp: NewBaseObjectFollow(),
			data: json.RawMessage(
				`{
					"speed": 100,
					"objectID": "obj123",
					"stopDistance": -1.0
				}`,
			),
			expectedErr:  true,
			expectedComp: NewBaseObjectFollow(),
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

			assert.Equal(t, tt.expectedComp, tt.baseComp)
		})
	}
}
