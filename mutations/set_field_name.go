package mutations

import "github.com/dungeon-code/structmap"

// SetFieldName find field name and define
func SetFieldName(tagName string) structmap.MutationFunc {
	return func(field *structmap.FieldPart) error {
		name, _ := structmap.ParseTag(field.Tag.Get(tagName))
		if name != "" {
			field.Name = name
		}

		return nil
	}
}
