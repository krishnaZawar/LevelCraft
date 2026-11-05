package collision

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_IsImmovable(t *testing.T) {
	collider := NewBaseCollider()
	assert.Equal(t, false, collider.IsImmovable())
}

func Test_MarshalAndUnmarshalCollider(t *testing.T) {
	expected := NewCollider(true)

	found := NewBaseCollider()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_Collider, details.Name)
	assert.Equal(t, expected, found)
}

func Test_UnmarshalCollider(t *testing.T) {
	tests := []struct {
		name             string
		baseComp         *Collider
		data             json.RawMessage
		expectedErr      bool
		expectedCollider *Collider
	}{
		{
			name:     "valid unmarshal",
			baseComp: NewBaseCollider(),
			data: json.RawMessage(
				`{
					"immovable": false
				}`,
			),
			expectedErr:      false,
			expectedCollider: NewBaseCollider(),
		},
		{
			name:     "invalid unmarshal",
			baseComp: NewBaseCollider(),
			data: json.RawMessage(
				`{
					"immovable": 1
				}`,
			),
			expectedErr:      true,
			expectedCollider: NewBaseCollider(),
		},
		{
			name:     "component is immovable",
			baseComp: NewBaseCollider(),
			data: json.RawMessage(
				`{
					"immovable": true
				}`,
			),
			expectedErr: false,
			expectedCollider: &Collider{
				immovable: true,
			},
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
			assert.Equal(t, tt.expectedCollider.IsImmovable(), tt.baseComp.IsImmovable())
		})
	}
}
