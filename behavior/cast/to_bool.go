package cast

import (
	"reflect"
	"strconv"
)

func toBool(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Bool:
		result = value.Bool()
	case reflect.Int:
		result = value.Int() != 0
	case reflect.Uint:
		result = value.Uint() != 0
	case reflect.Float32:
		result = value.Float() != 0
	case reflect.String:
		var b bool

		if b, err = strconv.ParseBool(value.String()); err == nil {
			result = b
		} else if value.String() == "" {
			err = nil
			result = false
		}
	default:
		err = errNoConvertible
	}

	return
}
