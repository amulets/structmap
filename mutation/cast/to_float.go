package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dungeon-code/structmap"
)

func toFloat(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	switch toKind(fieldValue.Type()) {
	case reflect.Float32:
		// Ignore is a float type
	case reflect.Int:
		field.Value = float64(fieldValue.Int())
	case reflect.Uint:
		field.Value = float64(fieldValue.Uint())
	case reflect.Bool:
		if fieldValue.Bool() {
			field.Value = float32(1)
		} else {
			field.Value = float32(0)
		}
	case reflect.String:
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		f, err := strconv.ParseFloat(fieldValue.String(), fieldType.Bits())
		if err == nil {
			field.Value = f
		} else {
			return fmt.Errorf("cannot parse '%s' as float: %s", field.Name, err)
		}
	default:
		return fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, fieldValue.Type())
	}

	return nil
}
