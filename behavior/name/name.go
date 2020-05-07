package name

import "github.com/dungeon-code/structmap"

// Discovery stop on first discover function, that's change field name
func Discovery(discoveries ...structmap.BehaviorFunc) structmap.BehaviorFunc {
	return func(field *structmap.FieldPart) error {
		currentName := field.Name

		for _, discover := range discoveries {
			if err := discover(field); err != nil {
				return err
			}

			if field.Name != currentName {
				break
			}
		}

		return nil
	}
}

// Noop do not change field name
func Noop(field *structmap.FieldPart) error {
	return nil
}
