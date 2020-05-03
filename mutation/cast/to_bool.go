package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dungeon-code/structmap"
)

func toBool(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	switch toKind(fieldValue.Type()) {
	case reflect.Bool:
		// Ignore is a bool type
	case reflect.Int:
		field.Value = fieldValue.Int() != 0
	case reflect.Uint:
		field.Value = fieldValue.Uint() != 0
	case reflect.Float32:
		field.Value = fieldValue.Float() != 0
	case reflect.String:
		b, err := strconv.ParseBool(fieldValue.String())
		if err == nil {
			field.Value = b
		} else if fieldValue.String() == "" {
			field.Value = false
		} else {
			return fmt.Errorf("cannot parse '%s' as bool: %s", field.Name, err)
		}
	default:
		return fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, fieldValue.Type())
	}

	return nil
}
