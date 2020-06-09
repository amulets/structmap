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
		if result, err = strconv.ParseBool(value.String()); err != nil {
			if value.String() == "" {
				err = nil
				result = false
			} else {
				err = fmt.Errorf("cannot parse to bool: %s", err)
			}
		}
	default:
		err = errNoConvertible
	}

	return
}
