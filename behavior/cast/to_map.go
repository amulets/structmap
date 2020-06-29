package cast

import (
	"reflect"

	"github.com/amulets/structmap/internal"
)

func toMap(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	if toKind(value.Type()) != reflect.Map {
		err = errNoConvertible
		return
	}

	// It is a same type
	if internal.Type(value.Type().Key()).Kind() == internal.Type(source.Key()).Kind() &&
		internal.Type(value.Type().Elem()).Kind() == internal.Type(source.Elem()).Kind() {
		result = value.Interface()
		return
	}

	mapType := reflect.MapOf(source.Key(), source.Elem())
	mapValue := reflect.MakeMap(mapType)

	mapIterator := value.MapRange()
	for mapIterator.Next() {
		mapKeyElem := mapIterator.Key()
		mapValueElem := mapIterator.Value()

		if mapKeyElem, err = toType(source.Key(), mapKeyElem); err != nil {
			switch err {
			case errEmptyValue, errNoCoveredType:
				err = nil
			default:
				return
			}
		}

		if mapValueElem, err = toType(source.Elem(), mapValueElem); err != nil {
			switch err {
			case errEmptyValue, errNoCoveredType:
				err = nil
			default:
				return
			}
		}

		mapValue.SetMapIndex(mapKeyElem, mapValueElem)
	}

	result = mapValue.Interface()

	return
}
