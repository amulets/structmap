package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dungeon-code/structmap"
)

func toInt(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	fieldValueType := fieldValue.Type()

	switch toKind(fieldValueType) {
	case reflect.Int:
		// Ignore is a int type
	case reflect.Uint:
		field.Value = int64(fieldValue.Uint())
	case reflect.Float32:
		field.Value = int64(fieldValue.Float())
	case reflect.Bool:
		if fieldValue.Bool() {
			field.Value = 1
		} else {
			field.Value = 0
		}
	case reflect.String:
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		i, err := strconv.ParseInt(fieldValue.String(), 0, fieldType.Bits())
		if err == nil {
			field.Value = i
		} else {
			return fmt.Errorf("cannot parse '%s' as int: %s", field.Name, err)
		}
	default:
		return fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, fieldValue.Type())
	}

	return nil
}
