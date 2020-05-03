package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dungeon-code/structmap"
)

func toString(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	converted := true
	switch kind := toKind(fieldValue.Type()); kind {
	case reflect.String:
		// Ignore is a string type
	case reflect.Bool:
		if fieldValue.Bool() {
			field.Value = "1"
		} else {
			field.Value = "0"
		}
	case reflect.Int:
		field.Value = strconv.FormatInt(fieldValue.Int(), 10)
	case reflect.Uint:
		field.Value = strconv.FormatUint(fieldValue.Uint(), 10)
	case reflect.Float32:
		field.Value = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
	case reflect.Slice, reflect.Array:
		fieldValueType := fieldValue.Type()
		elemKind := fieldValueType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if kind == reflect.Array {
				uints = make([]uint8, fieldValue.Len(), fieldValue.Len())
				for i := range uints {
					uints[i] = fieldValue.Index(i).Interface().(uint8)
				}
			} else {
				uints = fieldValue.Interface().([]uint8)
			}

			field.Value = string(uints)
		default:
			converted = false
		}
	default:
		converted = false
	}

	if !converted {
		return fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, fieldValue.Type())
	}

	return nil
}
