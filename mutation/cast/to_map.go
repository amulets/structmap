package cast

import (
	"fmt"
	"reflect"

	"github.com/dungeon-code/structmap"
)

func toMap(field *structmap.FieldPart) error {
	fieldValue := reflect.Indirect(reflect.ValueOf(field.Value))
	if fieldValue.Kind() == reflect.Invalid {
		return nil
	}

	// It is a same type
	if fieldValue.Type().Key().Kind() == field.Type.Key().Kind() && fieldValue.Type().Elem().Kind() == field.Type.Elem().Kind() {
		return nil
	}

	// Field value is map[string]interface{}
	if fieldValue.Type().Key().Kind() != reflect.String && fieldValue.Type().Elem().Kind() != reflect.Interface {
		return fmt.Errorf("supported only struct field of type map[string]interface{} to cast")
	}

	mapType := reflect.MapOf(field.Type.Key(), field.Type.Elem())
	mapValue := reflect.MakeMap(mapType)

	mapIterator := fieldValue.MapRange()
	for mapIterator.Next() {
		mapValueElem := mapIterator.Value().Elem()

		if !mapValueElem.Type().ConvertibleTo(field.Type.Elem()) {
			return fmt.Errorf("map value %v cannot is convertible to %s", mapValueElem.Interface(), field.Type.Elem().Kind())
		}

		mapValue.SetMapIndex(mapIterator.Key(), mapValueElem.Convert(field.Type.Elem()))
	}

	field.Value = mapValue.Interface()

	return nil
}
