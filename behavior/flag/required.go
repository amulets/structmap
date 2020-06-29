package flag

import (
	"fmt"

	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior"
)

// Required checks if the value to be filled in the structure
// is present, otherwise an exception will be thrown.
// Only for fields marked with the flag: required
// Example: `tagName:",required"`
func Required(tagName string) structmap.Behavior {
	return behavior.New(func(field *structmap.FieldPart) (err error) {
		if field.Value != nil {
			return
		}

		_, flags := structmap.ParseTag(field.Tag.Get(tagName))

		if flags.Has("required") {
			err = fmt.Errorf("field %s is required", field.Name)
		}

		return
	})
}
