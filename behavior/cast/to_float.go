package cast

import (
	"reflect"
	"strconv"
)

func toFloat(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Float32:
		result = value.Interface()
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
		var f float64

		if f, err = strconv.ParseFloat(value.String(), source.Bits()); err == nil {
			result = f
		}
	default:
		err = errNoConvertible
	}

	return
}
