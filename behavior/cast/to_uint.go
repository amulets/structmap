package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toUint(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch toKind(value.Type()) {
	case reflect.Uint:
		// Ignore is a uint type
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
		sourceType := source

		if source.Kind() == reflect.Ptr {
			sourceType = source.Elem()
		}

		if result, err = strconv.ParseUint(value.String(), 0, sourceType.Bits()); err != nil {
			err = fmt.Errorf("cannot parse to uint: %s", err)
		}
	default:
		err = errNoConvertible
	}

	return
}
