package component

// Component is a modular unit of data and functions attached to the object
//
// It models the behaviour of the object in the scene
type Component interface {
	// Returns the name of the component
	GetComponentName() string

	// Returns a snapshot of the complete data stored in the component
	//
	// Helpful in recursively building the game scene
	GetComponentDetails() map[string]interface{}
}
