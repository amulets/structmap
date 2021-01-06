package cast

import (
	"fmt"
	"reflect"

	"github.com/amulets/structmap/internal"
)

func (c *Cast) toList(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	valueKind := value.Type().Kind()

	switch valueKind {
	case reflect.Slice:
	case reflect.Array:
	default:
		err = ErrNoConvertible
		return
	}

	// It is a same type
	if value.Kind() == source.Kind() && internal.Type(value.Type().Elem()).Kind() == internal.Type(source.Elem()).Kind() {
		result = value.Interface()
		return
	}

	var (
		listType  reflect.Type
		listValue reflect.Value
	)
	switch source.Kind() {
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
		err = ErrNoConvertible
		return
	}

	for i := 0; i < value.Len(); i++ {
		itemValue := value.Index(i)

		if itemValue, err = c.toType(source.Elem(), itemValue); err != nil {
			switch err {
			case ErrEmptyValue, ErrNoCoveredType:
				err = nil
			default:
				return
			}
		}

		if source.Kind() == reflect.Slice {
			listValue.Set(reflect.Append(listValue, itemValue))
		} else {
			listValue.Index(i).Set(itemValue)
		}
	}

	result = listValue.Interface()

	return
}
