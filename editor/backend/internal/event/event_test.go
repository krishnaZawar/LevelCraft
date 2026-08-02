package event

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_Event(t *testing.T) {
	tests := []struct {
		name              string
		event             models.Event
		expectedEvent     models.Event
		expectedEventName string
	}{
		{
			name:              "AddGameobjectEventTest",
			event:             NewAddGameobjectEvent(),
			expectedEvent:     &AddGameobjectEvent{},
			expectedEventName: base.Event_AddGameobject,
		},
		{
			name:              "DeleteGameobjectEventTest",
			event:             NewDeleteGameobjectEvent("123"),
			expectedEvent:     &DeleteGameobjectEvent{ID: "123"},
			expectedEventName: base.Event_DeleteGameobject,
		},
		{
			name:              "AttachComponentEventTest",
			event:             NewAttachComponentEvent("123", "ok"),
			expectedEvent:     &AttachComponentEvent{ID: "123", ComponentName: "ok"},
			expectedEventName: base.Event_AttachComponent,
		},
		{
			name:              "DetachComponentEventTest",
			event:             NewDetachComponentEvent("123", "ok"),
			expectedEvent:     &DetachComponentEvent{ID: "123", ComponentName: "ok"},
			expectedEventName: base.Event_DetachComponent,
		},
		{
			name:              "UpdateComponentEventTest",
			event:             NewUpdateComponentEvent("123", "ok", map[string]interface{}{}),
			expectedEvent:     &UpdateComponentEvent{ID: "123", ComponentName: "ok", Data: map[string]interface{}{}},
			expectedEventName: base.Event_UpdateComponent,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expectedEvent, tt.event)
		assert.Equal(t, tt.expectedEventName, tt.event.GetEventName())
	}
}
