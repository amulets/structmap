package cast

import (
	"reflect"
	"time"
)

func toStructTime(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Int:
		result = time.Unix(0, value.Int()*int64(time.Millisecond))
	case reflect.Float32:
		result = time.Unix(0, int64(value.Float())*int64(time.Millisecond))
	case reflect.String:
		result, err = time.Parse(time.RFC3339, value.String())
	default:
		err = errNoConvertible
	}

	return
}
