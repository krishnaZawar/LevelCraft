package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_CommandFactory(t *testing.T) {
	tests := []struct {
		name          string
		factory       models.CommandFactory
		data          json.RawMessage
		ExpectedErr   bool
		ExpectedValue models.Command
	}{
		{
			name:          "AddGameobjectCommandFactoryTest",
			factory:       NewAddGameobjectCommandFactory(),
			data:          []byte{},
			ExpectedErr:   false,
			ExpectedValue: &AddGameobjectCommand{},
		},
		{
			name:    "DeleteGameobjectCommandFactoryTest valid marshal",
			factory: NewDeleteGameobjectCommandFactory(),
			data: json.RawMessage(`{
				"id": "123"
			}`),
			ExpectedErr:   false,
			ExpectedValue: &DeleteGameobjectCommand{ID: "123"},
		},
		{
			name:    "DeleteGameobjectCommandFactoryTest invalid marshal",
			factory: NewDeleteGameobjectCommandFactory(),
			data: json.RawMessage(`{
				"id": "123"
			`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
		{
			name:    "AttachComponentCommandFactoryTest valid marshal",
			factory: NewAttachComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123",
				"componentName": "ok"
			}`),
			ExpectedErr:   false,
			ExpectedValue: &AttachComponentCommand{ID: "123", ComponentName: "ok"},
		},
		{
			name:    "AttachComponentCommandFactoryTest invalid marshal",
			factory: NewAttachComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123"
			`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
		{
			name:    "DetachComponentCommandFactoryTest valid marshal",
			factory: NewDetachComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123",
				"componentName": "ok"
			}`),
			ExpectedErr: false,
			ExpectedValue: &DetachComponentCommand{
				ID:            "123",
				ComponentName: "ok",
			},
		},
		{
			name:    "DetachComponentCommandFactoryTest invalid marshal",
			factory: NewDetachComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123",
				"componentName": "ok"
			`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
		{
			name:    "UpdateComponentCommandFactoryTest valid marshal",
			factory: NewUpdateComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123",
				"componentName": "ok",
				"data" : {}
			}`),
			ExpectedErr: false,
			ExpectedValue: &UpdateComponentCommand{
				ID:            "123",
				ComponentName: "ok",
				Data:          map[string]interface{}{},
			},
		},
		{
			name:    "UpdateComponentCommandFactoryTest invalid marshal",
			factory: NewUpdateComponentCommandFactory(),
			data: json.RawMessage(`{
				"id": "123",
				"componentName": "ok",
				"data" : {}
			`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
	}

	for _, tt := range tests {
		comm, err := tt.factory.NewCommand(tt.data)
		if tt.ExpectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tt.ExpectedValue, comm)
	}
}
