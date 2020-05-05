package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toBool(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Bool:
		// Ignore is a bool type
	case reflect.Int:
		result = value.Int() != 0
	case reflect.Uint:
		result = value.Uint() != 0
	case reflect.Float32:
		result = value.Float() != 0
	case reflect.String:
		b, err := strconv.ParseBool(value.String())
		if err == nil {
			result = b
		} else if value.String() == "" {
			result = false
		} else {
			err = fmt.Errorf("cannot parse to bool: %s", err)
		}
	default:
		err = errNoConvertible
	}

	return
}
