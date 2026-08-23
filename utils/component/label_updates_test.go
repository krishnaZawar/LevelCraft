package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func contains(arr []string, value string) bool {
	for _, val := range arr {
		if val == value {
			return true
		}
	}
	return false
}

func Test_AssertLabelUpdates(t *testing.T) {
	tests := []struct {
		name  string
		arr   []string
		value string
	}{
		{
			name:  "Assert Transform_LabelsX update",
			arr:   transform_LabelsX,
			value: Transform_CurLabelX,
		},
		{
			name:  "Assert Transform_LabelsY update",
			arr:   transform_LabelsY,
			value: Transform_CurLabelY,
		},
		{
			name:  "Assert Transform_LabelsW update",
			arr:   transform_LabelsW,
			value: Transform_CurLabelW,
		},
		{
			name:  "Assert Transform_LabelsH update",
			arr:   transform_LabelsH,
			value: Transform_CurLabelH,
		},
		{
			name:  "Assert Color_LabelsR update",
			arr:   color_LabelsR,
			value: Color_CurLabelR,
		},
		{
			name:  "Assert Color_LabelsG update",
			arr:   color_LabelsG,
			value: Color_CurLabelG,
		},
		{
			name:  "Assert Color_LabelsB update",
			arr:   color_LabelsB,
			value: Color_CurLabelB,
		},
		{
			name:  "Assert Color_LabelsA update",
			arr:   color_LabelsA,
			value: Color_CurLabelA,
		},
	}

	for _, tt := range tests {
		testName := map[string]interface{}{
			"name": tt.name,
		}
		assert.Equal(t, true, contains(tt.arr, tt.value), testName)
	}
}
