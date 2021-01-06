package cast

import (
	"reflect"
	"strconv"
)

func (c *Cast) toBool(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	switch ToKind(value.Type()) {
	case reflect.Bool:
		result = value.Bool()
	case reflect.Int:
		result = value.Int() != 0
	case reflect.Uint:
		result = value.Uint() != 0
	case reflect.Float32:
		result = value.Float() != 0
	case reflect.String:
		var b bool

		if b, err = strconv.ParseBool(value.String()); err == nil {
			result = b
		} else if value.String() == "" {
			err = nil
			result = false
		}
	default:
		if convert, ok := c.convertToType[source]; ok {
			return convert(source, value)
		}

		err = ErrNoConvertible
	}

	return
}
