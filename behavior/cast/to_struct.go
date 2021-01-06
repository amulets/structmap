package cast

import (
	"reflect"
)

func (c *Cast) toStruct(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	// Is the same type
	if source == value.Type() {
		result = value.Interface()
		return
	}

	convert, ok := c.convertToType[source]
	if !ok {
		err = ErrNoCoveredType
		return
	}

	return convert(source, value)
}
