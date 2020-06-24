package cast

import "reflect"

func toInterface(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	result = value.Interface()

	return
}
