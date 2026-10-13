package eventgenerator

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_InitEventGeneratorRegistry(t *testing.T) {
	gsm := &gamestatemanager.GameStateManager{}
	registry := initEventGeneratorRegistry(gsm)

	assert.Equal(t, 1, registry.Size())

	var expected EventGenerator
	found, ok := registry.GetValue(compbase.ComponentName_InputMovement)
	expected = NewInputMovementEventGenerator(gsm)
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, found)
}
