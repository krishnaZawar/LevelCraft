package collision

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

// Collider is used to enable collision detection and resolution on the object.
type Collider struct {
	immovable bool // states whether positional correction should occur on this object or not during collision resolution
}

// intermediary structure of Collider used for marshalling and unmarshalling component details
type colliderJSON struct {
	Immovable bool `json:"immovable"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBaseCollider() *Collider {
	return &Collider{}
}

func NewCollider(immovable bool) *Collider {
	return &Collider{
		immovable: immovable,
	}
}

// Tells whether the object should move or not during collision resolution
func (c *Collider) IsImmovable() bool {
	return c.immovable
}

// Returns the name of the component
func (c *Collider) GetComponentName() string {
	return base.ComponentName_Collider
}

// Returns a snapshot of the complete data stored in the component
func (c *Collider) GetComponentDetails() models.ComponentDetails {
	data := colliderJSON{
		Immovable: c.immovable,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
		Name: c.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (c *Collider) BuildFromDetails(data json.RawMessage) error {
	var componentData colliderJSON
	err := json.Unmarshal(data, &componentData)
	if err != nil {
		return err
	}

	c.immovable = componentData.Immovable

	return nil
}

// no interaction possible
func (c *Collider) HandlesInput(input string) bool {
	return false
}

var _ models.Component = &Collider{}
