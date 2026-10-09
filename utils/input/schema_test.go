package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KeyExists(t *testing.T) {
	mapping := &InputMapping{
		Keyboard: []InputData{
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
		Mouse: []InputData{
			{
				Name:  "LEFT_MOUSE_BUTTON",
				Code:  1,
				Label: "left mouse button",
			},
		},
	}

	tests := []struct {
		name        string
		checkLabel  string
		expectedErr bool
	}{
		{
			name:        "label exists in keyboard",
			checkLabel:  "Key_A",
			expectedErr: false,
		},
		{
			name:        "label exists in mouse",
			checkLabel:  "LEFT_MOUSE_BUTTON",
			expectedErr: false,
		},
		{
			name:        "label does not exist",
			checkLabel:  "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapping.KeyExists(tt.checkLabel)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
