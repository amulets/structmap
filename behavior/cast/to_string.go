package cast

import (
	"reflect"
	"strconv"
)

func toString(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.String:
		result = value.String()
	case reflect.Bool:
		if value.Bool() {
			result = "true"
		} else {
			result = "false"
		}
	case reflect.Int:
		result = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		result = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		result = strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits())
	default:
		err = errNoConvertible
	}

	return
}
