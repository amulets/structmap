package name

import (
	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior"
)

// Discovery stop on first discover function, that's change field name
func Discovery(discoveries ...structmap.Behavior) structmap.Behavior {
	return behavior.New(func(field *structmap.FieldPart) error {
		currentName := field.Name

		for _, discover := range discoveries {
			if err := discover.Do(field); err != nil {
				return err
			}

			if field.Name != currentName {
				break
			}
		}

		return nil
	})
}

// Noop do not change field name
var Noop = behavior.New(func(field *structmap.FieldPart) error {
	return nil
})
