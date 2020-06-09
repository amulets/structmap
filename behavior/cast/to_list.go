package cast

import (
	"fmt"
	"reflect"
)

func toList(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	valueKind := toKind(value.Type())

	switch valueKind {
	case reflect.Slice:
	case reflect.Array:
	default:
		err = errNoConvertible
		return
	}

	// It is a same type
	if value.Type().Elem().Kind() == source.Elem().Kind() && value.Kind() == source.Kind() {
		return
	}

	var (
		listType   reflect.Type
		listValue  reflect.Value
		sourceKind = toKind(source)
	)
	switch sourceKind {
	case reflect.Slice:
		listType = reflect.SliceOf(source.Elem())
		sliceValue := reflect.MakeSlice(listType, 0, 0)

		listValue = reflect.Indirect(reflect.New(listType))
		listValue.Set(sliceValue)
	case reflect.Array:
		if value.Len() > source.Len() {
			err = fmt.Errorf("[%d]%s is greater than fixed [%d]array", value.Len(), valueKind, source.Len())
			return
		}

		listType = reflect.ArrayOf(source.Len(), source.Elem())
		listValue = reflect.Indirect(reflect.New(listType))
	default:
		err = errNoConvertible
		return
	}

	var (
		item      interface{}
		itemValue reflect.Value
	)
	for i := 0; i < value.Len(); i++ {
		itemValue = value.Index(i)

		item, err = toType(source.Elem(), itemValue)
		if err == errNoCoveredType {
			err = nil
		}

		if err != nil {
			return
		}

		if item != nil {
			itemValue = reflect.ValueOf(item)
		}

		// If do not is a slice, convert its type ex: int32 to int
		if itemValue.Kind() != reflect.Slice {
			if !itemValue.Type().ConvertibleTo(source.Elem()) {
				err = fmt.Errorf("slice item %v cannot is convertible to %s", itemValue.Interface(), source.Elem().Kind())
				return
			}

			itemValue = itemValue.Convert(source.Elem())
		}

		if sourceKind == reflect.Slice {
			listValue.Set(reflect.Append(listValue, itemValue))
		} else {
			listValue.Index(i).Set(itemValue)
		}
	}

	result = listValue.Interface()

	return
}
