package name

import "github.com/amulets/structmap"

// FromTag get field name from tagName, if has a name
func FromTag(tagName string) structmap.BehaviorFunc {
	return func(field *structmap.FieldPart) error {
		name, _ := structmap.ParseTag(field.Tag.Get(tagName))
		if name != "" {
			field.Name = name
		}

		return nil
	}
}
