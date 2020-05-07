package rule

import (
	"fmt"

	"github.com/dungeon-code/structmap"
)

// Required checks if the value to be filled in the structure
// is present, otherwise an exception will be thrown.
// Only for fields marked with the flag: required
// Example: `tagName:",required"`
func Required(tagName string) structmap.MutationFunc {
	return func(field *structmap.FieldPart) (err error) {
		if field.Value != nil {
			return
		}

		_, flags := structmap.ParseTag(field.Tag.Get(tagName))

		if flags.Has("required") {
			err = fmt.Errorf("field %s is required", field.Name)
		}

		return
	}
}
