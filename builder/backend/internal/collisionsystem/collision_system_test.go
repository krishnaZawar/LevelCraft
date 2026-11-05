package collisionsystem

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

type mockEvent struct {
	mockGetEventName func() string
}

func (me *mockEvent) GetEventName() string {
	return me.mockGetEventName()
}

func Test_CalculateOverlap(t *testing.T) {
	tests := []struct {
		name                   string
		primaryObj_Transform   *component.Transform
		secondaryObj_Transform *component.Transform
		expectedOverlap        collisionOverlap
	}{
		{
			name:                   "complete overlap",
			primaryObj_Transform:   component.NewTransform(100, 100, 100, 100),
			secondaryObj_Transform: component.NewTransform(100, 100, 100, 100),
			expectedOverlap: collisionOverlap{
				Left:   100,
				Right:  100,
				Top:    100,
				Bottom: 100,
			},
		},
		{
			name:                   "no overlap",
			primaryObj_Transform:   component.NewTransform(0, 0, 10, 10),
			secondaryObj_Transform: component.NewTransform(100, 100, 10, 10),
			expectedOverlap: collisionOverlap{
				Left:   -90,
				Right:  110,
				Top:    -90,
				Bottom: 110,
			},
		},
		{
			name:                   "partial overlap along x axis",
			primaryObj_Transform:   component.NewTransform(100, 100, 100, 100),
			secondaryObj_Transform: component.NewTransform(150, 100, 100, 100),
			expectedOverlap: collisionOverlap{
				Left:   50,
				Right:  150,
				Top:    100,
				Bottom: 100,
			},
		},
		{
			name:                   "partial overlap along y axis",
			primaryObj_Transform:   component.NewTransform(100, 100, 100, 100),
			secondaryObj_Transform: component.NewTransform(100, 150, 100, 100),
			expectedOverlap: collisionOverlap{
				Left:   100,
				Right:  100,
				Top:    50,
				Bottom: 150,
			},
		},
		{
			name:                   "partial overlap along both axes",
			primaryObj_Transform:   component.NewTransform(100, 100, 100, 100),
			secondaryObj_Transform: component.NewTransform(180, 180, 100, 100),
			expectedOverlap: collisionOverlap{
				Left:   20,
				Right:  180,
				Top:    20,
				Bottom: 180,
			},
		},
		{
			name:                   "a border of both the transforms touch",
			primaryObj_Transform:   component.NewTransform(100, 100, 100, 100),
			secondaryObj_Transform: component.NewTransform(200, 100, 100, 100),
			expectedOverlap: collisionOverlap{
				Left:   0,
				Right:  200,
				Top:    100,
				Bottom: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlap := calculateOverlap(tt.primaryObj_Transform, tt.secondaryObj_Transform)
			assert.Equal(t, tt.expectedOverlap, overlap)
		})
	}
}

func Test_DetectCollision(t *testing.T) {
	tests := []struct {
		name string

		obj1_transform *component.Transform
		obj1_collider  *component.Collider

		obj2_transform *component.Transform
		obj2_collider  *component.Collider

		collisionOccured bool
	}{
		{
			name:             "complete overlap",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: true,
		},
		{
			name:             "no overlap",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(300, 300, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: false,
		},
		{
			name:             "partial overlap along x axis",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(150, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: true,
		},
		{
			name:             "partial overlap along y axis",
			obj1_transform:   component.NewTransform(100, 150, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: true,
		},
		{
			name:             "partial overlap along both axis",
			obj1_transform:   component.NewTransform(20, 20, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: true,
		},
		{
			name:             "a border of both the objects touches",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(200, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: false,
		},
		{
			name:             "obj1 transform does not exist",
			obj1_transform:   nil,
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: false,
		},
		{
			name:             "obj1 collider does not exist",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    nil,
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    component.NewCollider(false),
			collisionOccured: false,
		},
		{
			name:             "obj2 transform does not exist",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   nil,
			obj2_collider:    component.NewCollider(false),
			collisionOccured: false,
		},
		{
			name:             "obj2 collider does not exist",
			obj1_transform:   component.NewTransform(100, 100, 100, 100),
			obj1_collider:    component.NewCollider(false),
			obj2_transform:   component.NewTransform(100, 100, 100, 100),
			obj2_collider:    nil,
			collisionOccured: false,
		},
	}

	compRegistry := component.NewComponentRegistry(&input.InputMapping{})
	collisionSystem := &CollisionSystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj1 := gameobject.NewGameobject(compRegistry)
			_ = obj1.AddComponent(tt.obj1_transform)
			_ = obj1.AddComponent(tt.obj1_collider)

			obj2 := gameobject.NewGameobject(compRegistry)
			_ = obj2.AddComponent(tt.obj2_transform)
			_ = obj2.AddComponent(tt.obj2_collider)

			collisionOccured := collisionSystem.DetectCollision(obj1, obj2)
			assert.Equal(t, tt.collisionOccured, collisionOccured)
		})
	}
}

func Test_CollisionResolution(t *testing.T) {
	test_dt := 1.0
	primaryObjID := "primary"
	secondaryObjID := "secondary"
	compRegistry := component.NewComponentRegistry(&input.InputMapping{})
	t.Run("wrong event passed", func(t *testing.T) {
		gsm := gamestatemanager.NewGameStateManager()

		testEvt := &mockEvent{
			mockGetEventName: func() string {
				return "testEvent"
			},
		}

		collisionSystem := NewCollisionSystem(gsm)
		resp := collisionSystem.Handle(testEvt, test_dt)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.NewErrEventTypeMismatch(base.Event_CollisionOccuredEvent, testEvt.GetEventName()), resp.Data)
	})

	t.Run("both objects present in scene", func(t *testing.T) {
		/*
			Covers:
				- no collision
				- both immovable
				- left, right, top, bottom resolution for cases:
					- both movable
					- primary movable
					- secondary movable
				- collider not present
					- primary, secondary object
				- transform not present
					- primary, secondary object
		*/
		tests := []struct {
			name string

			initialPrimaryTransform *component.Transform
			primaryCollider         *component.Collider

			initialSecondaryTransform *component.Transform
			secondaryCollider         *component.Collider

			expectedPrimaryTransform   *component.Transform
			expectedSecondaryTransform *component.Transform

			expectedResponse *models.EventResponse
		}{

			{
				name:                       "No collision between objects",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(300, 300, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(300, 300, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "collision between immovable objects",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(true),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(true),
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(100, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "left side collision between both movable objects",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(190, 180, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(95, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(195, 180, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "right side collision between both movable objects",
				initialPrimaryTransform:    component.NewTransform(190, 180, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(195, 180, 100, 100),
				expectedSecondaryTransform: component.NewTransform(95, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "top side collision between both movable objects",
				initialPrimaryTransform:    component.NewTransform(100, 20, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 15, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 115, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "bottom side collision between both movable objects",
				initialPrimaryTransform:    component.NewTransform(100, 200, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 205, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 105, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "left side collision between objects where only primary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(190, 180, 100, 100),
				secondaryCollider:          component.NewCollider(true),
				expectedPrimaryTransform:   component.NewTransform(90, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(190, 180, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "right side collision between objects where only primary is movable",
				initialPrimaryTransform:    component.NewTransform(190, 180, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(true),
				expectedPrimaryTransform:   component.NewTransform(200, 180, 100, 100),
				expectedSecondaryTransform: component.NewTransform(100, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "top side collision between objects where only primary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 20, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(true),
				expectedPrimaryTransform:   component.NewTransform(100, 10, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 110, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "bottom side collision between objects where only primary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 200, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(true),
				expectedPrimaryTransform:   component.NewTransform(100, 210, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 110, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "left side collision between objects where only secondary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(true),
				initialSecondaryTransform:  component.NewTransform(190, 180, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(200, 180, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "right side collision between objects where only secondary is movable",
				initialPrimaryTransform:    component.NewTransform(190, 180, 100, 100),
				primaryCollider:            component.NewCollider(true),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(190, 180, 100, 100),
				expectedSecondaryTransform: component.NewTransform(90, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "top side collision between objects where only secondary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 20, 100, 100),
				primaryCollider:            component.NewCollider(true),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 20, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 120, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "bottom side collision between objects where only secondary is movable",
				initialPrimaryTransform:    component.NewTransform(100, 200, 100, 100),
				primaryCollider:            component.NewCollider(true),
				initialSecondaryTransform:  component.NewTransform(150, 110, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 200, 100, 100),
				expectedSecondaryTransform: component.NewTransform(150, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: true,
					Data:    nil,
				},
			},
			{
				name:                       "collider is not present on primary object",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            nil,
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(100, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: false,
					Data:    base.NewErrComponentNotFound(compbase.ComponentName_Collider, primaryObjID),
				},
			},
			{
				name:                       "collider is not present on secondary object",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          nil,
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: component.NewTransform(100, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: false,
					Data:    base.NewErrComponentNotFound(compbase.ComponentName_Collider, secondaryObjID),
				},
			},
			{
				name:                       "transform is not present on primary object",
				initialPrimaryTransform:    nil,
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  component.NewTransform(100, 100, 100, 100),
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   nil,
				expectedSecondaryTransform: component.NewTransform(100, 100, 100, 100),
				expectedResponse: &models.EventResponse{
					Success: false,
					Data:    base.NewErrComponentNotFound(compbase.ComponentName_Transform, primaryObjID),
				},
			},
			{
				name:                       "transform is not present on secondary object",
				initialPrimaryTransform:    component.NewTransform(100, 100, 100, 100),
				primaryCollider:            component.NewCollider(false),
				initialSecondaryTransform:  nil,
				secondaryCollider:          component.NewCollider(false),
				expectedPrimaryTransform:   component.NewTransform(100, 100, 100, 100),
				expectedSecondaryTransform: nil,
				expectedResponse: &models.EventResponse{
					Success: false,
					Data:    base.NewErrComponentNotFound(compbase.ComponentName_Transform, secondaryObjID),
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gsm := gamestatemanager.NewGameStateManager()
				collisionSystem := NewCollisionSystem(gsm)

				primaryObj := gameobject.NewGameobjectWithID(primaryObjID, compRegistry)
				_ = primaryObj.AddComponent(tt.initialPrimaryTransform)
				_ = primaryObj.AddComponent(tt.primaryCollider)

				secondaryObj := gameobject.NewGameobjectWithID(secondaryObjID, compRegistry)
				_ = secondaryObj.AddComponent(tt.initialSecondaryTransform)
				_ = secondaryObj.AddComponent(tt.secondaryCollider)

				gsm.AddGameobject(primaryObj)
				gsm.AddGameobject(secondaryObj)

				collisionEvt := &event.CollisionOccuredEvent{
					Obj1_ID: primaryObjID,
					Obj2_ID: secondaryObjID,
				}

				resp := collisionSystem.Handle(collisionEvt, test_dt)
				assert.Equal(t, tt.expectedResponse.Success, resp.Success)
				assert.Equal(t, tt.expectedResponse.Data, resp.Data)

				assert.Equal(t, tt.expectedPrimaryTransform, tt.initialPrimaryTransform)
				assert.Equal(t, tt.expectedSecondaryTransform, tt.initialSecondaryTransform)
			})
		}
	})

	t.Run("primary object does not exist", func(t *testing.T) {
		gsm := gamestatemanager.NewGameStateManager()

		secondaryObj := gameobject.NewGameobjectWithID(secondaryObjID, compRegistry)

		_ = secondaryObj.AddComponent(component.NewTransform(100, 100, 100, 100))
		_ = secondaryObj.AddComponent(component.NewCollider(false))

		gsm.AddGameobject(secondaryObj)

		collisionSystem := NewCollisionSystem(gsm)

		collisionEvt := &event.CollisionOccuredEvent{
			Obj1_ID: primaryObjID,
			Obj2_ID: secondaryObjID,
		}

		resp := collisionSystem.Handle(collisionEvt, test_dt)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.NewErrGameobjectNotFound(primaryObjID), resp.Data)
	})

	t.Run("secondary object does not exist", func(t *testing.T) {
		gsm := gamestatemanager.NewGameStateManager()

		primaryObj := gameobject.NewGameobjectWithID(primaryObjID, compRegistry)

		_ = primaryObj.AddComponent(component.NewTransform(100, 100, 100, 100))
		_ = primaryObj.AddComponent(component.NewCollider(false))

		gsm.AddGameobject(primaryObj)

		collisionSystem := NewCollisionSystem(gsm)

		collisionEvt := &event.CollisionOccuredEvent{
			Obj1_ID: primaryObjID,
			Obj2_ID: secondaryObjID,
		}

		resp := collisionSystem.Handle(collisionEvt, test_dt)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.NewErrGameobjectNotFound(secondaryObjID), resp.Data)
	})
}
