package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dungeon-code/structmap"
)

func toUint(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	switch toKind(fieldValue.Type()) {
	case reflect.Uint:
		// Ignore is a uint type
	case reflect.Int:
		i := fieldValue.Int()
		if i < 0 {
			return fmt.Errorf("cannot parse '%s', %d overflows uint", field.Name, i)
		}

		field.Value = uint64(i)
	case reflect.Float32:
		f := fieldValue.Float()
		if f < 0 {
			return fmt.Errorf("cannot parse '%s', %f overflows uint", field.Name, f)
		}

		field.Value = uint64(f)
	case reflect.Bool:
		if fieldValue.Bool() {
			field.Value = uint(1)
		} else {
			field.Value = uint(0)
		}
	case reflect.String:
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		i, err := strconv.ParseUint(fieldValue.String(), 0, fieldType.Bits())
		if err != nil {
			return fmt.Errorf("cannot parse '%s' as uint: %s", field.Name, err)
		}

		field.Value = i
	default:
		return fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, fieldValue.Type())
	}

	return nil
}
