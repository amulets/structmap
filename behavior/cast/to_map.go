package cast

import (
	"reflect"

	"github.com/amulets/structmap/internal"
)

func (c *Cast) toMap(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	if value.Type().Kind() != reflect.Map {
		err = ErrNoConvertible
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

		if mapKeyElem, err = c.toType(source.Key(), mapKeyElem); err != nil {
			switch err {
			case ErrEmptyValue, ErrNoCoveredType:
				err = nil
			default:
				return
			}
		}

		if mapValueElem, err = c.toType(source.Elem(), mapValueElem); err != nil {
			switch err {
			case ErrEmptyValue, ErrNoCoveredType:
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
