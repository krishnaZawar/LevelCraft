package eventhandler

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

type mockEvent struct {
	mockGetEventName func() string
}

func (me *mockEvent) GetEventName() string {
	return "MockEvent"
}

var compRegistry = component.NewComponentRegistry(&input.InputMapping{
	Keyboard: []input.InputData{
		{
			Name:  "KEY_A",
			Code:  1,
			Label: "A",
		},
		{
			Name:  "KEY_D",
			Code:  2,
			Label: "D",
		},
		{
			Name:  "KEY_S",
			Code:  3,
			Label: "S",
		},
		{
			Name:  "KEY_W",
			Code:  4,
			Label: "W",
		},
	},
})

func Test_MoveEventHandler(t *testing.T) {
	t.Run("incompatible event type passed", func(t *testing.T) {
		handler := NewMoveEventHandler(&gamestatemanager.GameStateManager{})
		evt := &mockEvent{}

		resp := handler.Handle(evt, 1)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.MsgExpectedEventTypeNotFound, resp.Msg)
		assert.Equal(t, base.NewErrEventTypeMismatch(base.Event_MoveEvent, evt.GetEventName()), resp.Data)
	})

	t.Run("gameobject does not exist", func(t *testing.T) {
		handler := NewMoveEventHandler(&gamestatemanager.GameStateManager{})

		objID := "1234"
		dispX, dispY, dt := 12.0, 46.0, 1.0
		evt := &event.MoveEvent{
			GameobjectID: objID,
			MoveX:        dispX,
			MoveY:        dispY,
		}

		resp := handler.Handle(evt, dt)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.MsgGameobjectNotFound, resp.Msg)
		assert.Equal(t, base.NewErrGameobjectNotFound(objID), resp.Data)
	})

	t.Run("transform component on gameobject does not exist", func(t *testing.T) {
		gsm := gamestatemanager.NewGameStateManager()

		obj := gameobject.NewGameobject(compRegistry)
		gsm.AddGameobject(obj)

		handler := NewMoveEventHandler(gsm)

		dispX, dispY, dt := 12.0, 46.0, 1.0
		evt := &event.MoveEvent{
			GameobjectID: obj.GetID(),
			MoveX:        dispX,
			MoveY:        dispY,
		}

		resp := handler.Handle(evt, dt)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.MsgComponentNotFound, resp.Msg)
		assert.Equal(t, base.NewErrComponentNotFound(compbase.ComponentName_Transform, obj.GetID()), resp.Data)
	})

	t.Run("valid event and update", func(t *testing.T) {
		gsm := gamestatemanager.NewGameStateManager()

		x, y, w, h := 100.0, 100.0, 100.0, 100.0

		obj := gameobject.NewGameobject(compRegistry)
		obj.AddComponent(component.NewTransform(x, y, w, h))
		gsm.AddGameobject(obj)

		handler := NewMoveEventHandler(gsm)

		dispX, dispY, dt := 12.0, 46.0, 1.0
		evt := &event.MoveEvent{
			GameobjectID: obj.GetID(),
			MoveX:        dispX,
			MoveY:        dispY,
		}

		resp := handler.Handle(evt, dt)

		assert.Equal(t, true, resp.Success)
		assert.Nil(t, resp.Data)

		comp, found := obj.GetComponent(compbase.ComponentName_Transform)
		assert.Equal(t, true, found)

		transformComp, ok := comp.(*component.Transform)
		assert.Equal(t, true, ok)

		newX, newY := transformComp.GetPosition()

		assert.Equal(t, x+dispX*dt, newX)
		assert.Equal(t, y+dispY*dt, newY)
	})
}
