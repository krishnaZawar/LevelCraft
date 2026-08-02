package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_Commands(t *testing.T) {
	tests := []struct {
		name                string
		command             models.Command
		expectedCommandName string
		expectedEvents      []models.Event
	}{
		{
			name:                "AddGameobjectCommandTest",
			command:             &AddGameobjectCommand{},
			expectedCommandName: base.Command_AddGameobject,
			expectedEvents:      []models.Event{event.NewAddGameobjectEvent()},
		},
		{
			name:                "DeleteGameobjectCommandTest",
			command:             &DeleteGameobjectCommand{},
			expectedCommandName: base.Command_DeleteGameobject,
			expectedEvents:      []models.Event{event.NewDeleteGameobjectEvent("")},
		},
		{
			name:                "AttachComponentCommandTest",
			command:             &AttachComponentCommand{},
			expectedCommandName: base.Command_AttachComponent,
			expectedEvents:      []models.Event{event.NewAttachComponentEvent("", "")},
		},
		{
			name:                "DetachComponentCommandTest",
			command:             &DetachComponentCommand{},
			expectedCommandName: base.Command_DetachComponent,
			expectedEvents:      []models.Event{event.NewDetachComponentEvent("", "")},
		},
		{
			name: "UpdateComponentCommandTest",
			command: &UpdateComponentCommand{
				Data: map[string]interface{}{},
			},
			expectedCommandName: base.Command_UpdateComponent,
			expectedEvents:      []models.Event{event.NewUpdateComponentEvent("", "", map[string]interface{}{})},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expectedCommandName, tt.command.GetCommandName())
		assert.Equal(t, tt.expectedEvents, tt.command.Handle())
	}
}
