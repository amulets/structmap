package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toFloat(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Float32:
		// Ignore is a float type
	case reflect.Int:
		result = float64(value.Int())
	case reflect.Uint:
		result = float64(value.Uint())
	case reflect.Bool:
		if value.Bool() {
			result = float32(1)
		} else {
			result = float32(0)
		}
	case reflect.String:
		sourceType := source

		if source.Kind() == reflect.Ptr {
			sourceType = source.Elem()
		}

		f, err := strconv.ParseFloat(value.String(), sourceType.Bits())
		if err == nil {
			result = f
		} else {
			err = fmt.Errorf("cannot parse to float: %s", err)
		}
	default:
		err = errNoConvertible
	}

	return
}
