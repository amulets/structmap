package cast

import (
	"reflect"
	"strconv"
)

func (c *Cast) toInt(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch ToKind(value.Type()) {
	case reflect.Int:
		result = value.Int()
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
		var i int64

		if i, err = strconv.ParseInt(value.String(), 0, source.Bits()); err == nil {
			result = i
		}
	default:
		if convert, ok := c.convertToType[source]; ok {
			return convert(source, value)
		}

		err = ErrNoConvertible
	}

	return
}
