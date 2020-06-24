package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toUint(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Uint:
		result = value.Uint()
	case reflect.Int:
		i := value.Int()
		if i < 0 {
			err = fmt.Errorf("cannot parse, %d overflows uint", i)
			return
		}

		result = uint64(i)
	case reflect.Float32:
		f := value.Float()
		if f < 0 {
			err = fmt.Errorf("cannot parse, %f overflows uint", f)
			return
		}

		result = uint64(f)
	case reflect.Bool:
		if value.Bool() {
			result = uint(1)
		} else {
			result = uint(0)
		}
	case reflect.String:
		var u uint64

		if u, err = strconv.ParseUint(value.String(), 0, source.Bits()); err == nil {
			result = u
		}
	default:
		err = errNoConvertible
	}

	return
}
