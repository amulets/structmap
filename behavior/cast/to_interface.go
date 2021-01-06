package cast

import "reflect"

func (c *Cast) toInterface(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	result = value.Interface()

	return
}
