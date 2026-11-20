package simulationhandler

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

var compRegistry = component.NewComponentRegistry(&input.InputMapping{})

func Test_ObjectFollowSimulationHandler(t *testing.T) {
	dt := 1.0
	// comp is added explicitly to the follower to link the component to it
	// follower details never hold the objectFollow component here, that is injected
	tests := []struct {
		name               string
		comp               component.Component
		followerID         string
		followerObjDetails gameobject.GameobjectDetails
		targetObjDetails   gameobject.GameobjectDetails
		expectedErr        bool
		expectedEvents     []models.Event
	}{
		{
			name:       "nil component passed",
			comp:       nil,
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
			},
			expectedErr:    true,
			expectedEvents: nil,
		},
		{
			name:       "invalid component passed",
			comp:       component.NewColor(10, 10, 10, 10),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
			},
			expectedErr:    true,
			expectedEvents: nil,
		},
		{
			name:       "follower not found",
			comp:       component.NewObjectFollow(100, "target", 10),
			followerID: "non-existent",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
			},
			expectedErr:    true,
			expectedEvents: nil,
		},
		{
			name:       "target not found",
			comp:       component.NewObjectFollow(100, "non-existent", 10),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
			},
			expectedErr:    false,
			expectedEvents: []models.Event{},
		},
		{
			name:       "follower transform does not exist",
			comp:       component.NewObjectFollow(100, "target", 10),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			expectedErr:    true,
			expectedEvents: nil,
		},
		{
			name:       "target transform does not exist",
			comp:       component.NewObjectFollow(100, "target", 10),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
			},
			expectedErr:    true,
			expectedEvents: nil,
		},
		{
			name:       "distance less than stopDistance",
			comp:       component.NewObjectFollow(100, "target", 200),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 150,
							"y": 150,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			expectedErr:    false,
			expectedEvents: []models.Event{},
		},
		{
			name:       "generate valid events for movement",
			comp:       component.NewObjectFollow(100, "target", 20),
			followerID: "follower",
			followerObjDetails: gameobject.GameobjectDetails{
				Id: "follower",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			targetObjDetails: gameobject.GameobjectDetails{
				Id: "target",
				Components: []component.ComponentDetails{
					{
						Name: compbase.ComponentName_Transform,
						Data: json.RawMessage(`{
							"x": 200,
							"y": 200,
							"w": 100,
							"h": 100
						}`),
					},
				},
			},
			expectedErr: false,
			expectedEvents: []models.Event{
				&event.MoveEvent{
					GameobjectID: "follower",
					MoveX:        100,
					MoveY:        100,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gsm := gamestatemanager.NewGameStateManager()

			follower := gameobject.NewGameobjectWithID(tt.followerObjDetails.Id, compRegistry)
			err := follower.BuildFromDetails(tt.followerObjDetails)
			assert.Nil(t, err)
			_ = follower.AddComponent(tt.comp)

			target := gameobject.NewGameobjectWithID(tt.targetObjDetails.Id, compRegistry)
			err = target.BuildFromDetails(tt.targetObjDetails)
			assert.Nil(t, err)

			gsm.AddGameobject(follower)
			gsm.AddGameobject(target)

			simulator := NewObjectFollowSimulator(gsm)
			simulationCtx := SimulatorCtx{
				GameobjectID: tt.followerID,
				Component:    tt.comp,
				Dt:           dt,
			}

			evts, err := simulator.Simulate(simulationCtx)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.ElementsMatch(t, tt.expectedEvents, evts)
		})
	}
}
