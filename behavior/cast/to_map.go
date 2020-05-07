package cast

import (
	"fmt"
	"reflect"
)

func toMap(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	if toKind(value.Type()) != reflect.Map {
		err = errNoConvertible
		return
	}

	// It is a same type
	if value.Type().Key().Kind() == source.Key().Kind() && value.Type().Elem().Kind() == source.Elem().Kind() {
		return
	}

	mapType := reflect.MapOf(source.Key(), source.Elem())
	mapValue := reflect.MakeMap(mapType)

	var (
		mapKey  interface{}
		mapElem interface{}
	)
	mapIterator := value.MapRange()
	for mapIterator.Next() {
		mapKeyElem := mapIterator.Key()
		switch mapKeyElem.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			mapKeyElem = mapKeyElem.Elem()
		}
		mapValueElem := mapIterator.Value()
		switch mapValueElem.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			mapValueElem = mapValueElem.Elem()
		}

		mapKey, err = toType(source.Key(), mapKeyElem)
		if err == errNoCoveredType {
			err = nil
		}

		if err != nil {
			return
		}

		if mapKey != nil {
			mapKeyElem = reflect.ValueOf(mapKey)
		}

		mapElem, err = toType(source.Elem(), mapValueElem)
		if err == errNoCoveredType {
			err = nil
		}

		if err != nil {
			return
		}

		if mapElem != nil {
			mapValueElem = reflect.ValueOf(mapElem)
		}

		// If do not is a map, convert its type ex: int32 to int
		if mapValueElem.Kind() != reflect.Map {
			if !mapKeyElem.Type().ConvertibleTo(source.Key()) {
				err = fmt.Errorf("map key %v cannot is convertible to %s", mapKeyElem.Interface(), source.Key().Kind())
				return
			}

			mapKeyElem = mapKeyElem.Convert(source.Key())

			if !mapValueElem.Type().ConvertibleTo(source.Elem()) {
				err = fmt.Errorf("map value %v cannot is convertible to %s", mapValueElem.Interface(), source.Elem().Kind())
				return
			}

			mapValueElem = mapValueElem.Convert(source.Elem())
		}

		mapValue.SetMapIndex(mapKeyElem, mapValueElem)
	}

	result = mapValue.Interface()

	return
}
