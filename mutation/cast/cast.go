package cast

import (
	"reflect"

	"github.com/dungeon-code/structmap"
)

// Covered types to convert
var convert = map[reflect.Kind]func(*structmap.FieldPart) error{
	reflect.String: toString,
	reflect.Int:    toInt,
}

// ToType cast value to field type value
func ToType(field *structmap.FieldPart) error {
	kind := toKind(field.Type)

	convertTo, ok := convert[kind]
	if !ok {
		// Do not have a coverter to this type
		return nil
	}

	return convertTo(field)
}
