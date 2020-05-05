package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toInt(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Int:
		// Ignore is a int type
	case reflect.Uint:
		result = int64(value.Uint())
	case reflect.Float32:
		result = int64(value.Float())
	case reflect.Bool:
		if value.Bool() {
			result = 1
		} else {
			result = 0
		}
	case reflect.String:
		sourceType := source

		if source.Kind() == reflect.Ptr {
			sourceType = source.Elem()
		}

		i, err := strconv.ParseInt(value.String(), 0, sourceType.Bits())
		if err == nil {
			result = i
		} else {
			err = fmt.Errorf("cannot parse to int: %s", err)
		}
	default:
		err = errNoConvertible
	}

	return
}
