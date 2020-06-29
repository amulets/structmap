package behavior

import "github.com/amulets/structmap"

// Func that deals with the behavior of the field
type Func func(*structmap.FieldPart) error

// Do behavior
func (fn Func) Do(field *structmap.FieldPart) error {
	return fn(field)
}

// New instance of BehaviorFunc
func New(behavior Func) structmap.Behavior {
	return behavior
}
