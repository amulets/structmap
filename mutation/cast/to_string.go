package cast

import (
	"reflect"
	"strconv"
)

func toString(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch kind := toKind(value.Type()); kind {
	case reflect.String:
		// Ignore is a string type
	case reflect.Bool:
		if value.Bool() {
			result = "1"
		} else {
			result = "0"
		}
	case reflect.Int:
		result = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		result = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		result = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Slice, reflect.Array:
		valueType := value.Type()
		elemKind := valueType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if kind == reflect.Array {
				uints = make([]uint8, value.Len(), value.Len())
				for i := range uints {
					uints[i] = value.Index(i).Interface().(uint8)
				}
			} else {
				uints = value.Interface().([]uint8)
			}

			result = string(uints)
		default:
			err = errNoConvertible
		}
	default:
		err = errNoConvertible
	}

	return
}
