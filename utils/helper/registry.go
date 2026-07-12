package helper

// Registry is used to register the objects once at initialization.
//
// It is a generic registry that can be used to register Handlers, Factories, etc
type Registry[K comparable, V any] struct {
	items map[K]V
}

func NewRegistry[K comparable, V any]() *Registry[K, V] {
	return &Registry[K, V]{
		items: make(map[K]V),
	}
}

// Registers an item with the registry
func (r *Registry[K, V]) Register(key K, val V) {
	if r.items == nil {
		r.items = make(map[K]V)
	}
	r.items[key] = val
}

// Fetches the item associated with the key
//
// Return Values
//	- V : value of the item associated with key (if found)
//	- bool: True if the key was found else False
func (r *Registry[K, V]) GetValue(key K) (V, bool) {
	if r.items == nil {
		r.items = make(map[K]V)
	}
	val, ok := r.items[key]
	return val, ok
}
